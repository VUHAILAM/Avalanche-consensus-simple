package p2pnetworking

import (
	"avalanche-consensus/model"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	DiscoveryHost = "127.0.0.1"
	DiscoveryPort = 8080
)

type Discovery struct {
	mu          sync.Mutex
	Address     string
	Peers       []*model.Peer
	ginR        *gin.Engine
	restyClient *resty.Client
}

func InitDiscovery(r *gin.Engine) *Discovery {
	address := fmt.Sprintf("%s:%d", DiscoveryHost, DiscoveryPort)
	return &Discovery{
		Address:     address,
		Peers:       make([]*model.Peer, 0),
		ginR:        r,
		restyClient: resty.New(),
	}
}

func (d *Discovery) Router() {
	d.ginR.POST("/register-peer", d.RegisterPeer)
	d.ginR.GET("/peers", d.GetPeers)
}

func (d *Discovery) RegisterPeer(ginCtx *gin.Context) {
	var req model.RegisterPeerRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		logrus.Error("Can not bind json", err)
		ginCtx.JSON(400, nil)
	}
	isAvailablePeer := true
	for _, peer := range d.Peers {
		if peer.ID == req.Peer.ID {
			isAvailablePeer = false
		}
	}
	if isAvailablePeer {
		d.Peers = append(d.Peers, req.Peer)
	}
	ginCtx.JSON(200, model.RegisterPeerResponse{
		Peers: d.Peers,
	})
}

func (d *Discovery) GetPeers(ginCtx *gin.Context) {
	ginCtx.JSON(200, map[string]interface{}{
		"peers": d.Peers,
	})
}

func (d *Discovery) HealthCheckPeers() error {
	logrus.Infoln("healthy check start")
	var wg sync.WaitGroup
	peers := make([]*model.Peer, 0)

	for _, p := range d.Peers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := d.restyClient.R().Get(fmt.Sprintf("http://%s/health", p.Address))
			if err != nil {
				return
			}
			peers = append(peers, p)
		}()
	}
	wg.Wait()
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Peers = peers

	return nil
}
