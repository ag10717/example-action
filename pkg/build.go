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
	HandleError(err)
	minor, err := strconv.Atoi(o[1])
	HandleError(err)
	patch, err := strconv.Atoi(o[2])
	HandleError(err)

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
	HandleError(err)

	_, err = h.Repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
	})

	HandleError(err)
}

// see example: https://github.com/go-git/go-git/blob/master/_examples/find-if-any-tag-point-head/main.go
func (h *Handler) GetLatestBuild() string {
	var latestTag string

	ref, err := h.Repo.Head()
	HandleError(err)

	tags, err := h.Repo.Tags()
	HandleError(err)

	err = tags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := h.Repo.ResolveRevision(plumbing.Revision(t.Name()))
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

func (h *Handler) PushTag() {}
func (h *Handler) GetBuildEnv() string {
	if os.Getenv("BUILD_NUMBER") != "" {
		return os.Getenv("BUILD_NUMBER")
	}

	return ""
}
