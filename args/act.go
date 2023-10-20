package args

import (
	"fmt"
	"strings"
)

func ArgActRequest(message string) error {
	if strings.HasPrefix(message, "Act") && strings.HasSuffix(message, "Request") {
		entity := message[3 : len(message)-7]
		fmt.Printf("\texport function getAct%sRequest({action, id, item}) {\n", entity)
		fmt.Printf("\t\tconst request = new Act%sRequest().setActionId(getActionId({action, id, item}))\n", entity)
		fmt.Printf("\t\tif (action === INSERT || action === UPDATE) {\n")
		fmt.Printf("\t\t\trequest.set%s(set%sInfo(item))\n", entity, entity)
		fmt.Println("\t\t}")
		fmt.Println("\t\treturn request")
		fmt.Println("\t}")
	}
	return nil
}

//func getActRequest(message string) error {
//	//
//	//lines := strings.Split(string(bytes), "\n")
//	//found := false
//	//typeMessageStruct := fmt.Sprintf("type %s struct {", message)
//	//for _, line := range lines {
//	//	if line == typeMessageStruct {
//	//		fmt.Printf("// %s\n", filepath)
//	//		fmt.Printf("import {Act%sRequest, %sInfo} from \"@/assets/proto/js/api/%s/messages_pb\"\n", entity, entity, strings.ToLower(entity))
//	//		fmt.Println()
//	//		fmt.Printf("export function getAct%sRequest({action, id, item}) {\n", entity)
//	//		fmt.Printf("\tconst request = new Act%sRequest().setActionId(getActionId({action, id, item}))\n", entity)
//	//		fmt.Printf("\tif (action === INSERT || action === UPDATE) {\n")
//	//		fmt.Printf("\t\trequest.set%s(new %s()\n", entity, message)
//	//		found = true
//	//	} else {
//	//		if found {
//	//			if line == "}" {
//	//				fmt.Println("\t\t)")
//	//				fmt.Println("\t}")
//	//				fmt.Println("\treturn request")
//	//				fmt.Println("}")
//	//				return nil
//	//			} else {
//	//				if strings.Contains(line, "`protobuf") {
//	//					processActLine(line)
//	//				}
//	//			}
//	//		}
//	//	}
//	//}
//	//_ = found
//	return nil
//}

//func findNotInfoMessage(filepath, message string) error {
//	bytes, err := os.ReadFile(filepath)
//	if err != nil {
//		return err
//	}
//	lines := strings.Split(string(bytes), "\n")
//	found := false
//	typeMessageStruct := fmt.Sprintf("type %s struct {", message)
//	for _, line := range lines {
//		if line == typeMessageStruct {
//			fmt.Printf("// %s\n", filepath)
//			fmt.Printf("export function get%s(item) {\n", message)
//			fmt.Printf("\tif (item) {\n")
//			fmt.Printf("\t\treturn new %s()\n", message)
//			//fmt.Printf("\tconst request = new Act%sRequest().setActionId(getActionId({action, id, item}))\n", entity)
//			//fmt.Printf("\tif (action === INSERT || action === UPDATE) {\n")
//			//fmt.Printf("\t\trequest.set%s(new %s()\n", entity, message)
//			found = true
//		} else {
//			if found {
//				if line == "}" {
//					fmt.Println("\t}")
//					fmt.Println("}")
//					return nil
//				} else {
//					if strings.Contains(line, "`protobuf") {
//						processActLine(line)
//					}
//				}
//			}
//		}
//	}
//	_ = found
//	return nil
//}

//func processActLine(line string) {
//	line = strings.Trim(line, "\t")
//	lexemes := strings.Split(line, " ")
//	name := lexemes[0]
//	name_ := strcase.LowerCamelCase(name)
//	fmt.Print("\t\t\t")
//	fmt.Printf(".set%s(item.%s)\n", name, name_)
//}
