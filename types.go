package main

import "fmt"

type PullRequest struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	HeadRefName string `json:"headRefName"`
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("#%d - %s", pr.Number, pr.Title)
}
