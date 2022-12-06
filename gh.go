package main

import (
	"encoding/json"
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/cli/go-gh"
)

func checkoutPR(pr PullRequest) error {
	args := []string{"pr", "checkout", fmt.Sprintf("%d", pr.Number)}

	_, stdErr, err := gh.Exec(args...)
	if err != nil {
		fmt.Println(stdErr)
		return err
	}

	return nil
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

func selectPRs(interactive bool) ([]PullRequest, error) {
	args := []string{"pr", "list", "--search", queryFlag, "--limit", fmt.Sprintf("%d", limitFlag), "--json", "number,headRefName,title"}

	fmt.Println("Args:", args)

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

	if !interactive {
		// return the response from the API
		return prs, nil
	}

	// because we are in interactive mode, we need to prompt the user to select the PRs to merge

	prOptions := make([]string, len(prs))
	for i, pr := range prs {
		prOptions[i] = pr.String()
	}

	var selectedPrs []string
	survey.AskOne(&survey.MultiSelect{
		Message: "Please select the PRs to combine",
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
