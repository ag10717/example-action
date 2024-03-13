package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitHandler struct {
	Repo              *git.Repository
	BranchNameInput   string
	MajorVersionInput string
}

func (gh *GitHandler) IncrementTag(tag string) string {
	o := strings.Split(tag, ".")
	// convert to ints
	major, err := strconv.Atoi(o[0])
	HandleError(err)
	minor, err := strconv.Atoi(o[1])
	HandleError(err)
	patch, err := strconv.Atoi(o[2])
	HandleError(err)

	if gh.BranchNameInput == "main" {
		major = major + 1
		minor = 0
		patch = 0
	} else if StringContains(gh.BranchNameInput, []string{"feature", "task", "bugfix"}) {
		minor = minor + 1
	} else {
		patch = patch + 1
	}

	newTag := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	fmt.Printf("incrementing tag too: %s", newTag)
	return newTag
}

func (gh *GitHandler) SetTag(tag string) {
	head, err := gh.Repo.Head()
	HandleError(err)

	_, err = gh.Repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
	})

	HandleError(err)
}

// see example: https://github.com/go-git/go-git/blob/master/_examples/find-if-any-tag-point-head/main.go
func (gh *GitHandler) GetLatestTag() string {
	var latestTag string

	ref, err := gh.Repo.Head()
	HandleError(err)

	tags, err := gh.Repo.Tags()
	HandleError(err)

	err = tags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := gh.Repo.ResolveRevision(plumbing.Revision(t.Name()))
		HandleError(err)

		if *revHash == ref.Hash() {
			latestTag = t.Name().Short()
		}

		return nil
	})
	HandleError(err)

	fmt.Printf("found latest tag: %s", latestTag)
	return latestTag
}
