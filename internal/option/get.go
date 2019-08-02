package option

// GetCmdConfig is config for sum command
type GetCmdConfig struct {
	DBPath string
}

// NewGetCmdConfigFromViper generate config for sum command from viper
func NewGetCmdConfigFromViper() (*GetCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newGetCmdConfigFromRawConfig(rawConfig), err
}

func newGetCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *GetCmdConfig {
	return &GetCmdConfig{DBPath: rawConfig.DBPath}
}
