package option

// SetCmdConfig is config for sum command
type SetCmdConfig struct {
	DBPath string
}

// NewSetCmdConfigFromViper generate config for sum command from viper
func NewSetCmdConfigFromViper() (*SetCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newSetCmdConfigFromRawConfig(rawConfig), err
}

func newSetCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *SetCmdConfig {
	return &SetCmdConfig{DBPath: rawConfig.DBPath}
}
