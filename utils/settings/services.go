package settings

import (
	"bufio"
	"fmt"
	"github.com/timurkash/mcsdeploy/utils"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ServicesFilename = "services.yaml"
)

type (
	Service struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version"`
		Port     int16  `yaml:"port"`
		HttpPort int16  `yaml:"httpPort"`
	}
	Services struct {
		Services []Service `yaml:"services"`
	}
)

func (s *Services) Load() error {
	if !utils.IsFileExists(ServicesFilename) {
		return fmt.Errorf("filename %s not exists", ServicesFilename)
	}
	file, err := os.Open(ServicesFilename)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(s); err != nil {
		return err
	}
	return nil
}
