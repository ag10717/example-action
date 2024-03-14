package main

import (
	"fmt"
	"os"

	"github.com/ag10717/example-action/pkg"
	"github.com/go-git/go-git/v5"
)

func main() {
	wd, err := os.Getwd()
	pkg.HandleError(err, "get cwd")

	fmt.Println(os.Args[1:])

	fmt.Printf("Hello, Example-Action; %s: %s \n", os.Getenv("GITHUB_REF"), os.Args[2])
	fmt.Printf("Current Directory: %s \n", wd)

	r, err := git.PlainOpen(wd)
	pkg.HandleError(err, "open repo")

	// SETUP REPO
	gh := pkg.Handler{
		Repo:              r,
		BranchNameInput:   os.Getenv("GITHUB_REF"),
		MajorVersionInput: os.Args[2],
	}

	rem, err := gh.Repo.Remote("origin")
	pkg.HandleError(err, "get remote")
	fmt.Println(rem.Config().URLs)

	// GET & SET TAG
	var bn string
	var ib string

	// // check if the build number has already be injected into the container
	// // this might happen if you run this action without create_tag; and then run it again with create_tag
	bn = gh.GetBuildEnv()
	pkg.HandleError(err, "get build env")

	if bn == "" {
		bn = gh.GetLatestBuild()
		ib = gh.IncrementBuild(bn, os.Getenv("GITHUB_RUN_ID"))

		pkg.WriteGithubEnvValue("BUILD_NUMBER", ib)
	}

	if os.Args[1] == "true" {
		gh.SetTag(ib)
		gh.PushTag(ib)
	}
}
