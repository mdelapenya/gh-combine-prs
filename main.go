package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cli/go-gh"
)

var helpFlag bool

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")
}

func main() {
	flag.Parse()

	if helpFlag {
		fmt.Println("hi world, this is the gh-multi-merge-prs extension!")
		os.Exit(0)
	}

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("running as %s\n", response.Login)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
