package args

import (
	"bufio"
	"fmt"
	"github.com/stoewer/go-strcase"
	"github.com/timurkash/mcsdeploy/utils"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Field struct {
	Name       string
	Type       string
	Camel      string
	CamelLower string
	Counter    uint16
}

var fieldTypes = map[string]struct{}{
	"string": {},
	"uint32": {},
	"uint64": {},
	"bool":   {},
	"int32":  {},
	"int64":  {},
	"float":  {},
	"double": {},
}

const fieldsFilename = "fields.yaml"

func ArgSql(fieldsTable string) error {
	fieldsMap := make(map[string]string)
	if !utils.IsFileExists(fieldsFilename) {
		return fmt.Errorf("filename %s not exists", fieldsFilename)
	}
	file, err := os.Open(fieldsFilename)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(fieldsMap); err != nil {
		return err
	}
	for k, v := range fieldsMap {
		if _, ok := fieldTypes[v]; !ok {
			return fmt.Errorf("wrong type %s for fileld %s", v, k)
		}
	}
	split := strings.Split(fieldsTable, ":")
	if len(split) != 2 {
		return fmt.Errorf("wrong arg %s", fieldsTable)
	}
	fieldsString := split[0]
	table := split[1]
	pluralLower := utils.GetPlural(table)
	plural := utils.Title(pluralLower)
	fieldsSplit := strings.Split(fieldsString, ",")
	var sqlFieldsSplit []string
	for i := range fieldsSplit {
		sqlFieldsSplit = append(sqlFieldsSplit, fmt.Sprintf("  `%s`", fieldsSplit[i]))
	}
	fmt.Printf(`
--- sql

select 
%s
from %s`, strings.Join(sqlFieldsSplit, ", \n"), table)
	var fields []Field
	for counter, field := range fieldsSplit {
		typ, ok := fieldsMap[field]
		if !ok {
			return fmt.Errorf(`field "%s" not declared in fieldTypes`, field)
		}
		fields = append(fields, Field{
			Name:       field,
			Type:       typ,
			Camel:      strcase.UpperCamelCase(field),
			CamelLower: strcase.LowerCamelCase(field),
			Counter:    uint16(counter + 1),
		})
	}
	ucc := strcase.UpperCamelCase(table)
	ucc_ := strcase.LowerCamelCase(table)
	ucc__ := strings.ToLower(table)
	fmt.Printf(`

--- proto

mcsdeploy -prt %s | pbcopy

mcsdeploy -rep %s | pbcopy

message %sInfo {
`, ucc, ucc, ucc)
	for _, field := range fields {
		fmt.Printf("  %s %s = %d;\n", field.Type, field.Name, field.Counter)
	}
	fmt.Println("}")
	fmt.Printf(`

--- ent

ent init --target ./internal/data/ent/schema %s

`, ucc)
	for _, field := range fields {
		var entType string
		switch field.Type {
		case "string", "uint32", "bool", "uint64", "int32", "int64":
			entType = fmt.Sprintf("%s%s", strings.ToUpper(field.Type[:1]), field.Type[1:])
		case "double":
			entType = "Float"
		case "float":
			entType = "Float32"
		default:
			return fmt.Errorf("type %s not encountered", field.Type)
		}
		fmt.Printf(`        field.%s("%s"),
`, entType, field.Name)
	}
	fmt.Printf(`
// --- data.go

// get%sReply

        %s: &pb.%sInfo{
`, ucc, ucc, ucc)
	for _, field := range fields {
		fmt.Printf("            %s: record.%s,\n", field.Camel, field.Camel)
	}
	fmt.Println("        },")
	fmt.Printf(`
// set%s
`, ucc)
	for _, field := range fields {
		fmt.Printf(`
        Set%s(info.%s).`, field.Camel, field.Camel)
	}
	fmt.Println(`

// Or`)
	for _, field := range fields {
		if field.Type == "string" {
			fmt.Printf(`
            %s.%s(name),`, ucc__, field.Camel)
		}
	}
	fmt.Printf(`

// --- pinia actions

        async list%s() {
            const metadata = await getMetadata()
            if (!metadata) {
                throw notAuthorizedError
            }
            try {
                const the%s = Array()
                const reply = await client.list%s(new List%sRequest(), metadata)
                reply.get%sList().forEach(el => the%s.push(this.get%sItem(el)))
                return the%s
            } catch (err) {
                console.error(err)
            }
        },
`, plural, plural, plural, plural, plural, plural, ucc, plural)
	fmt.Printf(`        async act%s(action, %s) {
            const metadata = await getMetadata()
            if (!metadata) {
                throw notAuthorizedError
            }
            const request = new %sRequest()
                .setActionId(getActionId(action, %s))
                .set%s(new %sInfo()
                )
            try {
                const reply = await client.act%s(request, metadata)
                const the%s = [...this.%s]
                switch (action) {
                    case GET:
                        return this.get%sItem(reply)
                    case INSERT:
                        the%s.unshift(this.get%sItem(reply))
                        this.%s = the%s
                        return the%s
                    case UPDATE:
                        the%s[the%s.findIndex(el => el.idTimestamps.id === %s.id)] = this.get%sItem(reply)
                        this.%s = the%s
                        return the%s
                }
            } catch (err) {
                console.error(err)
            }
        },
`, ucc, ucc_, ucc, ucc_, ucc, ucc, ucc, plural, pluralLower, ucc, plural, ucc, pluralLower, plural, plural, plural,
		plural, ucc_, ucc, pluralLower, plural, plural)
	fmt.Printf(`
        get%sItem(el) {
            const item = el.get%s()
            return {
              idTimestamps: getIdTimestamp(el.getIdTimestamps()),`, ucc, ucc)
	for _, field := range fields {
		fmt.Printf(`
              %s: item.get%s(),`, field.CamelLower, field.Camel)
	}
	fmt.Println(`
            }
        },`)
	fmt.Printf(`
                .set%s(new %sInfo()`, ucc, ucc)
	for _, field := range fields {
		fmt.Printf(`
                    .set%s(%s.value.%s)`, field.Camel, ucc_, field.CamelLower)
	}
	fmt.Println(`
                )`)
	fmt.Printf(`
  formData.value = {`)
	for _, field := range fields {
		fmt.Printf(`
    %s: data.value.%s,`, field.CamelLower, field.CamelLower)
	}
	fmt.Printf(`
  }

`)
	return nil
}
