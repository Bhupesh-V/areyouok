package utils

import (
	"flag"
	"log"
)

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Check if element exists inside a slice
func In(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetUserDirectory() string {
	var userDir string

	if len(flag.Args()) == 0 {
		userDir = "."
	} else {
		userDir = flag.Args()[0]
	}
	return userDir
}
