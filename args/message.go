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
				fmt.Printf("export const get%s = item => getIfItem(item, item => ({\n", message)
				found = true
			} else {
				if found {
					if line == "}" {
						fmt.Println("}))")
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
		for i, line := range lines {
			if line == typeMessageStruct {
				fmt.Printf("// %s\n", filepath)
				fmt.Printf("const set%s = item => getIfItem(item, item => new %s()\n", message, message)
				found = true
			} else {
				if found {
					if line == "}" {
						fmt.Println("\r)")
						break
					} else {
						if strings.Contains(line, "`protobuf") {
							prev := lines[i-1]
							if !(strings.HasPrefix(prev, "\t// ro:") || strings.HasPrefix(prev, "\t//ro:")) {
								processLineSet(line)
							}
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
	fmt.Print("\t")
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
	name_ := strcase.LowerCamelCase(name)
	typ := getType(lexemes)
	switch {
	case strings.HasPrefix(typ, "[]"):
		//fmt.Printf("%sList(item)\n", name) //TODO
	case strings.HasPrefix(typ, "*common."):
		//fmt.Printf("%s(set%s(item))\n", name, typ[8:])
	case strings.HasPrefix(typ, "*"):
		//fmt.Printf("%s(set%s(item))\n", name, strings.Trim(typ, "*"))
	default:
		fmt.Printf("\t.set%s(item.%s)\n", name, name_)
	}
}
