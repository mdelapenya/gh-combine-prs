package main

import "fmt"

type PullRequest struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	HeadRefName string `json:"headRefName"`
	URL         string `json:"url"`
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("#%d - %s - %s", pr.Number, pr.Title, pr.URL)
}
