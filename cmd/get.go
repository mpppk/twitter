package cmd

import (
	"github.com/mpppk/twitter/internal/option"
	"github.com/mpppk/twitter/internal/repository"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func newGetCmd(fs afero.Fs) (*cobra.Command, error) {
	maxIDArg := "maxID"
	minIDArg := "minID"
	cmd := &cobra.Command{
		Use:       "get",
		Short:     "Print config property",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{maxIDArg, minIDArg},
		//Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			conf, err := option.NewGetCmdConfigFromViper()
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

			switch target {
			case maxIDArg:
				maxId, err := repo.GetMaxID()
				if err != nil {
					cmd.Println("maxID does not stored")
					return nil
				}
				cmd.Println(maxId)
				return nil
			case minIDArg:
				minId, err := repo.GetMinID()
				if err != nil {
					cmd.Println("minID does not stored")
					return nil
				}
				cmd.Println(minId)
				return nil
			}
			return xerrors.Errorf("unknown target: %s", target)
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newGetCmd)
}
