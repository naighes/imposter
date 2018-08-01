package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/naighes/imposter/functions"
)

func validateCmd() command {
	fs := flag.NewFlagSet("imposter validate", flag.ExitOnError)
	opts := validateOpts{}
	fs.StringVar(&opts.configFile, "config-file", "stdin", "The configuration file")
	fs.BoolVar(&opts.jsonEncoded, "json", false, "Enable JSON output instead of plain text")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return validateExec(&opts)
	}}
}

type validateOpts struct {
	configFile  string
	jsonEncoded bool
}

func validateExec(opts *validateOpts) error {
	var r []string
	config, err := readConfig(opts.configFile)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v", err)
	}
	var vars map[string]interface{}
	if config.Vars == nil {
		vars = make(map[string]interface{})
	} else {
		vars = config.Vars
	}
	defs := config.Defs
	for _, def := range defs {
		errors := def.validate(functions.ParseExpression, vars)
		if len(errors) > 0 {
			r = append(r, errors...)
		}
	}
	if l := len(r); l > 0 {
		if opts.jsonEncoded {
			rep := errorReport{Errors: r, Count: l}
			bytes, _ := json.MarshalIndent(&rep, "", "  ")
			fmt.Printf("%s", string(bytes))
		} else {
			const sep = "\n--------------------\n"
			fmt.Printf("found %d errors:%s", l, sep)
			fmt.Printf(strings.Join(r[:], sep))
		}
		os.Exit(1)
	}
	return nil
}

type errorReport struct {
	Count  int      `json:"count"`
	Errors []string `json:"errors"`
}
