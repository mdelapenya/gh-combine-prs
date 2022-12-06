package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

var helpFlag bool
var ghClient api.RESTClient

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show help for multi-merge-prs")

	client, err := gh.RESTClient(nil)
	if err != nil {
		panic(err)
	}
	ghClient = client
}

func main() {
	flag.Parse()

	if helpFlag {
		fmt.Println("hi world, this is the gh-multi-merge-prs extension!")
		os.Exit(0)
	}

	whoami()
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
