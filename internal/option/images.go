package option

// ImagesCmdConfig is config for sum command
type ImagesCmdConfig struct {
	Dir    string
	DBPath string
}

// NewImagesCmdConfigFromViper generate config for sum command from viper
func NewImagesCmdConfigFromViper() (*ImagesCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newImagesCmdConfigFromRawConfig(rawConfig), err
}

func newImagesCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *ImagesCmdConfig {
	imagesCmdConfig := &(rawConfig.ImagesCmdConfig)
	// FIXME
	imagesCmdConfig.DBPath = rawConfig.SearchCmdConfig.DBPath
	return imagesCmdConfig
}
