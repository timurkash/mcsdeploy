package args

import (
	"fmt"
	settingsPackage "github.com/timurkash/mcsdeploy/utils/settings"
)

func ArgMake() error {
	services := &settingsPackage.Services{}
	if err := services.Load(); err != nil {
		return err
	}
	fmt.Println("\n---")
	for _, service := range services.Services {
		fmt.Printf("      - %s\n", service.NameVersion)
	}
	fmt.Println("\n--- make docker-compose")
	for _, service := range services.Services {
		fmt.Println()
		fmt.Printf("dc-%s:\n", service.Name)
		fmt.Printf("\tdocker-compose build %s\n", service.NameVersion)
		fmt.Printf("\tdocker-compose stop %s\n", service.NameVersion)
		fmt.Printf("\tdocker-compose up -d %s\n", service.NameVersion)
	}
	fmt.Println("\n# make pull-all")
	fmt.Println("pull-all:")
	for _, service := range services.Services {
		fmt.Printf("\t@echo %s && git -C ../../back/%s pull || git clone https://%s/back/%s.git ../../back\n",
			service.NameVersion,
			service.NameVersion,
			services.ProjectRepo,
			service.NameVersion,
		)
	}
	return nil
}
