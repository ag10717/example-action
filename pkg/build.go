package pkg

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Handler struct {
	Repo              *git.Repository
	BranchNameInput   string
	MajorVersionInput string
}

func (h *Handler) IncrementBuild(tag, run_id string) string {
	buildType := GetBuildType(h.BranchNameInput)
	o := strings.Split(tag, ".")
	// convert to ints
	major, err := strconv.Atoi(o[0])
	HandleError(err, "major convert")
	minor, err := strconv.Atoi(o[1])
	HandleError(err, "minor convert")
	patch, err := strconv.Atoi(o[2])
	HandleError(err, "patch convert")

	if h.BranchNameInput == "main" {
		major = major + 1
		minor = 0
		patch = 0
	} else if StringContains(h.BranchNameInput, []string{"feature", "task", "bugfix"}) {
		minor = minor + 1
	} else {
		patch = patch + 1
	}

	buildNumber := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	if buildType == "feature" {
		buildNumber = fmt.Sprintf("%s.%s", buildNumber, fmt.Sprintf("PREVIEW-%s", run_id))
	}

	fmt.Printf("incrementing tag too: %s", buildNumber)
	return buildNumber
}

func (h *Handler) SetTag(tag string) {
	head, err := h.Repo.Head()
	HandleError(err, "get head")

	_, err = h.Repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
	})

	HandleError(err, "create tag")
}

// see example: https://github.com/go-git/go-git/blob/master/_examples/find-if-any-tag-point-head/main.go
func (h *Handler) GetLatestBuild() string {
	var latestTag string

	ref, err := h.Repo.Head()
	HandleError(err, "get head")

	tags, err := h.Repo.Tags()
	HandleError(err, "get tags")

	err = tags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := h.Repo.ResolveRevision(plumbing.Revision(t.Name()))
		HandleError(err, fmt.Sprintf("resolve revision: %s", t.Name()))

		if *revHash == ref.Hash() {
			latestTag = t.Name().Short()
		}

		return nil
	})
	HandleError(err, "tag iter")

	fmt.Printf("found latest tag: %s", latestTag)
	return latestTag
}

func (h *Handler) PushTag() {}
func (h *Handler) GetBuildEnv() string {
	if os.Getenv("BUILD_NUMBER") != "" {
		return os.Getenv("BUILD_NUMBER")
	}

	return ""
}
