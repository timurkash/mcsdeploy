package settings

import (
	"bufio"
	"fmt"
	"github.com/timurkash/mcsdeploy/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	DCServices struct {
		Services map[string]DCService `yaml:"services"`
	}
	DCService struct {
		Image string `yaml:"image"`
		Build Build  `yaml:"build"`
	}
	Build struct {
		Context    string `yaml:"context"`
		Dockerfile string `yaml:"dockerfile"`
	}
)

const (
	DockerCompose = "docker-compose.yml"
)

func (s *DCServices) Load() error {
	if !utils.IsFileExists(DockerCompose) {
		return fmt.Errorf("filename %s not exists", DockerCompose)
	}
	file, err := os.Open(DockerCompose)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(s); err != nil {
		return err
	}
	return nil
}
