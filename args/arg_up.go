package args

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/timurkash/mcsdeploy/utils/commands"
	"github.com/timurkash/mcsdeploy/utils/settings"
	"os"
	"strconv"
)

var badTagError = errors.New("bad tag")

func ArgUp(level int, serviceName string) error {

	dcServices := &settings.DCServices{}
	if err := dcServices.Load(); err != nil {
		return err
	}
	service, ok := dcServices.Services[serviceName]
	if !ok {
		fmt.Errorf("service %s not found", serviceName)
	}
	os.Chdir(service.Build.Context)
	status, err := commands.Exec("git", "status")
	if err != nil {
		return err
	}
	if !bytes.Contains(status, []byte("nothing to commit")) {
		errors.New("not committed")
	}
	tag, err := commands.Exec("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return err
	}
	tag = bytes.Trim(tag, "\n")
	fmt.Printf("version is %s", tag)
	if tag[0] != 'v' {
		return errors.New(" bad tag")
	}
	parts := bytes.Split(tag[1:], []byte("."))
	if len(parts) != 3 {
		return badTagError
	}
	major, err := getInt(parts[0])
	if err != nil {
		return err
	}
	minor, err := getInt(parts[1])
	if err != nil {
		return err
	}
	patch, err := getInt(parts[2])
	if err != nil {
		return err
	}
	switch level {
	case 0:
		major++
		minor = 0
		patch = 0
	case 1:
		minor++
		patch = 0
	case 2:
		patch++
	}
	versionNext := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	_, err = commands.Exec("git", "tag", "-a", versionNext, "-F", "tag_message")
	if err != nil {
		return err
	}

	//valuesBytes, err := os.ReadFile("values.yaml")
	//if err != nil {
	//	return err
	//}
	//valuesBytes = bytes.ReplaceAll(valuesBytes,
	//	[]byte(fmt.Sprintf("  tag: %s", string(version))),
	//	[]byte(fmt.Sprintf("  tag: %s", versionNext)))
	//if err := os.WriteFile("values.yaml", valuesBytes, 0644); err != nil {
	//	return err
	//}
	//if err := os.WriteFile("tag", []byte(versionNext), 0644); err != nil {
	//	return err
	//}
	return nil
}

func getInt(part []byte) (int, error) {
	return strconv.Atoi(string(part))
}
