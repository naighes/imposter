package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	commands := map[string]command{
		"start":   startCmd(),
		"version": versionCmd(),
	}
	fs := flag.NewFlagSet("imposter", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: imposter <command> [command flags]")
		for name, cmd := range commands {
			fmt.Printf("\n%s command:\n", name)
			cmd.fs.PrintDefaults()
		}
	}
	fs.Parse(os.Args[1:])
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}
	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}
