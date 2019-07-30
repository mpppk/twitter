package cmd

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
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

func newSearchCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewSearchCmdConfigFromViper()
			if err != nil {
				return err
			}

			api := anaconda.NewTwitterApiWithCredentials(
				conf.AccessToken,
				conf.AccessTokenSecret,
				conf.ConsumerKey,
				conf.ConsumerSecret)

			v := url.Values{}
			v.Set("count", "30")

			searchResult, err := api.GetSearch(conf.Tag, v)
			if err != nil {
				return err
			}
			for _, tweet := range searchResult.Statuses {
				cmd.Println(tweet.Text)
			}
			cmd.Println(conf.Tag)
			return nil
		},
	}
	if err := option.RegisterStringFlag(cmd, newTagFlag()); err != nil {
		return nil, err
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newSearchCmd)
}
