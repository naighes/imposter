package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func validateCmd() command {
	fs := flag.NewFlagSet("imposter validate", flag.ExitOnError)
	opts := validateOpts{}
	fs.StringVar(&opts.configFile, "config-file", "stdin", "The configuration file")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return validateExec(&opts)
	}}
}

type validateOpts struct {
	configFile string
}

func validateExec(opts *validateOpts) error {
	var r []string
	config, err := readConfig(opts.configFile)
	if err != nil {
		return fmt.Errorf("could not load configuration: %v\n", err)
	}
	var vars map[string]interface{}
	if config.Vars == nil {
		vars = make(map[string]interface{})
	} else {
		vars = config.Vars
	}
	defs := config.Defs
	for _, def := range defs {
		e, err := ParseExpression(def.RuleExpression)
		if err != nil {
			r = append(r, fmt.Sprintf("%v", err))
		} else {
			a, err := e.evaluate(vars, &http.Request{Header: http.Header{}})
			if err != nil {
				r = append(r, fmt.Sprintf("%v", err))
			}
			_, ok := a.(bool)
			if !ok {
				r = append(r, fmt.Sprintf("evaluation error: expected 'bool'; got '%v' instead", reflect.TypeOf(a)))
			}
		}
	}
	if l := len(r); l > 0 {
		const sep = "\n--------------------\n"
		fmt.Printf("found %d errors%s:", l, sep)
		fmt.Printf(strings.Join(r[:], sep))
		os.Exit(1)
	}
	return nil
}
