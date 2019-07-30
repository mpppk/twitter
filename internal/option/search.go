package option

// SearchCmdConfig is config for sum command
type SearchCmdConfig struct {
	Tag               string
	DBPath            string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// NewSearchCmdConfigFromViper generate config for sum command from viper
func NewSearchCmdConfigFromViper() (*SearchCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newSearchCmdConfigFromRawConfig(rawConfig), err
}

func newSearchCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *SearchCmdConfig {
	searchCmdConfig := &(rawConfig.SearchCmdConfig)
	searchCmdConfig.DBPath = rawConfig.DBPath
	return searchCmdConfig
}
