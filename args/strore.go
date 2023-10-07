package args

import (
	"fmt"
	"os"
	"strings"
)

func ArgStore(service string) error {
	files, err := os.ReadDir("api")
	if err != nil {
		return err
	}
	found := false
	for _, file := range files {
		if file.IsDir() && file.Name() == service {
			found = true
			if err := findRpc(fmt.Sprintf("api/%s/%s.proto", service, service)); err != nil {
				return err
			}
		}
	}
	if !found {
		fmt.Printf("service %s not found", service)
	}
	return nil
}

func findRpc(filepath string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	fmt.Println("actions: {")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "rpc ") {
			if err := processRpcLine(line); err != nil {
				return err
			}
		}
	}
	fmt.Println("},")
	return nil
}

func processRpcLine(line string) error {
	lexemes := clearSlice(strings.Split(line, " "))
	rpc := lexemes[1]
	if len(lexemes) < 5 {
		return fmt.Errorf("bad rpc line %s", rpc)
	}
	rpc_ := strings.ToLower(rpc[:1]) + rpc[1:]
	if lexemes[3] != "returns" {
		return fmt.Errorf("bad rpc line %s", rpc)
	}
	request := clearMessage(lexemes[2])
	param := "item"
	if strings.HasPrefix(request, "Act") {
		param = "actionIdItem"
	}
	reply := clearMessage(lexemes[4])
	fmt.Printf("    async %s(%s) {\n", rpc_, param)
	fmt.Printf("        try {\n")
	fmt.Printf("            const reply = await client.%s(get%s(%s), await getMetadata())\n", rpc_, request, param)
	fmt.Printf("            return get%s(reply)\n", reply)
	fmt.Printf("        } catch (err) {\n")
	fmt.Printf("            console.error(err)\n")
	fmt.Printf("        }\n")
	fmt.Printf("    },\n")
	return nil
}

func clearMessage(str string) string {
	str = strings.Trim(str, "(){;")
	pos := strings.Index(str, ".")
	if pos == -1 {
		return str
	}
	return str[pos+1:]
}
