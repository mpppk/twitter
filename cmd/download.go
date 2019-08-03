package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var downloadSubCmdGenerators []cmdGenerator

func newDownloadCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download twitter contents from local file DB",
		//Long: ``,
	}

	var subCmds []*cobra.Command
	for _, cmdGen := range downloadSubCmdGenerators {
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
	cmdGenerators = append(cmdGenerators, newDownloadCmd)
}
