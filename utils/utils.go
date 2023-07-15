package utils

import (
	"bufio"
	"fmt"
	"github.com/stoewer/go-strcase"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Util struct {
	Name    string
	Command string
}

func IsFileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil {
		if !fi.IsDir() {
			return true
		} else {
			log.Println(filename, "is directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type SinglePlural struct {
	Single           string
	SingleLower      string
	SingleLowerLower string
	Lower            string
	Plural           string
	PluralLower      string
	SnakeLower       string
	SnakePlural      string
	Service          string
	ServiceLower     string
}

var cas = cases.Title(language.English)

func GetSinglePlural(argStrings []string) (singlePlural *SinglePlural) {
	if len(argStrings) < 3 {
		return nil
	}
	var single string
	if len(argStrings) >= 3 {
		single = argStrings[2]
	}
	service := single
	if strings.Contains(single, ":") {
		split := strings.Split(single, ":")
		single = split[0]
		service = split[1]
	}
	plural := getPlural(single)
	title := func(str string) string {
		return cas.String(str[:1]) + str[1:]
	}
	lower := func(str string) string {
		return strings.ToLower(str[:1]) + str[1:]
	}
	return &SinglePlural{
		Single:           title(single),
		SingleLower:      lower(single),
		SingleLowerLower: strings.ToLower(single),
		Lower:            strings.ToLower(single),
		SnakeLower:       strcase.SnakeCase(single),
		SnakePlural:      strcase.SnakeCase(plural),
		Plural:           title(plural),
		PluralLower:      lower(plural),
		Service:          service,
		ServiceLower:     lower(service),
	}
}

func getPlural(single string) string {
	const irrPluralsFilename = "irr_plurals.yaml"
	plural := fmt.Sprintf("%ss", single)
	if !IsFileExists(irrPluralsFilename) {
		return plural
	}
	names := make(map[string]string)
	file, err := os.Open(irrPluralsFilename)
	if err != nil {
		return plural
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(names); err != nil {
		return plural
	}
	if pluralFound, ok := names[single]; ok {
		return pluralFound
	}
	return plural
}

func StdOut(template *template.Template, sp *SinglePlural) error {
	return template.Execute(os.Stdout, sp)
}

//type IrrPlurals struct {
//	Names map[string]string `yaml:"names"`
//}
//
//const irrPluralsFilename = "irr_plurals.yaml"
//
//func (ip *IrrPlurals) Load() error {
//	ip.Names = make(map[string]string)
//	if !utils.IsFileExists(irrPluralsFilename) {
//		return nil
//	}
//	file, err := os.Open(irrPluralsFilename)
//	if err != nil {
//		return err
//	}
//	return yaml.NewDecoder(bufio.NewReader(file)).Decode(ip)
//}
