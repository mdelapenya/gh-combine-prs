package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
)

var dryRunFlag bool
var helpFlag bool
var interactiveFlag bool
var limitFlag int
var queryFlag string
var skipPRCheckFlag bool
var verboseFlag bool

var ghClient api.RESTClient
var currentRepo repository.Repository

var extensionLogger Logger

const combinedPRsBranchName = "combined-pr-branch"

func init() {
	flag.BoolVar(&dryRunFlag, "dry-run", false, "If set, will not actually merge the PRs, forcing verbose mode to show internal steps. Defaults to false when not specified")
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")
	flag.BoolVar(&interactiveFlag, "interactive", false, "Enable interactive mode. If set, will prompt for selecting the PRs to merge")
	flag.IntVar(&limitFlag, "limit", 50, "Sets the maximum number of PRs that will be combined. Defaults to 50")
	flag.StringVar(&queryFlag, "query", "", `sets the query used to find combinable PRs. e.g. --query "author:app/dependabot to combine Dependabot PRs`)
	flag.BoolVar(&skipPRCheckFlag, "skip-pr-check", false, `if set, will combine matching PRs even if they are not passing checks. Defaults to false when not specified`)
	flag.BoolVar(&verboseFlag, "verbose", false, `if set, will print verbose output. Defaults to false when not specified`)
}

func main() {
	flag.Parse()

	if dryRunFlag {
		// force verbose mode when dry-running
		verboseFlag = true
	}

	extensionLogger = newLogger(verboseFlag)

	if helpFlag {
		usage(0)
	}

	if queryFlag == "" {
		usage(1, "ERROR: --query is required")
	}

	client, err := gh.RESTClient(nil)
	if err != nil {
		panic(err)
	}
	ghClient = client

	repo, err := gh.CurrentRepository()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current repository is %s/%s\n", repo.Owner(), repo.Name())

	currentRepo = repo

	extensionLogger.Debugf("Dry-run mode: %t\n", dryRunFlag)

	selectedPRs, err := fetchAndSelectPRs(interactiveFlag)
	if err != nil {
		extensionLogger.Errorf("while fetching the PRs. Exiting: %v\n", err)
		os.Exit(1)
	}

	if len(selectedPRs) == 0 {
		extensionLogger.Warnf("No PRs selected to merge. Exiting")
		os.Exit(0)
	}

	var confirmedPRs []PullRequest
	extensionLogger.Debugf("Selected PRs:")
	var errors []error
	for _, pr := range selectedPRs {
		if skipPRCheckFlag {
			extensionLogger.Debugf("%s\n", pr)
			confirmedPRs = append(confirmedPRs, pr)
			continue
		}

		passing, err := checkPassingChecks(pr)
		if err != nil {
			extensionLogger.Warnf("while fetching Github checks for #%d, skipping PR: %v\n", pr.Number, err)
			errors = append(errors, err)
			continue
		}

		if passing {
			extensionLogger.Debugf("%s\n", pr)
			confirmedPRs = append(confirmedPRs, pr)
		} else {
			extensionLogger.Warnf("Not all checks are passing for #%d, skipping PR", pr.Number)
		}
	}

	if len(errors) == len(selectedPRs) {
		extensionLogger.Errorf("All PRs failed to pass checks. Exiting")
		os.Exit(1)
	}

	// checkout default branch
	defaultBranch, err := defaultBranch()
	if err != nil {
		panic(err)
	}
	extensionLogger.Debugf("default branch is %s\n", defaultBranch)

	err = updateBranch(defaultBranch)
	if err != nil {
		panic(err)
	}

	var prTitle = combineTitles(confirmedPRs)

	branchName := fmt.Sprintf("%s-%s", combinedPRsBranchName, titlesHash(prTitle))

	err = createBranch(branchName, defaultBranch)
	if err != nil {
		panic(err)
	}

	executable := []string{"gh", "combine-prs"}
	executable = append(executable, os.Args[1:]...)

	command := strings.Join(executable, " ")
	disclaimer := "> [!NOTE]\n>This PR has been created with the [combine-prs](https://github.com/mdelapenya/gh-combine-prs) `gh` extension:\n\n>" + command + ".\n\n"
	body := disclaimer + "It combines the following PRs:\n\n"
	relatedIssuesText := "## Related Issues:\n\n"

	for _, pr := range confirmedPRs {
		err = checkoutPR(pr)
		if err != nil {
			panic(err)
		}
		err = mergeBranch(branchName, pr.HeadRefName)
		if err != nil {
			extensionLogger.Warnf("pull request #%d failed to merge into %s: %v. Skipping PR\n", pr.Number, branchName, err)
			continue
		}

		relatedIssuesText += fmt.Sprintf("- Closes #%d\n", pr.Number)
		prDescription, err := viewPR(pr)
		if err != nil {
			panic(err)
		}
		body += fmt.Sprintf("- %s\n", prDescription)
	}

	if len(confirmedPRs) > 0 {
		body += "\n" + relatedIssuesText
	}

	err = checkIfCreatePR(branchName, prTitle, body)
	if err != nil {
		panic(err)
	}
}

func defaultBranch() (string, error) {
	response := struct {
		DefaultBranch string `json:"default_branch"`
	}{}
	err := ghClient.Get("repos/"+currentRepo.Owner()+"/"+currentRepo.Name(), &response)
	if err != nil {
		return "", err
	}

	return response.DefaultBranch, nil
}

// titlesHash returns a hash of the PR titles
func titlesHash(prTitle string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(prTitle))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d", h.Sum32())
}

func usage(exitCode int, args ...string) {
	for _, arg := range args {
		fmt.Fprintln(os.Stderr, arg)
	}

	fmt.Println(`Usage: gh multi-merge-prs --query "QUERY" [--limit 50] [--skip-pr-check] [--verbose] [--interactive] [--help]
Arguments:
	`)
	maxLength := 0
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > maxLength {
			maxLength = len(f.Name)
		}
	})
	flag.VisitAll(func(f *flag.Flag) {
		currentLength := len(f.Name)
		fmt.Fprintf(os.Stderr, "  --%s%s%s\n", f.Name, strings.Repeat(" ", maxLength-currentLength+3), f.Usage)
	})

	// exit execution after printing usage
	os.Exit(exitCode)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
