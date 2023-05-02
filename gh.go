package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/cli/go-gh"
)

func checkoutPR(pr PullRequest) error {
	extensionLogger.Debugf("Checking out #%d\n", pr.Number)

	if !dryRunFlag {
		_, err := ghExec("pr", "checkout", fmt.Sprintf("%d", pr.Number))
		if err != nil {
			return err
		}
	}

	return nil
}

func checkPassingChecks(pr PullRequest) (bool, error) {
	extensionLogger.Debugf("Checking if #%d is passing Github checks\n", pr.Number)

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

func checkIfCreatePR(branch string, body string) error {
	create := false
	prompt := &survey.Confirm{
		Message: "Do you want to submit the combined PR?",
	}
	survey.AskOne(prompt, &create)

	const defaultPRTitle = "Combined dependencies PR"

	prTitle := ""
	titlePrompt := &survey.Input{
		Message: "Do you want to change the PR title?",
		Default: defaultPRTitle,
	}
	survey.AskOne(titlePrompt, &prTitle)

	const mainBranch = "main"
	targetBranch := ""
	branchPrompt := &survey.Input{
		Message: "Which branch do you want to send the PR against?",
		Default: mainBranch,
	}
	survey.AskOne(branchPrompt, &targetBranch)

	if create {
		err := pushBranch(branch)
		if err != nil {
			return err
		}

		extensionLogger.Infof("Creating combined PR:")
		extensionLogger.Infof("- Head branch: %s\n - Title: %s\n - Labels: dependencies\n - Body:\n%s\n", branch, prTitle, body)

		fork, err := extractFork()
		if err != nil {
			return err
		}

		whoIam := fmt.Sprintf("%s:%s", whoami(), branch)
		targetFork := ""
		forkPrompt := &survey.Select{
			Message: "Which remote do you want to send the PR against?",
			Options: []string{fmt.Sprintf("%s:%s", fork, branch), whoIam},
			Default: whoIam,
		}
		survey.AskOne(forkPrompt, &targetFork)

		createArgs := []string{"pr", "create", "-B", targetBranch, "--head", targetFork, "--title", prTitle, "--body", body, "--label", "dependencies"}
		extensionLogger.Infof("Running: gh %s\n", strings.Join(createArgs, " "))

		if !dryRunFlag {
			_, err := ghExec(createArgs...)
			if err != nil {
				return err
			}
		}

		extensionLogger.Successf("Done! The combined PR has been sent")
	}

	return nil
}

func extractFork() (string, error) {
	extensionLogger.Debugf("Extracting fork for the current repository\n")

	stdOut, err := ghExec("repo", "view", "--json", "owner", "--template", "{{.owner.login}}")
	if err != nil {
		return "", err
	}

	fork := stdOut.String()
	extensionLogger.Infof("Fork detected: %s\n", fork)
	return fork, nil
}

func fetchAndSelectPRs(interactive bool) ([]PullRequest, error) {
	extensionLogger.Debugf("Fetching pull requests using query: %s\n", queryFlag)

	stdOut, err := ghExec("pr", "list", "--search", queryFlag, "--limit", fmt.Sprintf("%d", limitFlag), "--json", "number,headRefName,title,url")
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
	repoOverride := os.Getenv("GH_REPO")
	if repoOverride != "" {
		args = append(args, "--repo", repoOverride) // for testing purposes
	}

	extensionLogger.Debugf("Args: %v\n", args)

	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		extensionLogger.Errorf("while executing gh: %v. Stderr: %s", err, &stdErr)
		return bytes.Buffer{}, err
	}

	return stdOut, nil
}

func viewPR(pr PullRequest) (string, error) {
	extensionLogger.Debugf("Viewing #%d\n", pr.Number)

	stdOut, err := ghExec("pr", "view", fmt.Sprintf("%d", pr.Number), "--json", "title,author,number", "--template", "{{.title}} (#{{.number}}) @{{.author.login}}")
	if err != nil {
		return "", err
	}

	return stdOut.String(), nil
}

func whoami() string {
	response := struct{ Login string }{}
	err := ghClient.Get("user", &response)
	if err != nil {
		extensionLogger.Errorf("%v", err)
		return "origin"
	}

	return response.Login
}
