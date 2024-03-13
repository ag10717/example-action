package pkg

import (
	"log"
	"strings"
)

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err)
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
