package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

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
			if err := createBucketIfNotExists(db, "maxID"); err != nil {
				return xerrors.Errorf("failed to create bucket which named %s: %w", "maxID", err)
			}

			api := anaconda.NewTwitterApiWithCredentials(
				conf.AccessToken,
				conf.AccessTokenSecret,
				conf.ConsumerKey,
				conf.ConsumerSecret)

			v := url.Values{}
			v.Set("count", "100")

			if maxID, err := getMaxID(db); err == nil {
				lastTweetIdStr := fmt.Sprintf("%d", maxID-1)
				v.Set("max_id", lastTweetIdStr)
				cmd.Println("retrieve max ID: ", maxID)
			}

			queries := []string{conf.Tag, "OR", "#" + conf.Tag, "exclude:retweets", "filter:images"}
			for {
				searchResult, err := api.GetSearch(strings.Join(queries, " "), v)
				if err != nil {
					cmd.Println("failed to search: %s", err)
					cmd.Println("sleep 60 sec...")
					time.Sleep(60 * time.Second)
					continue
				}

				if len(searchResult.Statuses) == 0 {
					return nil
				}

				tweets := searchResult.Statuses
				for _, tweet := range tweets {
					if err := saveTweet(db, &tweet); err != nil {
						return err
					}
					//cmd.Println(tweet.Text)
				}
				lastTweet := tweets[len(tweets)-1]
				if _, err := saveMaxId(db, lastTweet.Id-1); err != nil {
					return err
				}

				lastTweetIdStr := fmt.Sprintf("%d", lastTweet.Id-1)
				v.Set("max_id", lastTweetIdStr)
				cmd.Println("new max id: " + lastTweetIdStr)
			}
			//if err := printDB(cmd, db); err != nil {
			//	return err
			//}
			//return nil
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

func saveTweet(db *bolt.DB, tweet *anaconda.Tweet) error {
	idBytes := make([]byte, binary.MaxVarintLen64)
	tweetJsonBytes, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tweets"))
		binary.PutVarint(idBytes, tweet.Id)
		return b.Put(
			idBytes,
			tweetJsonBytes,
		)
	})
}

func saveMaxId(db *bolt.DB, maxId int64) (bool, error) {
	currentMaxID, err := getMaxID(db)
	if err == nil && currentMaxID <= maxId {
		return false, nil
	}

	idBytes := make([]byte, binary.MaxVarintLen64)
	return true, db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maxID"))
		binary.PutVarint(idBytes, maxId)
		return b.Put(
			[]byte("maxID"),
			idBytes,
		)
	})
}

func getMaxID(db *bolt.DB) (int64, error) {
	var maxID int64
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maxID"))
		maxIDBytes := b.Get([]byte("maxID"))
		mid, bufSize := binary.Varint(maxIDBytes) // FIXME error handling
		if mid == 0 && bufSize == 0 {
			return fmt.Errorf("buf too small")
		}
		if mid == 0 && bufSize < 0 {
			return fmt.Errorf("value larger than 64 bits (overflow)")
		}
		maxID = mid
		return nil
	})
	if err != nil {
		return 0, xerrors.Errorf("failed to retrieve tweet max id from db: %w", err)
	}
	return maxID, nil
}

func printDB(cmd *cobra.Command, db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tweets"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			cmd.Println(string(v))
		}
		return nil
	})
}

func createBucketIfNotExists(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %s", bucketName, err)
		}
		return nil
	})
}
