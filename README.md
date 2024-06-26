# gh-combine-prs

This is an extension for [GitHub CLI](https://cli.github.com/) that combines multiple PRs into one.
It is intended for use in repositories that receive many PRs that can be merged simultaneously, e.g. trivial Dependabot version bump PRs.

The tool will attempt to create one PR that contains all PRs that:

* match a provided query - e.g. `--query "author:app/dependabot"` so that only Dependabot PRs are processed
* and have checks passing
* and that can be merged cleanly - e.g. if two combinable PRs conflict with one another, it will allow you to resolve the conflicts and continue

This tool does not automerge into the `master`/`main` branch - it just attempts to create one unified PR for review and automated checks to run against.

*Note: When you merge the combined PR, it is recommended that you create a Merge Commit.
This allows GitHub to automatically detect that all of the original combined PRs have been merged, so that their state can be set correctly.*

## Inspiration

This tool has been created after using @rnorth's [combine-prs](https://github.com/rnorth/gh-combine-prs) tool for a while, but creating an interactive version that allows you to select which PRs to combine.

## Installation

Prerequisites:
 * `git` is installed (obviously)
 * [GitHub CLI](https://cli.github.com/) is already installed and authenticated

To install this extension:

```
gh extension install mdelapenya/gh-combine-prs
```

## Usage

```
cd $DIRECTORY_OF_YOUR_REPO

gh combine-prs --query "author:app/dependabot" --interactive --verbose --skip-pr-check
```

### Required arguments
    --query "QUERY"
            sets the query used to find combinable PRs.
            e.g. --query "author:app/dependabot"
            to combine Dependabot PRs

### Optional arguments
    --dry-run
            If set, will not actually merge the PRs, forcing verbose mode to show internal steps. Defaults to false when not specified
    --interactive
            Enable interactive mode. If set, will prompt for selecting the PRs to merge
    --limit LIMIT
            sets the maximum number of PRs that will be combined.
            Defaults to 50
    --skip-pr-check
            if set, will combine matching PRs even if they are not passing checks.
            Defaults to false when not specified
    --verbose
            if set, will print verbose output. Defaults to false when not specified

## License

See [LICENSE](./LICENSE)
