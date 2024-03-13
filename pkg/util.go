package pkg

import (
	"fmt"
	"log"
	"strings"
)

func HandleError(err error, message string) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("%s \n %v", message, err))
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
