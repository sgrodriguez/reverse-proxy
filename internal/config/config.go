package config

import (
	"reverseproxy/internal/blocker"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	TargetURL        string                     `toml:"TargetURL"`
	ReverseProxyPort int                        `toml:"ReverseProxyPort"`
	HeaderBlocker    *blocker.HeaderBlocker     `toml:"HeaderBlocker"`
	ParamBlocker     *blocker.QueryParamBlocker `toml:"ParamBlocker"`
	PathBlocker      *blocker.PathBlocker       `toml:"PathBlocker"`
	MethodBlocker    *blocker.MethodBlocker     `toml:"MethodBlocker"`
}

// LoadConfig from toml file
func LoadConfig(tomlPath string) (*Config, error) {
	conf := &Config{}
	_, err := toml.DecodeFile(tomlPath, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
