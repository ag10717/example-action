package pkg

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/go-version"
)

type Handler struct {
	Repo              *git.Repository
	BranchNameInput   string
	MajorVersionInput string
}

func (h *Handler) IncrementBuild(tag, run_id string) string {
	fmt.Println("increment build number")
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

func (h *Handler) GetLatestBuild() string {
	tags := h.getTags()

	latestTag := tags[0].String()

	fmt.Printf("found latest tag: %s \n", latestTag)
	return latestTag
}

func (h *Handler) PushTag() {}
func (h *Handler) GetBuildEnv() string {
	fmt.Println("check existing build number")

	ev := os.Environ()
	for _, e := range ev {
		fmt.Printf("VAR: %s \n", e)
	}
	fmt.Println(os.Getenv("GITHUB_ENV"))
	fmt.Println(os.Getenv("BUILD_NUMBER"))

	envPath := os.Getenv("GITHUB_ENV")
	file, err := os.ReadFile(envPath)
	HandleError(err, "open env file")

	fmt.Println("FILEDATA: ", string(file))

	if os.Getenv("BUILD_NUMBER") != "" {
		return os.Getenv("BUILD_NUMBER")
	}

	return ""
}

func (h *Handler) getTags() []*version.Version {
	fmt.Println("get all tags")
	var at []string

	tags, err := h.Repo.Tags()
	HandleError(err, "get tags")

	err = tags.ForEach(func(t *plumbing.Reference) error {
		at = append(at, t.Name().Short())

		return nil
	})
	HandleError(err, "tag iter")

	versions := make([]*version.Version, len(at))
	for i, r := range at {
		v, _ := version.NewVersion(r)
		versions[i] = v
	}
	sort.Sort(sort.Reverse(version.Collection(versions)))

	fmt.Printf("have versions: %v \n", versions)
	return versions
}
