package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
)

var helpFlag bool
var limitFlag int
var queryFlag string
var skipPRCheckFlag bool

var ghClient api.RESTClient
var currentRepo repository.Repository

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")
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

	selectedPRs, err := selectPRsPrompt()
	if err != nil {
		panic(err)
	}

	if len(selectedPRs) == 0 {
		fmt.Println("No PRs selected to merge. Exiting")
		os.Exit(0)
	}

	var passingPRs []PullRequest
	fmt.Printf("Selected PRs: %v\n", selectedPRs)
	for _, pr := range selectedPRs {
		passing, err := checkPassingChecks(pr)
		if err != nil {
			panic(err)
		}

		if passing {
			passingPRs = append(passingPRs, pr)
		} else {
			fmt.Printf("Not all checks are passing for #%d, skipping PR", pr.Number)
		}
	}

	whoami()
}

func checkPassingChecks(pr PullRequest) (bool, error) {
	args := []string{"pr", "checks", fmt.Sprintf("%d", pr.Number)}

	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		fmt.Println(stdErr)
		return false, err
	}

	checks := stdOut.String()
	checksList := strings.Split(checks, "\n")
	for _, check := range checksList {
		if strings.Contains(check, "fail") || strings.Contains(check, "pending") {
			return false, nil
		}
	}

	return true, nil
}

func selectPRsPrompt() ([]PullRequest, error) {
	args := []string{"pr", "list", "--search", queryFlag, "--limit", fmt.Sprintf("%d", limitFlag), "--json", "number,headRefName,title"}

	fmt.Println("Args:", args)

	fmt.Println("The following PRs will be evaluated for inclusion:")
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		fmt.Println(stdErr)
		return nil, err
	}

	var prs []PullRequest
	err = json.Unmarshal(stdOut.Bytes(), &prs)
	if err != nil {
		return nil, err
	}

	prOptions := make([]string, len(prs))
	for i, pr := range prs {
		prOptions[i] = pr.String()
	}

	var selectedPrs []string
	survey.AskOne(&survey.MultiSelect{
		Message: "Select the PRs to combine",
		Options: prOptions,
	}, &selectedPrs, survey.WithRemoveSelectAll())

	result := []PullRequest{}
	for _, selectedPr := range selectedPrs {
		for _, pr := range prs {
			if pr.String() == selectedPr {
				result = append(result, pr)
			}
		}
	}

	return result, nil
}

func usage(exitCode int, args ...string) {
	for _, arg := range args {
		fmt.Fprintln(os.Stderr, arg)
	}

	fmt.Println(`Usage: gh combine-prs --query "QUERY" [--limit 50] [--selected-pr-numbers 42,13,78] [--skip-pr-check] [--help]
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
