package main

import "testing"

func TestCombineTitles(t *testing.T) {
	tests := []struct {
		name string
		prs  []PullRequest
		want string
	}{
		{
			name: "Dependabot: single PR",
			prs: []PullRequest{
				{
					Title: "chore(deps): update baaz from 3.2.1 to 3.2.2 in /",
				},
			},
			want: "chore(deps): update baaz from 3.2.1 to 3.2.2 in /",
		},
		{
			name: "Non-Dependabot: single PR",
			prs: []PullRequest{
				{
					Title: "feat: implement features",
				},
			},
			want: "feat: implement features",
		},
		{
			name: "Multiple dependabot PRs with bump",
			prs: []PullRequest{
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
			},
			want: "chore(deps): bump foo from 1.0.0 to 1.0.1 in /, bar from 1.0.0 to 1.0.1 in /modules/bar, baaz from 3.2.1 to 3.2.2 in /, faag from 1.0.0 to 1.0.1 in /foo",
		},
		{
			name: "Multiple dependabot PRs with update",
			prs: []PullRequest{
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
			},
			want: "chore(deps): update foo from 1.0.0 to 1.0.1 in /, bar from 1.0.0 to 1.0.1 in /modules/bar, baaz from 3.2.1 to 3.2.2 in /, faag from 1.0.0 to 1.0.1 in /foo",
		},
		{
			name: "Multiple PRs with bump",
			prs: []PullRequest{
				{
					Title: "feat 1",
				},
				{
					Title: "feat 2",
				},
				{
					Title: "chore 1",
				},
				{
					Title: "chore 2",
				},
			},
			want: "feat 1, feat 2, chore 1, chore 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title := combineTitles(tt.prs)
			if title != tt.want {
				t.Errorf("expected %q, got %q", tt.want, title)
			}
		})
	}
}
