package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/sitemap/command"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Michael Boudreau"
	app.Email = "https://github.com/mkboudreau/sitemap"
	app.Version = "1.0"
	app.Usage = "Builds a sitemap for a url argument and formats its output"
	app.Action = command.AppAction()
	app.Flags = command.AppCliFlags()

	err := app.Run(os.Args)
	if err != nil {
		panic(fmt.Errorf("Couldn't start application %v", err))
	}
}
