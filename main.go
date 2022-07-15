package main

import (
	"flag"
	"fmt"
	app "github.com/mtwig/gh-notification-purge/app"
	"os"
)

func main() {

	var cfg app.Config
	flag.BoolVar(&cfg.PrintSubjects, "display-subjects", true, "display the PR subject for ")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "do not really mark as read")
	flag.Parse()

	app, err := app.CreateNew(cfg)
	if err != nil {
		fmt.Printf("unable to create app.\n%s\n", err)
		os.Exit(1)
	}

	app.Run()
}
