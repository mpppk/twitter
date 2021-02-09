package option

import (
	"strings"

	"golang.org/x/xerrors"
)

// SearchCmdConfig is config for search command
type SearchCmdConfig struct {
	*SearchRawCmdConfig
	IgnoreKeywords []string
	Excludes       []string
	Filters        []string
	DBPath         string
}

func (s *SearchCmdConfig) validate() error {
	if s.DBPath == "" {
		return xerrors.Errorf("--db-path must be provided")
	}

	if s.ConsumerKey == "" {
		return xerrors.Errorf("ConsumerKey must be provided via config file. Put .twitter.yml to ~/.config")
	}
	if s.ConsumerSecret == "" {
		return xerrors.Errorf("ConsumerSecret must be provided via config file. Put .twitter.yml to ~/.config")
	}
	if s.AccessToken == "" {
		return xerrors.Errorf("AccessToken must be provided via config file. Put .twitter.yml to ~/.config")
	}
	if s.AccessTokenSecret == "" {
		return xerrors.Errorf("AccessTokenSecret must be provided via config file. Put .twitter.yml to ~/.config")
	}
	return nil
}

// SearchCmdConfig is raw config for search command
type SearchRawCmdConfig struct {
	Query             string
	Ignore            string
	Exclude           string
	Interval          int
	Filter            string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// NewSearchCmdConfigFromViper generate config for search command from viper
func NewSearchCmdConfigFromViper() (*SearchCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	if err != nil {
		return nil, err
	}
	searchCmdConfig := newSearchCmdConfigFromRawConfig(rawConfig)
	return searchCmdConfig, searchCmdConfig.validate()
}

func newSearchCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *SearchCmdConfig {
	searchCmdConfig := &SearchCmdConfig{
		SearchRawCmdConfig: &(rawConfig.SearchRawCmdConfig),
	}
	searchCmdConfig.DBPath = rawConfig.DBPath

	searchCmdConfig.IgnoreKeywords = strings.Split(searchCmdConfig.Ignore, ",")
	searchCmdConfig.Excludes = strings.Split(searchCmdConfig.Exclude, ",")
	searchCmdConfig.Filters = strings.Split(searchCmdConfig.Filter, ",")

	return searchCmdConfig
}
