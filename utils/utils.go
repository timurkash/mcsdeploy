package utils

import (
	"log"
	"os"
	"text/template"

	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	Single      string
	SingleLower string
	Plural      string
	PluralLower string
}

var cas = cases.Title(language.English)

func GetSinglePlural(argStrings []string) (singlePlural *SinglePlural) {
	if len(argStrings) < 2 {
		return nil
	}
	var single, plural string
	if len(argStrings) >= 3 {
		single = argStrings[2]
	}
	if len(argStrings) >= 4 {
		plural = argStrings[3]
	}
	if plural == "" {
		plural = single + "s"
	}
	return &SinglePlural{
		Single:      strcase.UpperCamelCase(cas.String(single)),
		SingleLower: single,
		Plural:      strcase.UpperCamelCase(cas.String(plural)),
		PluralLower: plural,
	}
}

func StdOut(template *template.Template, sp *SinglePlural) error {
	return template.Execute(os.Stdout, sp)
}
