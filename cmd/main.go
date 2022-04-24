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
	// load config
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	deeplCfg, err := cfg.GetSection("deepl")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	flag.Parse()
	jsonCmd := flag.NewFlagSet("json", flag.ExitOnError)

	switch os.Args[1] {
	case "json":
		jsonCmd.Parse(os.Args[2:])
		j, err := ddl.NewJsonObj([]byte(jsonCmd.Arg(0)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = j.Translate(jsonCmd.Args()[1:], translate.NewDeepl(deeplCfg.Key("key").String()))
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
		fmt.Fprintln(os.Stderr, "Invalid Subcommand")
		os.Exit(1)
	}
}
