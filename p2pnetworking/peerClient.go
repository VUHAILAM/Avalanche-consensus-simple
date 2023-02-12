package p2pnetworking

import (
	"avalanche-consensus/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Client struct {
	conf      Config
	peer      *model.Peer
	ginR      *gin.Engine
	discovery *Discovery
	resty     *resty.Client
}

func InitPeerClient(ctx context.Context, cfg Config, discovery *Discovery) (*Client, error) {
	r := gin.Default()
	client := &Client{
		conf:      cfg,
		ginR:      r,
		discovery: discovery,
		resty:     resty.New(),
	}

	p2pClient, err := client.InitP2P()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to start the client with host: %s, port: %d", cfg.Host, cfg.Port))
	}

	client.peer = p2pClient
	return client, nil
}

func (c *Client) InitP2P() (*model.Peer, error) {
	address := fmt.Sprintf("%s:%d", c.conf.Host, c.conf.Port)
	p := &model.Peer{
		Address: address,
		ID:      uuid.New().String(),
	}
	go func() {
		err := c.ginR.Run(address)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return p, nil
}

func (c *Client) Router() {
	//c.ginR.POST("/get-data-by-index", c.GetDataByIndex)
	c.ginR.GET("/health", c.Health)
}

func (c *Client) Health(r *gin.Context) {
	r.JSON(200, "ok")
}

func (c *Client) RegisterDiscovery(ctx context.Context, peer *model.Peer) ([]*model.Peer, error) {
	resp, err := c.resty.R().SetBody(map[string]interface{}{
		"peer": peer,
	}).Post(fmt.Sprintf("http://%s/register-peer", c.discovery.Address))
	if err != nil {
		return nil, err
	}
	response := model.RegisterPeerResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}
	return response.Peers, nil
}
