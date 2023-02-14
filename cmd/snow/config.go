package main

import (
	"avalanche-consensus/consensus"
	"avalanche-consensus/p2pnetworking"

	"github.com/spf13/viper"
)

type SnowConfig struct {
	P2p           p2pnetworking.Config `json:"p2p" mapstructure:"P2P"`
	Consensus     consensus.Config     `json:"consensus" mapstructure:"CONSENSUS"`
	NumberOfNode  int                  `json:"numberOfNode" mapstructure:"NUMBER_OF_NODE"`
	NumberOfBlock int                  `json:"numberOfBlock" mapstructure:"NUMBER_OF_BLOCK"`
}

func LoadConfig(path string) (*SnowConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("snow")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := SnowConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
