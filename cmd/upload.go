package cmd

import (
	"errors"
	"os"

	"github.com/meilisearch/meilisearch-go"
	"github.com/rs/zerolog/log"
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

var cmdUpload = &cli.Command{
	Name: "upload",
	Flags: []cli.Flag{
		flagHost, flagFile, flagIndex,
	},
	Action: func(c *cli.Context) error {
		client := meilisearch.NewClient(meilisearch.ClientConfig{
			Host: c.String("host"),
		})

		indexName := c.String("index")
		task, err := client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        indexName,
			PrimaryKey: "id",
		})
		if err != nil {
			return err
		}

		index := client.Index(indexName)
		fh, err := os.Open(c.String("file"))
		if err != nil {
			return err
		}
		defer func() { _ = fh.Close() }()

		task, err = index.AddDocumentsNdjsonFromReader(fh, "id")
		if err != nil {
			return err
		}

		log.Info().Interface("task", task).Msg("ok")
		return nil
	},
}
