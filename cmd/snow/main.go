package main

import (
	"avalanche-consensus/api"
	"avalanche-consensus/chain"
	"avalanche-consensus/consensus"
	"avalanche-consensus/model"
	"avalanche-consensus/p2pnetworking"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	discovery, err := runDiscovery()
	if err != nil {
		Logger.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	config, err := LoadConfig(".")
	if err != nil {
		Logger.Info(err)
	}
	api.InitApp(discovery, config.P2p, config.Consensus, config.NumberOfBlock, config.AppAddress)
	var wg sync.WaitGroup

	for i := 0; i < config.NumberOfNode; i++ {
		wg.Add(1)
		go SetupNode(&wg, sigs, i, discovery, config)
	}
	wg.Wait()

}

func SetupNode(wg *sync.WaitGroup, sigs chan os.Signal, i int, discovery *p2pnetworking.Discovery, conf *SnowConfig) {
	if discovery == nil {
		Logger.Fatal("Discovery is nil")
	}
	doneChan := make(chan bool, 1)
	defer wg.Done()
	go func() {
		<-sigs
		doneChan <- true
	}()
	ctx := context.Background()
	freePort, err := freeport.GetFreePort()
	if err != nil {
		Logger.Fatal(err)
		return
	}
	conf.P2p.Port = freePort

	node, err := chain.InitNode(ctx, conf.P2p, conf.Consensus, discovery)
	if err != nil {
		Logger.Fatal(err)
		return
	}
	time.Sleep(1 * time.Second)
	// Random init data for blockchain
	for j := 0; j < conf.NumberOfBlock; j++ {
		r := rand.Intn(int(float32(2) * 1.5))
		block := &chain.Block{}
		if r < 2 {
			block.SetData(model.DataType(r))
		} else {
			block.SetData(model.DataType(j))
		}
		node.Add(block)
	}
	//Get data before sync
	beforeBlockChainState := ""
	for _, b := range node.Blocks {
		beforeBlockChainState += fmt.Sprintf("%d ", int(b.Data))
	}
	Logger.Infof("Before sync, data of node: %d is %s", i, beforeBlockChainState)

	err = node.Sync(ctx)
	if err != nil {
		Logger.Fatal(err)
		return
	}

	// Get data after sync
	blockChainState := ""
	for _, b := range node.Blocks {
		blockChainState += fmt.Sprintf("%d ", int(b.Data))
	}

	Logger.Infof("client: %d, block: %s", i, blockChainState)
	<-doneChan
}

type SnowConfig struct {
	P2p           p2pnetworking.Config `yaml:"p2p" mapstructure:"p2p"`
	Consensus     consensus.Config     `yaml:"consensus" mapstructure:"consensus"`
	NumberOfNode  int                  `yaml:"numberOfNode" mapstructure:"numberOfNode"`
	NumberOfBlock int                  `yaml:"numberOfBlock" mapstructure:"numberOfBlock"`
	AppAddress    string               `yaml:"appAddress" mapstructure:"appAddress"`
}

func LoadConfig(path string) (*SnowConfig, error) {
	viperDeafault()
	//viper.AddConfigPath(path)
	//viper.SetConfigName("snow")
	//viper.SetConfigType("yaml")
	//viper.AutomaticEnv()
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	return nil, err
	//}

	config := SnowConfig{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func runDiscovery() (*p2pnetworking.Discovery, error) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	discovery := p2pnetworking.InitDiscovery(r)
	discovery.Router()
	go func() {
		err := r.Run(discovery.Address)
		if err != nil {
			Logger.Fatal(err)
		}
	}()
	return discovery, nil
}

func viperDeafault() {
	viper.SetDefault("p2p.name", "avalanche-consensus")
	viper.SetDefault("p2p.protocolId", "avalanche-consensus/1.0.0")
	viper.SetDefault("p2p.host", "127.0.0.1")
	viper.SetDefault("consensus.k", "3")
	viper.SetDefault("consensus.alphal", "2")
	viper.SetDefault("consensus.beta", "10")
	viper.SetDefault("numberOfNode", "200")
	viper.SetDefault("numberOfBlock", "4")
	viper.SetDefault("appAddress", "127.0.0.1:8080")
}
