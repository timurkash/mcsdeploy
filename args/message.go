package args

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
)

func ArgMes(message string) error {
	if err := filepath.Walk("gen/go", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.HasSuffix(path, "/messages.pb.go") || strings.HasSuffix(path, "/common.pb.go")) {
			if err := findMessage(path, message); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func findMessage(filepath, message string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	found := false
	typeMessageStruct := fmt.Sprintf("type %s struct {", message)
	for _, line := range lines {
		if line == typeMessageStruct {
			fmt.Println()
			fmt.Printf("// %s\n", filepath)
			fmt.Println()
			fmt.Printf("export function get%s(item) {\n", message)
			fmt.Println("\tif (item) {")
			fmt.Println("\t\treturn {")
			found = true
		} else {
			if found {
				if line == "}" {
					fmt.Println("\t\t}")
					fmt.Println("\t}")
					fmt.Println("}")
					return nil
				} else {
					if strings.Contains(line, "`protobuf") {
						processLine(line)
					}
				}
			}
		}
	}
	_ = found
	return nil
}

func processLine(line string) {
	line = strings.Trim(line, "\t")
	lexemes := strings.Split(line, " ")
	name := lexemes[0]
	name_ := strcase.LowerCamelCase(name)
	typ := ""
	for i, lexeme := range lexemes {
		if i > 0 && lexeme != "" {
			typ = lexeme
			break
		}
	}
	fmt.Print("\t\t\t")
	switch {
	case strings.HasPrefix(typ, "[]"):
		fmt.Printf("%s: item.get%sList(),\n", name_, name)
	case strings.HasPrefix(typ, "*common."):
		getFun := "get" + typ[8:]
		fmt.Printf("%s: get%s(item.%s()), // import {%s} from '@/assets/json/common'\n", name_, getFun, name, getFun)
	case strings.HasPrefix(typ, "*"):
		fun := typ[1:]
		fmt.Printf("%s: get%s(item.get%s()),\n", name_, fun, name)
	default:
		fmt.Printf("%s: item.get%s(),\n", name_, name)
	}
}
