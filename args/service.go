package args

import (
	"fmt"
	"os"
	"strings"
)

func ArgService(service string) error {
	files, err := os.ReadDir("api")
	if err != nil {
		return err
	}
	found := false
	for _, file := range files {
		if file.IsDir() && file.Name() == service {
			found = true
			if err := findMessagesAndEnums(fmt.Sprintf("api/%s/messages.proto", service)); err != nil {
				return err
			}
		}
	}
	if !found {
		fmt.Printf("service %s not found", service)
	}
	return nil
}

const (
	unknown = 0
	message = 1
	enum    = 2
)

func findMessagesAndEnums(filepath string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		typ := unknown
		switch {
		case strings.HasPrefix(line, "message "):
			typ = message
		case strings.HasPrefix(line, "enum "):
			typ = enum
		}
		if typ == message || typ == enum {
			lexemes := clearSlice(strings.Split(line, " "))
			name := lexemes[1]
			if typ == message {
				if strings.HasSuffix(name, "Request") {
					name = name[:len(name)-7]
					if strings.HasPrefix(name, "Act") {
						name = name[3:]
					}
					fmt.Printf("mcsdeploy -req %sInfo | pbcopy\n", name)
				} else {
					fmt.Printf("mcsdeploy -msg %s | pbcopy\n", name)
				}
			} else {
				fmt.Printf("mcsdeploy -enm %s | pbcopy\n", name)
			}
		}
	}
	return nil
}

func clearSlice(slice []string) (result []string) {
	for _, sl := range slice {
		if sl != "" && sl != "=" && sl != "{" {
			result = append(result, sl)
		}
	}
	return
}
