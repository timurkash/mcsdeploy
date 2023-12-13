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

func processLineGet(line string) {
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
	fmt.Print("\t\t")
	switch {
	case strings.HasPrefix(typ, "[]"):
		fmt.Printf("%s: item.get%sList(),\n", name_, name)
	case strings.HasPrefix(typ, "*common."):
		fmt.Printf("%s: get%s(item.get%s()), // import {get%s} from '@/assets/json/common'\n", name_, typ[8:], name, typ[8:])
	case strings.HasPrefix(typ, "*"):
		//ix := strings.Index(typ, ".")
		//if ix == -1 {
		//	typ = typ[ix:]
		//} else {
		//	typ = typ[1:]
		//}
		fmt.Printf("%s: get%s(item.get%s()),\n", name_, name, name)
	default:
		fmt.Printf("%s: item.get%s(),\n", name_, name)
	}
}

func processLineSet(line string) {
	line = strings.Trim(line, "\t")
	lexemes := strings.Split(line, " ")
	name := lexemes[0]
	//name_ := strcase.LowerCamelCase(name)
	typ := ""
	for i, lexeme := range lexemes {
		if i > 0 && lexeme != "" {
			typ = lexeme
			break
		}
	}
	fmt.Print("\t\t.set")
	switch {
	case strings.HasPrefix(typ, "[]"):
		fmt.Printf("%sList(item)\n", name) //TODO
	case strings.HasPrefix(typ, "*common."):
		fmt.Printf("%s(set%s(item))\n", name, typ[8:])
	case strings.HasPrefix(typ, "*"):
		fmt.Printf("%s(set%s(item))\n", name, name)
	default:
		fmt.Printf("%s(item)\n", name)
	}
}
