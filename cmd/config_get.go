package cmd

import (
	"fmt"

	"github.com/mpppk/twitter/internal/option"
	"github.com/mpppk/twitter/internal/repository"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func newGetCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Print specified config value",
		Long: fmt.Sprintf(`Print specified config value.
Available values are %s, %s.`, repository.MinIDKey, repository.MaxIDKey),
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{repository.MaxIDKey, repository.MinIDKey},
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
			case repository.MaxIDKey:
				maxId, err := repo.GetMaxID()
				if err != nil {
					cmd.Println("maxID does not stored")
					return nil
				}
				cmd.Println(maxId)
				return nil
			case repository.MinIDKey:
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
	configSubCmdGenerators = append(configSubCmdGenerators, newGetCmd)
}
