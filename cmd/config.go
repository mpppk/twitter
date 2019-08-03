package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var configSubCmdGenerators []cmdGenerator

func newConfigCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration of twitter CLI and DB file",
		//Long: ``,
	}

	var subCmds []*cobra.Command
	for _, cmdGen := range configSubCmdGenerators {
		subCmd, err := cmdGen(fs)
		if err != nil {
			return nil, err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newConfigCmd)
}
