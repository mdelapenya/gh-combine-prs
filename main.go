package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

var helpFlag bool
var limitFlag int
var queryFlag string
var selectedPRNumbersFlag string
var skipPRCheckFlag bool

var ghClient api.RESTClient

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")
	flag.IntVar(&limitFlag, "limit", 50, "Sets the maximum number of PRs that will be combined. Defaults to 50")
	flag.StringVar(&queryFlag, "query", "", `sets the query used to find combinable PRs. e.g. --query "author:app/dependabot to combine Dependabot PRs`)
	flag.StringVar(&selectedPRNumbersFlag, "selected-pr-numbers", "", `COMMA,SEPARATED,LIST. If set, will only work on PRs with the selected numbers. e.g. --selected-pr-numbers 42,13,78`)
	flag.BoolVar(&skipPRCheckFlag, "skip-pr-check", false, `if set, will combine matching PRs even if they are not passing checks. Defaults to false when not specified`)

	client, err := gh.RESTClient(nil)
	if err != nil {
		panic(err)
	}
	ghClient = client
}

func main() {
	flag.Parse()

	if helpFlag {
		usage(0)
	}

	if queryFlag == "" {
		usage(1, "ERROR: --query is required")
	}

	whoami()
}

func usage(exitCode int, args ...string) {
	for _, arg := range args {
		fmt.Fprintln(os.Stderr, arg)
	}

	fmt.Println(`Usage: gh combine-prs --query "QUERY" [--limit 50] [--selected-pr-numbers 42,13,78] [--skip-pr-check] [--help]
Arguments:
	`)
	maxLength := 0
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > maxLength {
			maxLength = len(f.Name)
		}
	})
	flag.VisitAll(func(f *flag.Flag) {
		currentLength := len(f.Name)
		fmt.Fprintf(os.Stderr, "  --%s%s%s\n", f.Name, strings.Repeat(" ", maxLength-currentLength+3), f.Usage)
	})

	// exit execution after printing usage
	os.Exit(exitCode)
}

func whoami() {
	response := struct{ Login string }{}
	err := ghClient.Get("user", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("running as %s\n", response.Login)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
