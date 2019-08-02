package option

// ImagesRawCmdConfig is raw config for sum command
type ImagesRawCmdConfig struct {
	Dir string
}

// ImagesCmdConfig is config for sum command
type ImagesCmdConfig struct {
	*ImagesRawCmdConfig
	DBPath string
}

// NewImagesCmdConfigFromViper generate config for sum command from viper
func NewImagesCmdConfigFromViper() (*ImagesCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newImagesCmdConfigFromRawConfig(rawConfig), err
}

func newImagesCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *ImagesCmdConfig {
	imagesCmdConfig := &ImagesCmdConfig{
		ImagesRawCmdConfig: &(rawConfig.ImagesRawCmdConfig),
		DBPath:             rawConfig.DBPath,
	}
	return imagesCmdConfig
}
