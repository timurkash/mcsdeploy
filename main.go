package main

import (
	"github.com/timurkash/mcsdeploy/args"
	"log"
	"os"
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
	case "-prt":
		var (
			single string
			plural string
		)
		if len(argStrings) >= 3 {
			single = argStrings[2]
		}
		if len(argStrings) >= 4 {
			plural = argStrings[3]
		}
		if single == "" {
			args.ShowDescription()
			return
		}
		if err := args.ArgProto(single, plural); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalf("option %s not defined\n", arg)
	}
}
