package cmd

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
)

var flagFile = &cli.StringFlag{
	Name:    "file",
	Aliases: []string{"f"},
	Action: func(context *cli.Context, s string) error {
		_, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		return nil
	},
}

var flagHost = &cli.StringFlag{
	Name:    "host",
	Aliases: []string{"d"},
	Action: func(context *cli.Context, s string) error {
		if s != "" {
			return nil
		}
		return errors.New("invalid dst addr")
	},
}

var flagIndex = &cli.StringFlag{
	Name:    "index",
	Aliases: []string{"i"},
}

var flagFetchBatch = &cli.Int64Flag{
	Name:    "fetch-batch",
	Aliases: []string{"b"},
}
