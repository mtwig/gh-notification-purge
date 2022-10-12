package main

import (
	"flag"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/mtwig/gh-notification-purge/app"
	"os"
)

func main() {

	var cfg app.Config
	flag.BoolVar(&cfg.PrintSubjects, "display-subjects", true, "display the PR subject for ")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "do not really mark as read")
	flag.Parse()

	instance, err := app.CreateNew(cfg)
	if err != nil {
		fmt.Printf("unable to create instance.\n%s\n", err)
		os.Exit(1)
	}

	err = instance.Run()
	if err != nil {
		fmt.Printf(color.Colorize(color.Red, fmt.Sprintf("Error: %s", err)))
	}
}
