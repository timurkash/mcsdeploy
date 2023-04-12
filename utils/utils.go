package utils

import (
	"log"
	"os"
	"strings"
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

func GetTag(image string) string {
	p := strings.LastIndex(image, ":")
	return image[p:]
}

//
//func IsDirExists(dirname string) bool {
//	fi, err := os.Stat(dirname)
//	if err == nil {
//		if fi.IsDir() {
//			return true
//		} else {
//			log.Println(dirname, "is not directory")
//			return false
//		}
//	}
//	if os.IsNotExist(err) {
//		return false
//	}
//	return false
//}
//
//func IsExists(filename string) bool {
//	_, err := os.Stat(filename)
//	if err != nil {
//		if os.IsNotExist(err) {
//			return false
//		}
//	}
//	return true
//}
