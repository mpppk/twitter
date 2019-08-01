package cmd

import (
	"os"

	"github.com/mpppk/twitter/internal/repository"

	"github.com/mpppk/twitter/internal/option"
	"github.com/spf13/afero"

	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
)

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

			repo, err := repository.New(conf.DBPath, cmd)
			if err != nil {
				return xerrors.Errorf("failed to create repo: %w", err)
			}
			defer func() {
				cerr := repo.Close()
				if cerr == nil {
					return
				}
				err = xerrors.Errorf("failed to close db(path: %s): %w", conf.DBPath, err)
			}()

			_ = os.Mkdir(conf.Dir, 0777)

			if err := repo.DownloadImageFromDB(conf.Dir); err != nil {
				return xerrors.Errorf("failed to download images: %w", err)
			}

			return err
		},
	}
	dirFlag := &option.StringFlag{
		Flag: &option.Flag{
			IsDirName: true,
			Name:      "dir",
			Usage:     "downloaded images destination directory path",
		},
		Value: "images",
	}
	if err := option.RegisterStringFlag(cmd, dirFlag); err != nil {
		return nil, err
	}

	intervalFlag := &option.IntFlag{
		Flag: &option.Flag{
			Name:  "interval",
			Usage: "interval between download images",
		},
		Value: 10,
	}
	if err := option.RegisterIntFlag(cmd, intervalFlag); err != nil {
		return nil, err
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, ImagesCmd)
}
