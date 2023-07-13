package args

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	settingsPackage "github.com/timurkash/mcsdeploy/utils/settings"
	"strings"
)

type Field struct {
	Name    string
	Type    string
	Camel   string
	Camel_  string
	Counter uint16
}

func ArgSql(fieldsTable string) error {
	fieldTypes := &settingsPackage.Fields{}
	if err := fieldTypes.Load(); err != nil {
		return err
	}
	split := strings.Split(fieldsTable, ":")
	if len(split) != 2 {
		return fmt.Errorf("wrong arg %s", fieldsTable)
	}
	fieldsString := split[0]
	table := split[1]
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
		typ, ok := fieldTypes.Fields[field]
		if !ok {
			return fmt.Errorf("field %s not declared in fieldTypes", field)
		}
		fields = append(fields, Field{
			Name:    field,
			Type:    typ,
			Camel:   strcase.UpperCamelCase(field),
			Camel_:  strcase.LowerCamelCase(field),
			Counter: uint16(counter + 1),
		})
	}
	ucc := strcase.UpperCamelCase(table)
	ucc_ := strcase.LowerCamelCase(table)
	fmt.Printf(`

--- proto

mcsdeploy -rep %s | pbcopy

mcsdeploy -prt %s | pbcopy

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
		case "string", "uint32", "bool", "uint64", "float", "double":
			entType = fmt.Sprintf("%s%s", strings.ToUpper(field.Type[:1]), field.Type[1:])
		//case "float64":
		//	entType = "Double"
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
            %s.%s(name),`, ucc_, field.Camel)
		}
	}
	fmt.Printf(`

// --- pinia actions

        async list%s() {
            const metadata = await getMetadata()
            if (!metadata) {
                this.products = Array()
                return
            }
            try {
                const the%s = Array()
                const reply = await client.list%s(new List%sRequest(), metadata)
                reply.getItemsList().forEach(el => the%s.push(this.get%sItem(el)))
                this.%s = the%s
            } catch (err) {
                console.error(err)
            }
        },
`, ucc, ucc, ucc, ucc, ucc, ucc, ucc_, ucc)
	fmt.Printf(`
        get%sItem(el) {
            const item = el.getItem()
            return {
              idTimestamps: getIdTimestamp(el.getIdTimestamps()),`, ucc)
	for _, field := range fields {
		fmt.Printf(`
              %s: item.get%s(),`, field.Camel_, field.Camel)
	}
	fmt.Println(`
            }
        },`)
	fmt.Printf(`
                .set%s(new %sInfo()`, ucc, ucc)
	for _, field := range fields {
		fmt.Printf(`
                    .set%s(%s.value.%s)`, field.Camel, ucc_, field.Camel_)
	}
	fmt.Println(`
                )`)
	fmt.Printf(`
  formData.value = {`)
	for _, field := range fields {
		fmt.Printf(`
    %s: data.value.%s,`, field.Camel_, field.Camel_)
	}
	fmt.Printf(`
  }`)
	return nil
}
