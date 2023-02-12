package p2pnetworking

import (
	"avalanche-consensus/model"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	DiscoveryHost = "0.0.0.0"
	DiscoveryPort = 8080
)

type Discovery struct {
	mu          sync.Mutex
	Address     string
	Peers       []*model.Peer
	ginR        *gin.Engine
	restyClient *resty.Client
}

func InitDiscovery() *Discovery {
	ginEng := gin.Default()
	address := fmt.Sprintf("%s:%d", DiscoveryHost, DiscoveryPort)
	return &Discovery{
		Address:     address,
		Peers:       make([]*model.Peer, 0),
		ginR:        ginEng,
		restyClient: resty.New(),
	}
}

func (d *Discovery) Router() {
	d.ginR.POST("/register-peer", d.RegisterPeer)
	d.ginR.GET("/peers", d.GetPeers)
}

func (d *Discovery) RegisterPeer(c *gin.Context) {
	var req model.RegisterPeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, nil)
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
	c.JSON(200, model.RegisterPeerResponse{
		Peers: d.Peers,
	})
}

func (d *Discovery) GetPeers(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"peers": d.Peers,
	})
}

func (d *Discovery) HealthCheckPeers() error {
	fmt.Println("healthy check start")
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
