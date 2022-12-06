package main

import (
	"fmt"
	"os/exec"
)

func checkoutBranch(branch string) error {
	fmt.Println("Checking out branch", branch)
	err := gitExec("checkout", branch)
	if err != nil {
		return err
	}

	fmt.Println("Branch", branch, "checked out")
	return nil
}

func createBranch(name string, base string) error {
	fmt.Println("Creating branch", name, "from", base)
	err := deleteBranch(name)
	if err != nil {
		fmt.Printf(">> failed to delete branch, ignoring: %s\n", err)
	}

	err = gitExec("checkout", "-b", name, base)
	if err != nil {
		return err
	}

	fmt.Println("Branch", name, "created from", base)
	return nil
}

func deleteBranch(branch string) error {
	fmt.Println("Deleting branch", branch)
	err := gitExec("branch", "-D", branch)
	if err != nil {
		return err
	}

	fmt.Println("Branch", branch, "deleted")
	return nil
}

func mergeBranch(branch string, target string) error {
	fmt.Println("Merging branch ", target, "into", branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("merge", target, "--no-edit")
	if err != nil {
		fmt.Println(">> unable to merge", err)
		return gitExec("merge", "--abort")
	}

	fmt.Println("Branch", target, "merged into", branch)
	return nil
}

func updateBranch(branch string) error {
	fmt.Println("Updating branch ", branch)

	err := checkoutBranch(branch)
	if err != nil {
		return err
	}

	err = gitExec("pull", "origin", branch, "--ff-only")
	if err != nil {
		fmt.Printf(">> failed to pull from origin, trying upstream: %s\n", err)
		err = gitExec("pull", "upstream", branch, "--ff-only")
		if err != nil {
			return err
		}
	}

	fmt.Println("Branch", branch, "updated")
	return nil
}

func gitExec(args ...string) error {
	fmt.Printf("Executing git %s\n", args)
	cmd := exec.Command("git", args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
