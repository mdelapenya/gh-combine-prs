package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/cli/go-gh"
)

func checkoutPR(pr PullRequest) error {
	extensionLogger.Printf("Checking out #%d\n", pr.Number)

	if !dryRunFlag {
		_, err := ghExec("pr", "checkout", fmt.Sprintf("%d", pr.Number))
		if err != nil {
			return err
		}
	}

	return nil
}

func checkPassingChecks(pr PullRequest) (bool, error) {
	extensionLogger.Printf("Checking if #%d is passing Github checks\n", pr.Number)

	stdOut, err := ghExec("pr", "checks", fmt.Sprintf("%d", pr.Number))
	if err != nil {
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

func fetchAndSelectPRs(interactive bool) ([]PullRequest, error) {
	extensionLogger.Printf("Fetching pull rquests using query: %s\n", queryFlag)

	stdOut, err := ghExec("pr", "list", "--search", queryFlag, "--limit", fmt.Sprintf("%d", limitFlag), "--json", "number,headRefName,title")
	if err != nil {
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

func ghExec(args ...string) (bytes.Buffer, error) {
	extensionLogger.Println("Args:", args)

	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		extensionLogger.Printf(">> error while executing gh: %v. Stderr: %s", err, &stdErr)
		return bytes.Buffer{}, err
	}

	return stdOut, nil
}
