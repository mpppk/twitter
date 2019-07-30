package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"

	bolt "github.com/mpppk/bbolt"
	"github.com/mpppk/twitter/internal/option"
	"github.com/spf13/afero"

	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
)

func newDirFlag() *option.StringFlag {
	return &option.StringFlag{
		Flag: &option.Flag{
			IsDirName: true,
			Name:      "dir",
			Usage:     "downloaded images destination directory path",
		},
		Value: "images",
	}
}

func ImagesCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "Download images from DB file",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			conf, err := option.NewImagesCmdConfigFromViper()
			if err != nil {
				return err
			}

			cmd.Println("db path", conf.DBPath)
			db, err := bolt.Open(conf.DBPath, 0666, nil)
			if err != nil {
				return xerrors.Errorf("failed to open db file from %s: %w", conf.DBPath, err)
			}
			defer func() {
				cerr := db.Close()
				if cerr == nil {
					return
				}
				err = xerrors.Errorf("failed to close db(path: %s): %w", conf.DBPath, err)
			}()

			_ = os.Mkdir(conf.Dir, 0777)

			if err := downloadImageFromDB(cmd, db, conf.Dir); err != nil {
				return xerrors.Errorf("failed to download images: %w", err)
			}

			return err
		},
	}
	if err := option.RegisterStringFlag(cmd, newDirFlag()); err != nil {
		return nil, err
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, ImagesCmd)
}

func downloadImageFromDB(cmd *cobra.Command, db *bolt.DB, downloadDir string) error {
	return db.View(func(tx *bolt.Tx) error {

		//err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
		//	cmd.Println(string(name))
		//	return nil
		//})
		//if err != nil {
		//	return err
		//}
		//cmd.Println("-----")

		b := tx.Bucket([]byte("tweets"))
		if b == nil {
			return fmt.Errorf("failed to retrieve bucket which named %s", "tweets")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tweet anaconda.Tweet
			if err := json.Unmarshal(v, &tweet); err != nil {
				return xerrors.Errorf("failed to unmarshal tweet json: %w", err)
			}
			for i, entityMedia := range tweet.Entities.Media {
				mediaRawUrl := entityMedia.Media_url_https
				mediaUrl, err := url.Parse(mediaRawUrl)
				if err != nil {
					return xerrors.Errorf("failed to parse media url(%s): %w", mediaRawUrl, err)
				}
				mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
				if len(mediaUrlPaths) == 0 {
					return xerrors.Errorf("invalid mediaUrl: %s", mediaRawUrl)
				}
				mediaFileName := mediaUrlPaths[len(mediaUrlPaths)-1]
				fileName := fmt.Sprintf("%d-%d-%s", tweet.Id, i, mediaFileName)
				if isExist(fileName) {
					continue
				}
				downloadPath := path.Join(downloadDir, fileName)
				if err := downloadFile(mediaRawUrl, downloadPath); err != nil {
					return xerrors.Errorf("failed to download file to %s", downloadPath)
				}
				cmd.Printf("media is downloaded to %s\n", downloadPath)
				time.Sleep(10 * time.Second)

			}
		}
		return nil
	})
}

func downloadFile(fileUrl, downloadPath string) (err error) {
	response, err := http.Get(fileUrl)
	if err != nil {
		return xerrors.Errorf("failed to request http get to %s: %w", fileUrl, err)
	}
	defer func() {
		cerr := response.Body.Close()
		if cerr == nil {
			return
		}
		err = xerrors.Errorf("failed to close http response body: %w", err)
	}()

	file, err := os.Create(downloadPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		cerr := file.Close()
		if cerr == nil {
			return
		}
		err = xerrors.Errorf("failed to close file(%s): %w", downloadPath, err)
	}()

	if _, err := io.Copy(file, response.Body); err != nil {
		return xerrors.Errorf("failed to write download file to local file(%s): %w", downloadPath, err)
	}
	return nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
