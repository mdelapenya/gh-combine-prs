package main

import "strings"

func combineTitles(prs []PullRequest) string {
	if len(prs) == 1 {
		return prs[0].Title
	}

	var title string

	const bumpPrefix string = "chore(deps): bump"
	const udpatePrefix string = "chore(deps): update"
	if strings.HasPrefix(prs[0].Title, bumpPrefix) {
		title = bumpPrefix
	} else if strings.HasPrefix(prs[0].Title, udpatePrefix) {
		title = udpatePrefix
	}

	for i, pr := range prs {
		t := strings.ReplaceAll(pr.Title, bumpPrefix, "")
		t = strings.ReplaceAll(t, udpatePrefix, "")
		t = strings.TrimSpace(t)

		title += " " + t
		if i < len(prs)-1 {
			title += ","
		}
	}

	return title
}
