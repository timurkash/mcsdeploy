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
						processEnumLine(enum, line)
					}
				}
			}
		}
	}
	_ = found
	return nil
}

func processEnumLine(enum, line string) {
	line = strings.Trim(line, "\t")
	lexemesClear := clearSlice(strings.Split(line, " "))
	fmt.Printf("\t\tcase %s:\n\t\t\treturn \"%s\"\n", lexemesClear[2], lexemesClear[0][len(enum)+1:])
}
