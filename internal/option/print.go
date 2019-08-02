package option

// PrintCmdConfig is config for sum command
type PrintCmdConfig struct {
	DBPath string
}

// NewPrintCmdConfigFromViper generate config for sum command from viper
func NewPrintCmdConfigFromViper() (*PrintCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newPrintCmdConfigFromRawConfig(rawConfig), err
}

func newPrintCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *PrintCmdConfig {
	return &PrintCmdConfig{DBPath: rawConfig.DBPath}
}
