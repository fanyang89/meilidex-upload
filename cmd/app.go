package cmd

import "github.com/urfave/cli/v2"

func NewApp() *cli.App {
	return &cli.App{
		Name: "meilidex-upload",
		Commands: []*cli.Command{
			cmdUpload,
			cmdClean,
		},
	}
}
