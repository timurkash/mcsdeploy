package utils

import (
	"log"
	"os"
)

type Util struct {
	Name    string
	Command string
}

func IsFileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil {
		if !fi.IsDir() {
			return true
		} else {
			log.Println(filename, "is directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
