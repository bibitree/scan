package app

import (
	"errors"
	"ethgo/chainFinder"
	"ethgo/model"
	"ethgo/util/logx"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ChainFinder *chainFinder.Config `toml:"ChainFinder" json:"chainFinder"`
	Logger      *logx.Config        `toml:"logger" json:"logger"`
	Redis       *model.Config       `toml:"redis" json:"redis"`
}

func NewConfig(filepath string) (*Config, error) {
	var c = new(Config)
	if _, err := toml.DecodeFile(filepath, c); err != nil {
		return nil, err
	}
	if err := c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Init() error {

	// if err := c.Tyche.Init(); err != nil {
	// 	return err
	// }

	if c.Redis.Domain == "" {
		return errors.New("domain cannot be set to empty")
	}

	var namespace = fmt.Sprintf("%v", c.Redis.Domain)
	c.Redis.Namespace = namespace
	return nil
}
