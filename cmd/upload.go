package cmd

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var cmdUpload = &cli.Command{
	Name: "upload",
	Flags: []cli.Flag{
		flagHost, flagFile, flagIndex,
	},
	Action: func(c *cli.Context) error {
		client := meilisearch.New(c.String("host"))

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
