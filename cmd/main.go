package main

import (
	"fmt"
	"os"

	"github.com/ag10717/example-action/pkg"
	"github.com/go-git/go-git/v5"
)

func main() {
	wd, err := os.Getwd()
	pkg.HandleError(err)

	fmt.Println(os.Args[1:])

	fmt.Printf("Hello, Example-Action; %s: %s \n", os.Getenv("GITHUB_REF"), os.Args[2])
	fmt.Printf("Current Directory: %s \n", wd)

	r, err := git.PlainOpen(wd)
	pkg.HandleError(err)

	// SETUP REPO
	gh := pkg.Handler{
		Repo:              r,
		BranchNameInput:   os.Getenv("GITHUB_REF"),
		MajorVersionInput: os.Args[2],
	}

	gh.Repo.Fetch(&git.FetchOptions{
		Tags: git.AllTags,
	})

	// GET & SET TAG
	var bn string

	// check if the build number has already be injected into the container
	// this might happen if you run this action without create_tag; and then run it again with create_tag
	bn = gh.GetBuildEnv()
	pkg.HandleError(err)

	if bn == "" {
		bn = gh.GetLatestBuild()
		gh.IncrementBuild(bn, os.Getenv("GITHUB_RUN_ID"))

		os.Setenv("GITHUB_ENV", fmt.Sprintf("BUILD_NUMBER=%s", bn))
	}

	// if os.Args[1] == "true" {
	// 	// gh.SetTag(bn)
	// 	// gh.PushTag()
	// }
}
