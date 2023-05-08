package app

import (
	"errors"
	"ethgo/eth"
	"ethgo/model"
	"ethgo/tyche"
	"ethgo/util/logx"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger  *logx.Config  `toml:"logger" json:"logger"`
	Backend *eth.Config   `toml:"backend" json:"backend"`
	Redis   *model.Config `toml:"redis" json:"redis"`
	Tyche   *tyche.Config `toml:"tyche" json:"tyche"`
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
	if err := c.Backend.Init(); err != nil {
		return err
	}

	if err := c.Tyche.Init(); err != nil {
		return err
	}

	if c.Redis.Domain == "" {
		return errors.New("domain cannot be set to empty")
	}

	var namespace = fmt.Sprintf("%v:%v", c.Redis.Domain, c.Tyche.Account)
	c.Redis.Namespace = namespace
	return nil
}
