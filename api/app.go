package api

import (
	"avalanche-consensus/chain"
	"avalanche-consensus/consensus"
	"avalanche-consensus/model"
	"avalanche-consensus/p2pnetworking"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/phayes/freeport"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const (
	appHost = "127.0.0.1"
	appPort = 8080
)

type App struct {
	discovery       *p2pnetworking.Discovery
	p2pConfig       p2pnetworking.Config
	consensusConfig consensus.Config
	numberOfBlock   int
	address         string
}

func InitApp(discovery *p2pnetworking.Discovery, p2pConf p2pnetworking.Config, consensusConf consensus.Config, numberOfNode int, address string) {
	ginR := gin.Default()
	gin.SetMode(gin.DebugMode)

	go func() {
		err := ginR.Run(address)
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	app := App{
		discovery:       discovery,
		p2pConfig:       p2pConf,
		consensusConfig: consensusConf,
		numberOfBlock:   numberOfNode,
		address:         address,
	}
	appRouter := ginR.Group("/api/v1")

	appRouter.GET("/health", app.Health)
	appRouter.POST("/create", app.CreateNode)
}

func (a *App) Health(ginCtx *gin.Context) {
	err := a.discovery.HealthCheckPeers()
	if err != nil {
		logrus.Error(err)
		ginCtx.JSON(400, nil)
		return
	}

	ginCtx.JSON(200, "ok")
}

func (a *App) CreateNode(ginCtx *gin.Context) {
	ctx := context.Background()
	freePort, err := freeport.GetFreePort()
	if err != nil {
		ginCtx.JSON(500, "Can not find free port")
		return
	}

	a.p2pConfig.Port = freePort

	node, err := chain.InitNode(ctx, a.p2pConfig, a.consensusConfig, a.discovery)
	if err != nil {
		ginCtx.JSON(500, "Can not Init Node")
		return
	}
	time.Sleep(1 * time.Second)

	for j := 0; j < a.numberOfBlock; j++ {
		r := rand.Intn(int(float32(2) * 1.5))
		block := &chain.Block{}
		if r < 2 {
			block.SetData(model.DataType(r))
		} else {
			block.SetData(model.DataType(j))
		}
		node.Add(block)
	}
	beforeBlockChainState := ""
	for _, b := range node.Blocks {
		beforeBlockChainState += fmt.Sprintf("%d ", int(b.Data))
	}
	logrus.Infof("Before sync, data of new node is %s", beforeBlockChainState)
	err = node.Sync(ctx)
	if err != nil {
		ginCtx.JSON(500, "Can not sync Node")
		return
	}
	blockChainState := ""
	for _, b := range node.Blocks {
		blockChainState += fmt.Sprintf("%d ", int(b.Data))
	}

	logrus.Infof("new node block: %s", blockChainState)
	ginCtx.JSON(200, "Create and Sync success")
}
