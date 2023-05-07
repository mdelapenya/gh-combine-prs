package main

import (
	"os/exec"

	survey "github.com/AlecAivazis/survey/v2"
)

func checkoutBranch(branch string) error {
	extensionLogger.Debugf("Checking out branch %s\n", branch)
	err := gitExec("checkout", branch)
	if err != nil {
		return err
	}

	extensionLogger.Successf("Branch %s checked out\n", branch)
	return nil
}

func createBranch(name string, base string) error {
	extensionLogger.Debugf("Creating branch %s from %s\n", name, base)
	err := deleteBranch(name)
	if err != nil {
		extensionLogger.Errorf("failed to delete branch, ignoring: %s\n", err)
	}

	err = gitExec("checkout", "-b", name, base)
	if err != nil {
		return err
	}

	extensionLogger.Successf("Branch %s created from %s\n", name, base)
	return nil
}

func deleteBranch(branch string) error {
	extensionLogger.Debugf("Deleting branch %s\n", branch)
	err := gitExec("branch", "-D", branch)
	if err != nil {
		return err
	}

	extensionLogger.Successf("Branch %s deleted\n", branch)
	return nil
}

func mergeBranch(branch string, target string) error {
	extensionLogger.Debugf("Merging branch  %s into %s\n", target, branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("merge", target, "--no-edit")
	if err != nil {
		extensionLogger.Errorf("unable to merge: %v", err)

		resolved := false
		prompt := &survey.Confirm{
			Message: "We have found merge conflicts. Please go back to your editor, resolve them including commit, and come back to this prompt to press 'Yes'. Did you resolve them?",
		}
		survey.AskOne(prompt, &resolved)

		if !resolved {
			return gitExec("merge", "--abort")
		}
	}

	extensionLogger.Successf("Branch %s merged into %s\n", target, branch)
	return nil
}

func pushBranch(branch string) error {
	extensionLogger.Debugf("Pushing branch %s\n", branch)

	err := gitExec("push", "origin", branch)
	if err != nil {
		extensionLogger.Errorf("unable to push: %v", err)
		return err
	}

	extensionLogger.Successf("Branch %s pushed to origin", branch)
	return nil
}

func updateBranch(branch string) error {
	extensionLogger.Debugf("Updating branch %s\n", branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("pull", "origin", branch, "--ff-only")
	if err != nil {
		extensionLogger.Errorf("failed to pull from origin, trying upstream: %s\n", err)
		err = gitExec("pull", "upstream", branch, "--ff-only")
		if err != nil {
			return err
		}
	}

	extensionLogger.Successf("Branch %s updated\n", branch)
	return nil
}

func gitExec(args ...string) error {
	extensionLogger.Debugf("Executing git %s\n", args)

	if !dryRunFlag {
		cmd := exec.Command("git", args...)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	return nil
}
