package cmd

import (
	"sync"

	"github.com/mpppk/twitter/internal/option"
	"github.com/mpppk/twitter/internal/repository"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newPrintCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print tweets in local file DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewPrintCmdConfigFromViper()
			if err != nil {
				return err
			}
			repo, err := repository.New(conf.DBPath, cmd)
			defer func() {
				err = repo.Close()
			}()

			ch := make(chan string)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func(ch chan string) {
				for tweetStr := range ch {
					cmd.Println(tweetStr)
				}
				wg.Done()
			}(ch)

			if err := repo.SendTweetStrToChannel(ch); err != nil {
				return err
			}
			close(ch)
			wg.Wait()

			return nil
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newPrintCmd)
}
