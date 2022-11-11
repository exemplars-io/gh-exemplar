package main

/* import (
	"fmt"

	"github.com/cli/go-gh"
	"github.com/exemplars-io/gh-exemplar/internal"
)

func main() {
	fmt.Println("hi world, this is the gh-exemplar extension!")

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
	fmt.Printf("you are running as %s\n", response.Login)
	var profile = internal.GetNilUserProfile() // List repos
	internal.ListExemplars("k8s", true, profile)
} */

import "github.com/exemplars-io/gh-exemplar/cmd"

func main() {
	cmd.Execute()
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
