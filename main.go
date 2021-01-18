package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/weberc2/neon/build"
	"github.com/weberc2/neon/serve"
	"gopkg.in/urfave/cli.v2"
)

func init() {
	// noop by  default
	build.DefaultLogFunc = func(v ...interface{}) {}
}

func main() {
	const verboseFlagName = "verbose"
	const portFlagName = "port"

	app := cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "build",
				Usage: "Build the current project",
				Action: func(c *cli.Context) error {
					if c.Bool(verboseFlagName) {
						build.DefaultLogFunc = func(v ...interface{}) {
							log.Println(v...)
							log.Println(string(debug.Stack()))
						}
					}
					return build.Build()
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Aliases: []string{"v"},
						Name:    verboseFlagName,
						Usage:   "Set to enable verbose error messages",
						Value:   false,
					},
				},
			},
			&cli.Command{
				Name:  "serve",
				Usage: "Build and serve the project's output directory",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    portFlagName,
						Aliases: []string{"p"},
						Value:   8080,
						Usage:   "The port the HTTP server will listen on",
					},
				},
				Action: func(context *cli.Context) error {
					return serve.Serve(uint16(context.Int(portFlagName)))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
