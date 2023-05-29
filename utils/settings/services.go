package settings

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/timurkash/mcsdeploy/utils"

	"gopkg.in/yaml.v3"
)

const (
	ServicesFilename = "services.yaml"
)

type (
	Service struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		NameVersion string
		Port        int16 `yaml:"port"`
		HttpPort    int16 `yaml:"httpPort"`
	}
	Services struct {
		Services    []Service `yaml:"services"`
		ProjectRepo string
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
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	goPath := os.Getenv("GOPATH")
	projectRepo := strings.ReplaceAll(wd, fmt.Sprintf("%s/src/", goPath), "")
	projectRepo = strings.ReplaceAll(projectRepo, "/deploy/local", "")
	for i, service := range s.Services {
		if s.Services[i].Version == "" {
			s.Services[i].Version = "v1"
			service.Version = "v1"
		}
		s.Services[i].NameVersion = fmt.Sprintf("%s-%s", service.Name, service.Version)
	}
	s.ProjectRepo = projectRepo
	return nil
}
