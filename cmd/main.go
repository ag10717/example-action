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

	fmt.Printf("Hello, Example-Action; %s: %s \n", os.Args[1], os.Args[2])
	fmt.Printf("Current Directory: %s \n", wd)

	r, err := git.PlainOpen(wd)
	pkg.HandleError(err)

	// SETUP REPO
	gh := pkg.GitHandler{
		Repo:              r,
		BranchNameInput:   os.Args[1],
		MajorVersionInput: os.Args[2],
	}

	gh.Repo.Fetch(&git.FetchOptions{
		Tags: git.AllTags,
	})

	// GET & SET TAG
	lt := gh.GetLatestTag()
	gh.IncrementTag(lt)

	// gh.SetTag(it)
}
