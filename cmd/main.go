package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nikkely/ddl-translate/pkg/ddl"
	"github.com/nikkely/ddl-translate/pkg/translate"
	"gopkg.in/ini.v1"
)

func main() {
	cfgPath := flag.String("config", "./config.ini", "/path/to/config.ini")
	flag.Parse()
	jsonCmd := flag.NewFlagSet("json", flag.ExitOnError)

	// load config
	cfg, err := ini.Load(*cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	deeplCfg, err := cfg.GetSection("deepl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// subcommands
	// HACK: no need to loop
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[1+i] {
		case "json":
			jsonPath := jsonCmd.String("path", "", "path to json")
			jsonCmd.Parse(os.Args[2+i:])

			var j *ddl.JSONObj
			pos := 0
			if *jsonPath == "" {
				j, err = ddl.NewJSONObj([]byte(jsonCmd.Arg(0)))
				pos++
			} else {
				j, err = ddl.NewJSONObjFromFile(*jsonPath)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}

			err = j.Translate(jsonCmd.Args()[pos:], translate.NewDeepl(deeplCfg))
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			result, err := j.ToString()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Fprintln(os.Stdout, result)
			os.Exit(0)
		default:
			continue
		}
		fmt.Fprintln(os.Stderr, "Invalid Subcommand")
		os.Exit(1)
	}
}
