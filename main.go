package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
)

var helpFlag bool
var interactiveFlag bool
var limitFlag int
var queryFlag string
var skipPRCheckFlag bool

var ghClient api.RESTClient
var currentRepo repository.Repository

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")
	flag.BoolVar(&interactiveFlag, "interactive", false, "Enable interactive mode. If set, will prompt for selecting the PRs to merge")
	flag.IntVar(&limitFlag, "limit", 50, "Sets the maximum number of PRs that will be combined. Defaults to 50")
	flag.StringVar(&queryFlag, "query", "", `sets the query used to find combinable PRs. e.g. --query "author:app/dependabot to combine Dependabot PRs`)
	flag.BoolVar(&skipPRCheckFlag, "skip-pr-check", false, `if set, will combine matching PRs even if they are not passing checks. Defaults to false when not specified`)

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
}

func main() {
	flag.Parse()

	if helpFlag {
		usage(0)
	}

	if queryFlag == "" {
		usage(1, "ERROR: --query is required")
	}

	selectedPRs, err := selectPRs(interactiveFlag)
	if err != nil {
		panic(err)
	}

	if len(selectedPRs) == 0 {
		fmt.Println("No PRs selected to merge. Exiting")
		os.Exit(0)
	}

	var confirmedPRs []PullRequest
	fmt.Println("Selected PRs:")
	for _, pr := range selectedPRs {
		if skipPRCheckFlag {
			fmt.Printf("%s\n", pr)
			confirmedPRs = append(confirmedPRs, pr)
			continue
		}

		passing, err := checkPassingChecks(pr)
		if err != nil {
			panic(err)
		}

		if passing {
			fmt.Printf("%s\n", pr)
			confirmedPRs = append(confirmedPRs, pr)
		} else {
			fmt.Printf("Not all checks are passing for #%d, skipping PR", pr.Number)
		}
	}

	// checkout default branch
	defaultBranch, err := defaultBranch()
	if err != nil {
		panic(err)
	}
	fmt.Printf("default branch is %s\n", defaultBranch)

	mergeBranchName := "multi-merge-pr-branch"

	err = updateBranch(defaultBranch)
	if err != nil {
		panic(err)
	}
	err = createBranch(mergeBranchName, defaultBranch)
	if err != nil {
		panic(err)
	}

	for _, pr := range confirmedPRs {
		err = checkoutPR(pr)
		if err != nil {
			panic(err)
		}
		err = mergeBranch(mergeBranchName, pr.HeadRefName)
		if err != nil {
			fmt.Printf(">> Pull request #%d failed to merge into %s. Skipping PR\n", pr.Number, mergeBranchName)
			continue
		}
	}

	// send PR to merge the new branch into the default branch

	whoami()
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

func usage(exitCode int, args ...string) {
	for _, arg := range args {
		fmt.Fprintln(os.Stderr, arg)
	}

	fmt.Println(`Usage: gh multi-merge-prs --query "QUERY" [--limit 50] [--skip-pr-check] [--help]
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

func whoami() {
	response := struct{ Login string }{}
	err := ghClient.Get("user", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("running as %s\n", response.Login)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
