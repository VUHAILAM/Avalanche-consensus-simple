package consensus

import "errors"

type Config struct {
	//K is sample size
	K int `yaml:"k" mapstructure:"k"`
	//Alphal is quorum size
	Alphal int `yaml:"alphal" mapstructure:"alphal"`
	//Beta is decision threshold
	Beta int `yaml:"beta" mapstructure:"beta"`
}

func (c *Config) Verify() error {
	if c.Alphal <= c.K/2 || c.Alphal > c.K {
		return errors.New("alpha must be in (k/2, k]")
	}

	if c.Beta < 1 {
		return errors.New("beta must be >= 1")
	}
	return nil
}
