package main

import (
	"avalanche-consensus/chain"
	"avalanche-consensus/p2pnetworking"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})

	discovery, err := runDiscovery()
	if err != nil {
		Logger.Fatal(err)
	}

	config, err := LoadConfig(".")
	if err != nil {
		Logger.Fatal(err)
	}
	var wg sync.WaitGroup

	for i := 0; i < config.NumberOfNode; i++ {
		wg.Add(1)
		go SetupNode(&wg, i, discovery, config)
	}
	wg.Wait()
}

func SetupNode(wg *sync.WaitGroup, i int, discovery *p2pnetworking.Discovery, conf *SnowConfig) {
	defer wg.Done()
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

	for j := 0; j < conf.NumberOfBlock; i++ {
		r := rand.Intn(int(float32(2) * 1.5))
		block := &chain.Block{}
		if r < 2 {
			block.SetData(chain.DataType(r))
		} else {
			block.SetData(chain.DataType(j))
		}
		node.Add(block)
	}
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
	blockChainState := ""
	for _, b := range node.Blocks {
		blockChainState += fmt.Sprintf("%d ", int(b.Data))
	}

	Logger.Infof("client: %d, block: %s", i, blockChainState)
}

func runDiscovery() (*p2pnetworking.Discovery, error) {
	discovery := p2pnetworking.InitDiscovery()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	discovery.Router()
	go func() {
		err := r.Run(discovery.Address)
		if err != nil {
			Logger.Fatal(err)
		}
	}()
	return discovery, nil
}
