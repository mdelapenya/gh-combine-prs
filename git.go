package main

import (
	"os/exec"
)

func checkoutBranch(branch string) error {
	extensionLogger.Println("Checking out branch", branch)
	err := gitExec("checkout", branch)
	if err != nil {
		return err
	}

	extensionLogger.Println("Branch", branch, "checked out")
	return nil
}

func createBranch(name string, base string) error {
	extensionLogger.Println("Creating branch", name, "from", base)
	err := deleteBranch(name)
	if err != nil {
		extensionLogger.Printf(">> failed to delete branch, ignoring: %s\n", err)
	}

	err = gitExec("checkout", "-b", name, base)
	if err != nil {
		return err
	}

	extensionLogger.Println("Branch", name, "created from", base)
	return nil
}

func deleteBranch(branch string) error {
	extensionLogger.Println("Deleting branch", branch)
	err := gitExec("branch", "-D", branch)
	if err != nil {
		return err
	}

	extensionLogger.Println("Branch", branch, "deleted")
	return nil
}

func mergeBranch(branch string, target string) error {
	extensionLogger.Println("Merging branch ", target, "into", branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("merge", target, "--no-edit")
	if err != nil {
		extensionLogger.Println(">> unable to merge", err)
		return gitExec("merge", "--abort")
	}

	extensionLogger.Println("Branch", target, "merged into", branch)
	return nil
}

func updateBranch(branch string) error {
	extensionLogger.Println("Updating branch ", branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("pull", "origin", branch, "--ff-only")
	if err != nil {
		extensionLogger.Printf(">> failed to pull from origin, trying upstream: %s\n", err)
		err = gitExec("pull", "upstream", branch, "--ff-only")
		if err != nil {
			return err
		}
	}

	extensionLogger.Println("Branch", branch, "updated")
	return nil
}

func gitExec(args ...string) error {
	extensionLogger.Printf("Executing git %s\n", args)

	if !dryRunFlag {
		cmd := exec.Command("git", args...)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	return nil
}
