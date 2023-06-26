package settings

import (
	"bufio"
	"fmt"
	"github.com/timurkash/mcsdeploy/utils"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	FieldsFilename = "fields.yaml"
)

var fields = map[string]struct{}{
	"string": {},
	"uint32": {},
	"uint64": {},
	"bool":   {},
	"int32":  {},
	"int64":  {},
	"float":  {},
	"double": {},
}

type (
	Fields struct {
		Fields map[string]string `yaml:"fields"`
	}
)

func (f *Fields) Load() error {
	if !utils.IsFileExists(FieldsFilename) {
		return fmt.Errorf("filename %s not exists", FieldsFilename)
	}
	file, err := os.Open(FieldsFilename)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(f); err != nil {
		return err
	}
	for k, v := range f.Fields {
		if _, ok := fields[v]; !ok {
			return fmt.Errorf("wrong type %s for fileld %s", v, k)
		}
	}
	return nil
}
