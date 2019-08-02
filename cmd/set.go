package cmd

import (
	"fmt"
	"strconv"

	"github.com/mpppk/twitter/internal/option"
	"github.com/mpppk/twitter/internal/repository"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func newSetCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set config property",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return xerrors.Errorf("accepts 2 args, received %d", len(args))
			}
			if !repository.IsValidKey(args[0]) {
				return fmt.Errorf("invalid key specified: %s", args[0])
			}

			if _, err := strconv.Atoi(args[1]); err != nil {
				return fmt.Errorf("invalid value specified: %s", args[1])
			}
			return nil
		},
		//Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			conf, err := option.NewSetCmdConfigFromViper()
			if err != nil {
				return err
			}
			repo, err := repository.New(conf.DBPath, cmd)
			if err != nil {
				return err
			}
			defer func() {
				err = repo.Close()
			}()

			// error is already checked in Args func
			v, _ := strconv.Atoi(args[1])

			switch target {
			case repository.MaxIDKey:
				if err := repo.SetMaxId(int64(v)); err != nil {
					cmd.Println("maxID does not stored")
				}
				return nil
			case repository.MinIDKey:
				err := repo.SetMinId(int64(v))
				if err != nil {
					cmd.Println("minID does not stored")
				}
				return nil
			}
			return xerrors.Errorf("unknown target: %s", target)
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newSetCmd)
}
