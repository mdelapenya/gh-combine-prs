package main

import "testing"

func TestCombineTitle(t *testing.T) {
	prs := []PullRequest{
		{
			Title: "chore(deps): update baaz from 3.2.1 to 3.2.2 in /",
		},
	}

	title := combineTitles(prs)
	expected := "chore(deps): update baaz from 3.2.1 to 3.2.2 in /"
	if title != expected {
		t.Errorf("expected %q, got %q", expected, title)
	}
}

func TestCombineTitlesBumps(t *testing.T) {
	prs := []PullRequest{
		{
			Title: "chore(deps): bump foo from 1.0.0 to 1.0.1 in /",
		},
		{
			Title: "chore(deps): bump bar from 1.0.0 to 1.0.1 in /modules/bar",
		},
		{
			Title: "chore(deps): update baaz from 3.2.1 to 3.2.2 in /",
		},
		{
			Title: "chore(deps): bump faag from 1.0.0 to 1.0.1 in /foo",
		},
	}

	title := combineTitles(prs)
	expected := "chore(deps): bump foo from 1.0.0 to 1.0.1 in /, bar from 1.0.0 to 1.0.1 in /modules/bar, baaz from 3.2.1 to 3.2.2 in /, faag from 1.0.0 to 1.0.1 in /foo"
	if title != expected {
		t.Errorf("expected %q, got %q", expected, title)
	}
}

func TestCombineTitlesUpdates(t *testing.T) {
	prs := []PullRequest{
		{
			Title: "chore(deps): update foo from 1.0.0 to 1.0.1 in /",
		},
		{
			Title: "chore(deps): bump bar from 1.0.0 to 1.0.1 in /modules/bar",
		},
		{
			Title: "chore(deps): bump baaz from 3.2.1 to 3.2.2 in /",
		},
		{
			Title: "chore(deps): bump faag from 1.0.0 to 1.0.1 in /foo",
		},
	}

	title := combineTitles(prs)
	expected := "chore(deps): update foo from 1.0.0 to 1.0.1 in /, bar from 1.0.0 to 1.0.1 in /modules/bar, baaz from 3.2.1 to 3.2.2 in /, faag from 1.0.0 to 1.0.1 in /foo"
	if title != expected {
		t.Errorf("expected %q, got %q", expected, title)
	}
}
