package repository

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mpppk/twitter/pkg/io"

	"github.com/ChimeraCoder/anaconda"
	bolt "github.com/mpppk/bbolt"
	"golang.org/x/xerrors"
)

type Repository struct {
	db     *bolt.DB
	logger Logger
}

type Logger interface {
	Println(i ...interface{})
	Printf(format string, i ...interface{})
}

func New(dbPath string, logger Logger) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to open db file from %s: %w", dbPath, err)
	}

	if err := createBucketIfNotExists(db, tweetsBucketName); err != nil {
		return nil, xerrors.Errorf("failed to create bucket which named %s: %w", tweetsBucketName, err)
	}
	if err := createBucketIfNotExists(db, "maxID"); err != nil {
		return nil, xerrors.Errorf("failed to create bucket which named %s: %w", "maxID", err)
	}

	return &Repository{db: db, logger: logger}, nil
}

func (r *Repository) Close() error {
	err := r.db.Close()
	if err != nil {
		return xerrors.Errorf("failed to close repository: %w", err)
	}
	return nil
}

func (r *Repository) DownloadImageFromDB(downloadDir string) error {
	return r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tweetsBucketName))
		if b == nil {
			return fmt.Errorf("failed to retrieve bucket which named %s", tweetsBucketName)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tweet anaconda.Tweet
			if err := json.Unmarshal(v, &tweet); err != nil {
				r.logger.Println(xerrors.Errorf("failed to unmarshal tweet json: %w", err))
				continue
			}

			for i, entityMedia := range tweet.Entities.Media {
				downloadPath, err := io.DownloadEntityMedia(&tweet, &entityMedia, i, downloadDir)
				if err != nil {
					r.logger.Println(xerrors.Errorf("failed to download entity media from %s: %w", entityMedia.Media_url_https, err))
					continue
				}
				if downloadPath != "" {
					r.logger.Printf("media is downloaded to %s\n", downloadPath)
					time.Sleep(10 * time.Second)
				}
			}
			for i, entityMedia := range tweet.ExtendedEntities.Media {
				downloadPath, err := io.DownloadEntityMedia(&tweet, &entityMedia, i, downloadDir)
				if err != nil {
					r.logger.Println(xerrors.Errorf("failed to download extended entity media from %s: %w", entityMedia.Media_url_https, err))
					continue
				}
				if downloadPath != "" {
					r.logger.Printf("extended media is downloaded to %s\n", downloadPath)
					time.Sleep(10 * time.Second)
				}
			}
		}
		return nil
	})
}

func (r *Repository) GetMaxID() (int64, error) {
	return getInt64(r.db, idBucketName, MaxIDKey)
}

func (r *Repository) GetMinID() (int64, error) {
	return getInt64(r.db, idBucketName, MinIDKey)
}

func (r *Repository) SetMinId(minId int64) error {
	return saveInt64(r.db, idBucketName, MinIDKey, minId)
}

func (r *Repository) SaveMinId(minId int64) (bool, error) {
	currentMinID, err := getInt64(r.db, idBucketName, MinIDKey)
	if err == nil && currentMinID >= minId {
		return false, nil
	}
	return true, r.SetMinId(minId)
}

func (r *Repository) SetMaxId(maxId int64) error {
	return saveInt64(r.db, idBucketName, MaxIDKey, maxId)
}

func (r *Repository) SaveMaxId(maxId int64) (bool, error) {
	currentMaxID, err := getInt64(r.db, idBucketName, MaxIDKey)
	if err == nil && currentMaxID <= maxId {
		return false, nil
	}
	return true, r.SetMaxId(maxId)
}

func (r *Repository) SaveTweet(tweet *anaconda.Tweet) error {
	idBytes := make([]byte, binary.MaxVarintLen64)
	tweetJsonBytes, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tweetsBucketName))
		binary.PutVarint(idBytes, tweet.Id)
		return b.Put(
			idBytes,
			tweetJsonBytes,
		)
	})
}

// SendTweetStrToChannel send tweet strings in DB to provided channel
func (r *Repository) SendTweetStrToChannel(ch chan string) error {
	return r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tweetsBucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			ch <- string(v)
		}
		return nil
	})
}
