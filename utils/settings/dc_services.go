package settings

import (
	"fmt"
	"github.com/timurkash/mcsdeploy/utils"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	DCServices struct {
		Services map[string]*DCService `yaml:"services"`
		Dat      []byte
	}
	DCService struct {
		Image string `yaml:"image"`
		Build *Build `yaml:"build"`
	}
	Build struct {
		Context    string `yaml:"context"`
		Dockerfile string `yaml:"dockerfile"`
	}
)

const DockerCompose = "docker-compose.yml"

func (s *DCServices) Load() error {
	if !utils.IsFileExists(DockerCompose) {
		return fmt.Errorf("filename %s not exists", DockerCompose)
	}
	var err error
	s.Dat, err = os.ReadFile(DockerCompose)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(s.Dat, s)
}

func (s *DCServices) Save(dat []byte) error {
	return os.WriteFile(DockerCompose, dat, 0644)
}
