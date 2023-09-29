package args

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ArgEnum(enum string) error {
	if err := filepath.Walk("gen/go", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.HasSuffix(path, "/messages.pb.go") || strings.HasSuffix(path, "/common.pb.go")) {
			if err := findEnum(path, enum); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func findEnum(filepath, enum string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	found := false
	typeMessageStruct := fmt.Sprintf("type %s int32", enum)
	for _, line := range lines {
		if line == typeMessageStruct {
			fmt.Println()
			fmt.Printf("// %s\n", filepath)
			fmt.Println()
			fmt.Printf("export function get%sString(id) {\n", enum)
			fmt.Printf("\tswitch (id) {\n")
			found = true
		} else {
			if found {
				if line == ")" {
					fmt.Println("\t}")
					fmt.Println("}")
					return nil
				} else {
					if strings.Contains(line, enum+"_") {
						processEnumLine(line)
					}
				}
			}
		}
	}
	_ = found
	return nil
}

func processEnumLine(line string) {
	line = strings.Trim(line, "\t")
	lexemes := strings.Split(line, " ")
	lexemesClear := make([]string, 0, 3)
	for _, lex := range lexemes {
		if lex != "" && lex != "=" {
			lexemesClear = append(lexemesClear, lex)
		}
	}
	values := strings.Split(lexemesClear[0], "_")
	//value := values[1]
	//name := lexemes[0]
	fmt.Printf("\t\tcase %s:\n", lexemesClear[2])
	fmt.Printf("\t\t\treturn \"%s\"\n", values[1])
}
