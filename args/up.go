package args

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/timurkash/mcsdeploy/utils/commands"
	"github.com/timurkash/mcsdeploy/utils/settings"
	"os"
	"strconv"
	"strings"
)

func ArgUp(level int, serviceName string) error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	dcServices := &settings.DCServices{}
	if err := dcServices.Load(); err != nil {
		return err
	}
	service, ok := dcServices.Services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found", serviceName)
	}
	if service.Build == nil {
		return fmt.Errorf("no build section")
	}
	if err := os.Chdir(service.Build.Context); err != nil {
		return err
	}
	status, err := commands.Exec("git", "status")
	if err != nil {
		return err
	}
	if !bytes.Contains(status, []byte("nothing to commit")) {
		return fmt.Errorf("%s not committed", serviceName)
	}
	tag, err := commands.Exec("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return err
	}
	tag = bytes.Trim(tag, "\n")
	if tag[0] != 'v' {
		return errors.New("version tag must begin with v")
	}
	parts := bytes.Split(tag[1:], []byte("."))
	if len(parts) != 3 {
		return fmt.Errorf("bad tag")
	}
	fmt.Printf("version tag was %s\n", tag)
	imageTag := getTag(service.Image)
	fmt.Printf("  image tag was %s\n", imageTag)
	if string(tag) != imageTag {
		return errors.New("tags not according")
	}
	major := 0
	minor := 0
	patch := 0
	if level == 0 {
		major, err = getInt(parts[0])
		if err != nil {
			return err
		}
	}
	if level <= 1 {
		minor, err = getInt(parts[1])
		if err != nil {
			return err
		}
	}
	if level <= 2 {
		patch, err = getInt(parts[2])
		if err != nil {
			return err
		}
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
	_, err = commands.Exec("git", "tag", "-a", versionNext,
		"-m", fmt.Sprintf("version tag up to %s", versionNext))
	if err != nil {
		return err
	}
	if err := os.Chdir(curDir); err != nil {
		return err
	}
	newImageTag := strings.ReplaceAll(service.Image, string(tag), versionNext)
	newDat := bytes.ReplaceAll(dcServices.Dat, []byte(service.Image), []byte(newImageTag))
	return dcServices.Save(newDat)
}

func getInt(part []byte) (int, error) {
	return strconv.Atoi(string(part))
}

func getTag(image string) string {
	p := strings.LastIndex(image, ":")
	return image[p+1:]
}
