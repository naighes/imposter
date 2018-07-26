package main

import (
	"bytes"
	"flag"
	"fmt"
)

func versionCmd() command {
	fs := flag.NewFlagSet("imposter version", flag.ExitOnError)
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return versionExec()
	}}
}

func versionExec() error {
	var r bytes.Buffer
	fmt.Fprintf(&r, "%s v%s", ProductName, Version)
	if VersionPrerelease != "" {
		fmt.Fprintf(&r, "-%s", VersionPrerelease)
	}

	fmt.Printf(r.String())
	return nil
}
