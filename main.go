package main

import (
	"errors"
	"github.com/timurkash/mcsdeploy/args"
	"github.com/timurkash/mcsdeploy/utils"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetPrefix("[>error<] ")
	log.SetFlags(0)
	argStrings := os.Args
	if len(argStrings) == 1 {
		args.ShowDescription()
		return
	}
	arg := argStrings[1]
	if len(arg) != 4 {
		log.Fatalln("argument must be 4 characters")
	}
	if arg[0] != '-' {
		log.Fatalln("argument must begin with dash")
	}
	if len(argStrings) == 1 {
		args.ShowDescription()
	}
	var service string
	if len(argStrings) == 3 {
		service = argStrings[2]
	}
	pwd, err := getPwd()
	if err != nil {
		log.Fatalln(err)
	}
	if pwd != "proto" {
		log.Fatalln("you are not in proto")
	}
	switch arg {
	case "-upv":
		if err := args.ArgUp(2, service); err != nil {
			log.Fatalln(err)
		}
	case "-uvp":
		if err := args.ArgUp(1, service); err != nil {
			log.Fatalln(err)
		}
	case "-vup":
		if err := args.ArgUp(0, service); err != nil {
			log.Fatalln(err)
		}
	case "-env":
		if err := args.ArgEnvoy(); err != nil {
			log.Fatalln(err)
		}
	case "-doc":
		if err := args.ArgDocker(); err != nil {
			log.Fatalln(err)
		}
	case "-mak":
		if err := args.ArgMake(); err != nil {
			log.Fatalln(err)
		}
	case "-rep":
		sp := utils.GetSinglePlural(argStrings)
		if sp == nil {
			args.ShowDescription()
			return
		}
		if err := utils.StdOut(args.RepoTemp, sp); err != nil {
			log.Fatalln(err)
		}
	case "-prt":
		sp := utils.GetSinglePlural(argStrings)
		if sp == nil {
			args.ShowDescription()
			return
		}
		if err := utils.StdOut(args.ProtoTemp, sp); err != nil {
			log.Fatalln(err)
		}
	case "-sql":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgSql(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	case "-msg":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgMessage(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	case "-req":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgActRequest(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	case "-enm":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgEnum(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	case "-srv":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgService(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	case "-str":
		if len(argStrings) <= 2 {
			args.ShowDescription()
			return
		}
		if err := args.ArgStore(argStrings[2]); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalf("option %s not defined\n", arg)
	}
}

func getPwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := strings.LastIndex(pwd, "/")
	if p == -1 {
		return "", errors.New("bad dir")
	}
	return pwd[p+1:], nil
}
