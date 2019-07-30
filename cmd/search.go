package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/xerrors"

	"github.com/ChimeraCoder/anaconda"
	bolt "github.com/mpppk/bbolt"
	"github.com/mpppk/twitter/internal/option"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newTagFlag() *option.StringFlag {
	return &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: true,
			Name:       "tag",
			Usage:      "hash tag",
		},
		Value: option.DefaultStringValue,
	}
}

func newDBPathFlag() *option.StringFlag {
	return &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: true,
			IsFileName: true,
			Name:       "dbPath",
			Usage:      "DB file path",
		},
		Value: option.DefaultStringValue,
	}
}

func newSearchCmd(fs afero.Fs) (cmd *cobra.Command, err error) {
	cmd = &cobra.Command{
		Use:   "search",
		Short: "search",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewSearchCmdConfigFromViper()
			if err != nil {
				return err
			}

			db, err := bolt.Open(conf.DBPath, 0666, nil)
			if err != nil {
				cmd.Println(conf)
				return xerrors.Errorf("failed to open db file from %s: %w", conf.DBPath, err)
			}
			defer func() {
				err = db.Close()
			}()

			if err := createBucketIfNotExists(db, "tweets"); err != nil {
				return xerrors.Errorf("failed to create bucket which named %s: %w", "tweets", err)
			}

			api := anaconda.NewTwitterApiWithCredentials(
				conf.AccessToken,
				conf.AccessTokenSecret,
				conf.ConsumerKey,
				conf.ConsumerSecret)

			v := url.Values{}
			v.Set("count", "100")

			queries := []string{conf.Tag, "exclude:retweets", "filter:images"}
			searchResult, err := api.GetSearch(strings.Join(queries, " "), v)
			if err != nil {
				return xerrors.Errorf("failed to search tweets: %w", err)
			}
			for _, tweet := range searchResult.Statuses {
				idBytes := make([]byte, binary.MaxVarintLen64)
				tweetJsonBytes, err := json.Marshal(tweet)
				if err != nil {
					return err
				}
				err = db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("tweets"))
					binary.PutVarint(idBytes, tweet.Id)
					return b.Put(
						idBytes,
						tweetJsonBytes,
					)
				})
				if err != nil {
					return err
				}
				//cmd.Println(tweet.Text)
			}

			err = db.View(func(tx *bolt.Tx) error {
				// Assume bucket exists and has keys
				b := tx.Bucket([]byte("tweets"))

				c := b.Cursor()

				for k, v := c.First(); k != nil; k, v = c.Next() {
					cmd.Println(string(v))
				}
				return nil
			})
			if err != nil {
				return err
			}

			return err
		},
	}
	if err := option.RegisterStringFlag(cmd, newTagFlag()); err != nil {
		return nil, err
	}
	if err := option.RegisterStringFlag(cmd, newDBPathFlag()); err != nil {
		return nil, err
	}
	return cmd, err
}

func init() {
	cmdGenerators = append(cmdGenerators, newSearchCmd)
}
func createBucketIfNotExists(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %s", bucketName, err)
		}
		return nil
	})
}
