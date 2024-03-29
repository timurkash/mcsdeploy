package args

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
)

func ArgMessage(message string) error {
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
	if !strings.HasSuffix(message, "Request") {
		for _, line := range lines {
			if line == typeMessageStruct {
				fmt.Printf("// %s\n", filepath)
				fmt.Printf("export const get%s = (item) => {\n", message)
				fmt.Println("\tif (item) return {")
				found = true
			} else {
				if found {
					if line == "}" {
						fmt.Println("\t}")
						fmt.Println("}")
						break
					} else {
						if strings.Contains(line, "`protobuf") {
							processLineGet(line)
						}
					}
				}
			}
		}
	}
	found = false
	if !strings.HasSuffix(message, "Reply") {
		for _, line := range lines {
			if line == typeMessageStruct {
				fmt.Printf("// %s\n", filepath)
				fmt.Printf("export const set%s = (item) => {\n", message)
				fmt.Printf("\tif (item) return new %s()\n", message)
				found = true
			} else {
				if found {
					if line == "}" {
						fmt.Println("}")
						break
					} else {
						if strings.Contains(line, "`protobuf") {
							processLineSet(line)
						}
					}
				}
			}
		}
	}
	_ = found
	return nil
}

func getType(lexemes []string) string {
	for i, lexeme := range lexemes {
		if i > 0 && lexeme != "" {
			return lexeme
		}
	}
	return ""
}

func processLineGet(line string) {
	lexemes := strings.Split(strings.Trim(line, "\t"), " ")
	name := lexemes[0]
	name_ := strcase.LowerCamelCase(name)
	typ := getType(lexemes)
	_typ := strings.Trim(typ, "[]*")
	fmt.Print("\t\t")
	switch {
	case strings.HasPrefix(typ, "[]"):
		fmt.Printf("%s: listToArray(item.get%sList(), get%s),\n", name_, name, _typ)
	case strings.HasPrefix(typ, "*common."):
		fmt.Printf("%s: get%s(item.get%s()), // import {get%s} from '@/assets/json/common'\n", name_, typ[8:], name, typ[8:])
	case strings.HasPrefix(typ, "*"):
		fmt.Printf("%s: get%s(item.get%s()),\n", name_, _typ, name)
	default:
		fmt.Printf("%s: item.get%s(),\n", name_, name)
	}
}

func processLineSet(line string) {
	lexemes := strings.Split(strings.Trim(line, "\t"), " ")
	name := lexemes[0]
	//name_ := strcase.LowerCamelCase(name)
	typ := getType(lexemes)
	fmt.Print("\t\t.set")
	switch {
	case strings.HasPrefix(typ, "[]"):
		fmt.Printf("%sList(item)\n", name) //TODO
	case strings.HasPrefix(typ, "*common."):
		fmt.Printf("%s(set%s(item))\n", name, typ[8:])
	case strings.HasPrefix(typ, "*"):
		fmt.Printf("%s(set%s(item))\n", name, strings.Trim(typ, "*"))
	default:
		fmt.Printf("%s(item)\n", name)
	}
}
