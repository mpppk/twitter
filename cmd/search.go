package cmd

import (
	"time"

	"github.com/mpppk/twitter/internal/twitter"

	"github.com/mpppk/twitter/internal/repository"

	"github.com/mpppk/twitter/internal/option"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newSearchCmd(fs afero.Fs) (cmd *cobra.Command, err error) {
	cmd = &cobra.Command{
		Use:   "search",
		Short: "search",
		Long: `Search tweets by query and some options.
Results are stored in local file DB. (You can specify the DB path by --db-path flag.
If you want to download images which contained in tweets, execute 'download images' command after search command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewSearchCmdConfigFromViper()
			if err != nil {
				return err
			}

			cmd.Println("db path", conf.DBPath)
			repo, err := repository.New(conf.DBPath, cmd)
			defer func() {
				err = repo.Close()
			}()

			client := twitter.NewClient(
				conf.AccessToken,
				conf.AccessTokenSecret,
				conf.ConsumerKey,
				conf.ConsumerSecret)

			maxId, err := repo.GetMaxID()
			if err != nil {
				cmd.Println("saved maxID does not found")
				maxId = -1
			} else {
				cmd.Println("Retrieved maxID: ", maxId)
			}

			query := twitter.BuildQuery(conf.Query, conf.IgnoreKeywords, conf.Excludes, conf.Filters)
			cmd.Printf("Search query: %q\n", query)

			for {
				tweets, err := client.SearchTweets(query, maxId, -1)
				if err != nil {
					cmd.Println("failed to search:", err)
					cmd.Printf("sleep %d sec...\n", conf.Interval)
					time.Sleep(time.Duration(conf.Interval) * time.Second)
					continue
				}

				if len(tweets) == 0 {
					break
				}

				for _, tweet := range tweets {
					if err := repo.SaveTweet(&tweet); err != nil {
						return err
					}
				}

				lastTweet := tweets[len(tweets)-1]
				if err := repo.SetMaxId(lastTweet.Id - 1); err != nil {
					return err
				}
				maxId = lastTweet.Id - 1
				cmd.Printf("%d tweets are saved. (%d-%d)\n", len(tweets), tweets[0].Id, lastTweet.Id)
				cmd.Printf("maxID is updated => %d\n", maxId)
			}

			// 最新のtweet取得部分を実装
			minId, err := repo.GetMinID()
			if err != nil {
				return nil
			}

			oldestTweetId := int64(-1)
			latestTweetId := int64(-1)
			for {
				tweets, err := client.SearchTweets(query, oldestTweetId, minId)
				if err != nil {
					cmd.Println("failed to search:", err)
					cmd.Printf("sleep %d sec...\n", conf.Interval)
					time.Sleep(time.Duration(conf.Interval) * time.Second)
					continue
				}

				if len(tweets) == 0 {
					break
				}

				for _, tweet := range tweets {
					if err := repo.SaveTweet(&tweet); err != nil {
						return err
					}
				}

				latestTweet := tweets[0]
				if latestTweetId < latestTweet.Id {
					latestTweetId = latestTweet.Id
				}
				oldestTweet := tweets[len(tweets)-1]
				if oldestTweetId < 0 || oldestTweetId > oldestTweet.Id {
					oldestTweetId = oldestTweet.Id
				}

				cmd.Printf("%d tweets are saved. (%d-%d)\n", len(tweets), oldestTweet.Id, latestTweet.Id)
			}
			if err := repo.SetMinId(latestTweetId); err != nil {
				return err
			}
			cmd.Printf("minID is updated => %d\n", minId)
			return nil
		},
	}
	queryFlag := &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: true,
			Name:       "query",
			Usage:      "search query",
		},
	}
	if err := option.RegisterStringFlag(cmd, queryFlag); err != nil {
		return nil, err
	}

	ignoreFlag := &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: false,
			Name:       "ignore",
			Usage:      "ignore keywords",
		},
		Value: "",
	}
	if err := option.RegisterStringFlag(cmd, ignoreFlag); err != nil {
		return nil, err
	}

	excludeFlag := &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: true,
			Name:       "exclude",
			Usage:      "exclude tweet type",
		},
	}
	if err := option.RegisterStringFlag(cmd, excludeFlag); err != nil {
		return nil, err
	}

	filterFlag := &option.StringFlag{
		Flag: &option.Flag{
			IsRequired: true,
			Name:       "filter",
			Usage:      "filter tweet type",
		},
	}
	if err := option.RegisterStringFlag(cmd, filterFlag); err != nil {
		return nil, err
	}

	intervalFlag := &option.IntFlag{
		Flag: &option.Flag{
			Name:  "interval",
			Usage: "Interval sec between API request failure and rerun",
		},
		Value: 60,
	}
	if err := option.RegisterIntFlag(cmd, intervalFlag); err != nil {
		return nil, err
	}

	return cmd, err
}

func init() {
	cmdGenerators = append(cmdGenerators, newSearchCmd)
}
