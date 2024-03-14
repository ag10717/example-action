package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func HandleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s \n %v \n", message, err)
	}
}

func StringContains(baseValue string, values []string) bool {
	for _, s := range values {
		if strings.Contains(baseValue, s) {
			return true
		}
	}

	return false
}

func GetBuildType(branchName string) string {
	if strings.Contains(branchName, "main") {
		return "release"
	}

	return "feature"
}

func WriteGithubEnvValue(name, value string) {
	filePath := os.Getenv("GITHUB_ENV")

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	HandleError(err, "open file")

	defer f.Close()

	fmt.Fprintf(f, "%s=%s \n", name, value)
}
