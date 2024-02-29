package cmd

import (
	"bufio"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-json"
	"github.com/meilisearch/meilisearch-go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func getDocumentCount(idx *meilisearch.Index) (int64, error) {
	total := int64(-1)
	{
		q := meilisearch.DocumentsQuery{
			Offset: 0,
			Limit:  1,
		}
		rsp := meilisearch.DocumentsResult{}
		err := idx.GetDocuments(&q, &rsp)
		if err != nil {
			return 0, err
		}
		total = rsp.Total
	}
	return total, nil
}

var cmdClean = &cli.Command{
	Name: "clean",
	Flags: []cli.Flag{
		flagFile, flagIndex, flagHost, flagFetchBatch,
	},
	Action: func(c *cli.Context) error {
		set := make(map[string]struct{})

		file, err := os.Open(c.String("file"))
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var m map[string]interface{}
			err := json.Unmarshal([]byte(scanner.Text()), &m)
			if err != nil {
				return errors.Wrap(err, "Unmarshal failed")
			}
			id := m["id"].(string)
			set[id] = struct{}{}
		}
		if err := scanner.Err(); err != nil {
			return errors.Wrap(err, "scan failed")
		}

		toDelete := make(map[string]struct{})

		indexName := c.String("index")
		client := meilisearch.NewClient(meilisearch.ClientConfig{
			Host: c.String("host"),
		})

		index := client.Index(indexName)
		if index == nil {
			return errors.Newf("index %s not exists", indexName)
		}

		total, err := getDocumentCount(index)
		if err != nil {
			return err
		}
		log.Info().Int64("total", total).Msg("document count")

		offset := int64(0)
		batch := c.Int64("fetch-batch")
		for offset < total {
			q := meilisearch.DocumentsQuery{
				Offset: offset,
				Limit:  batch,
			}
			rsp := meilisearch.DocumentsResult{}
			err := index.GetDocuments(&q, &rsp)
			if err != nil {
				return err
			}

			offset += rsp.Limit

			for _, v := range rsp.Results {
				id := v["id"].(string)
				_, ok := set[id]
				if !ok {
					toDelete[id] = struct{}{}
				}
			}
		}

		log.Info().Int("delete-count", len(toDelete)).Msg("delete plan prepared")

		deleteList := make([]string, 0)
		for id, _ := range toDelete {
			deleteList = append(deleteList, id)
		}

		if len(deleteList) == 0 {
			log.Info().Msg("Nothing to delete")
			return nil
		}

		rsp, err := index.DeleteDocuments(deleteList)
		if err != nil {
			return err
		}
		log.Info().Interface("response", rsp).Msg("ok")

		return nil
	},
}
