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
		Long:  ``,
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
				cmd.Println("saved max ID does not found")
				maxId = -1
			} else {
				cmd.Println("retrieved max ID: ", maxId)
			}

			query := twitter.BuildQuery(conf.Query, conf.Excludes, conf.Filters)
			cmd.Printf("Search query: %q\n", query)

			for {
				tweets, err := client.SearchTweets(query, maxId, -1)
				if err != nil {
					cmd.Println("failed to search: %s", err)
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
				if _, err := repo.SaveMaxId(lastTweet.Id - 1); err != nil {
					return err
				}

				maxId = lastTweet.Id - 1
				cmd.Printf("new max id: %d\n", maxId)
			}

			// 最新のtweet取得部分を実装
			if minId, err := repo.GetMinID(); err == nil {
				for {
					tweets, err := client.SearchTweets(query, -1, minId)
					if err != nil {
						cmd.Println("failed to search: %s", err)
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

					firstTweet := tweets[0]
					if _, err := repo.SaveMinId(firstTweet.Id); err != nil {
						return err
					}

					minId = firstTweet.Id
					cmd.Printf("new min id: %d\n", minId)
				}
			}
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
