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

type IdValue struct {
	Id    string
	Value string
}

func findEnum(filepath, enum string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	found := false
	typeMessageStruct := fmt.Sprintf("type %s int32", enum)
	var idValues []IdValue
	for _, line := range lines {
		if line == typeMessageStruct {
			fmt.Printf("// %s\n", filepath)
			fmt.Printf("export const get%sString = id => {\n", enum)
			fmt.Printf("\tswitch (id) {\n")
			found = true
		} else {
			if found {
				if line == ")" {
					fmt.Println("\t}")
					fmt.Println("}")
					break
				} else {
					if strings.Contains(line, enum+"_") {
						idValues = append(idValues, processEnumLine(enum, line))
					}
				}
			}
		}
	}
	if found {
		fmt.Printf("export const %sArray = [\n", toLowerFirst(enum))
		for _, idValue := range idValues {
			fmt.Printf("\t{id: %s, name: \"%s\"},\n", idValue.Id, idValue.Value)
		}
		fmt.Println("]")
	}
	return nil
}

func processEnumLine(enum, line string) IdValue {
	line = strings.Trim(line, "\t")
	lexemesClear := clearSlice(strings.Split(line, " "))
	id := lexemesClear[2]
	value := lexemesClear[0][len(enum)+1:]
	idValue := IdValue{Id: id, Value: value}
	fmt.Printf("\t\tcase %s:\n\t\t\treturn \"%s\"\n", lexemesClear[2], value)
	return idValue
}
