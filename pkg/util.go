package pkg

import (
	"log"
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
