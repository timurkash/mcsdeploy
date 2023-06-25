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
		if !(v == "string" || v == "uint32") { // TODO
			return fmt.Errorf("wrong type %s for fileld %s", v, k)
		}
	}
	return nil
}
