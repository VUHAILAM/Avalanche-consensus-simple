package chain

import (
	"avalanche-consensus/consensus"
	"avalanche-consensus/model"
	"avalanche-consensus/p2pnetworking"
	"context"
	"os"
	"testing"

	"github.com/phayes/freeport"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var p2pConfig p2pnetworking.Config
var consensusConfig consensus.Config
var discovery *p2pnetworking.Discovery

func TestMain(m *testing.M) {
	r := gin.New()
	gin.SetMode(gin.DebugMode)
	discovery = p2pnetworking.InitDiscovery(r)
	discovery.Router()
	go func() {
		err := r.Run(discovery.Address)
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	freePort, err := freeport.GetFreePort()
	if err != nil {
		logrus.Fatal(err)
	}
	p2pConfig = p2pnetworking.Config{
		Host: "127.0.0.1",
		Port: freePort,
	}
	consensusConfig = consensus.Config{
		K:      3,
		Alphal: 2,
		Beta:   10,
	}
	mt := m.Run()
	os.Exit(mt)
}

func TestInitBlockchain(t *testing.T) {
	blockchain, err := InitBlockchain(context.Background(), p2pConfig, consensusConfig, discovery)
	assert.NotNil(t, blockchain)
	assert.NoError(t, err)
}

func TestBlockchain_Add(t *testing.T) {
	blockchain, _ := InitBlockchain(context.Background(), p2pConfig, consensusConfig, discovery)
	block := Block{Data: 1}
	blockchain.Add(&block)

	assert.Equal(t, 1, len(blockchain.Blocks))
}

func TestBlockchain_GetBlockData(t *testing.T) {
	blockchain, _ := InitBlockchain(context.Background(), p2pConfig, consensusConfig, discovery)
	block1 := Block{Data: 1}
	block2 := Block{Data: 2}
	blockchain.Add(&block1)
	blockchain.Add(&block2)

	canNotCase, err := blockchain.GetBlockData(-1)
	assert.Equal(t, model.DataType(-1), canNotCase)
	assert.NotNil(t, err)

	data, err := blockchain.GetBlockData(1)
	assert.Equal(t, model.DataType(2), data)
	assert.NoError(t, err)
}
