package p2pnetworking

import (
	"avalanche-consensus/chain"
	"avalanche-consensus/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Client struct {
	conf                 Config
	peer                 *model.Peer
	ginR                 *gin.Engine
	discovery            *Discovery
	resty                *resty.Client
	getBlockDataCallback func(int) (chain.DataType, error)
}

func InitPeerClient(
	ctx context.Context,
	cfg Config, discovery *Discovery,
	getBlockDataCb func(int) (chain.DataType, error),
) (*Client, error) {
	r := gin.Default()
	client := &Client{
		conf:                 cfg,
		ginR:                 r,
		discovery:            discovery,
		resty:                resty.New(),
		getBlockDataCallback: getBlockDataCb,
	}

	p2pClient, err := client.InitP2P()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to start the client with host: %s, port: %d", cfg.Host, cfg.Port))
	}

	client.peer = p2pClient

	err = client.RegisterDiscovery(ctx, p2pClient)
	if err != nil {
		return nil, err
	}
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
	c.ginR.POST("/get-data-by-index", c.GetBlockDataByIndex)
	c.ginR.GET("/health", c.Health)
}

func (c *Client) Health(ginCtx *gin.Context) {
	ginCtx.JSON(200, "ok")
}

func (c *Client) RegisterDiscovery(ctx context.Context, peer *model.Peer) error {
	resp, err := c.resty.R().SetBody(map[string]interface{}{
		"peer": peer,
	}).Post(fmt.Sprintf("http://%s/register-peer", c.discovery.Address))
	if err != nil {
		return err
	}
	response := model.RegisterPeerResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetBlockDataByIndex(ginCtx *gin.Context) {
	req := model.GetBlockDataByIndexRequest{}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		ginCtx.JSON(400, "Can not read request")
		return
	}

	blockData, err := c.getBlockDataCallback(req.Index)
	if err != nil {
		ginCtx.JSON(400, "Can not get block data")
		return
	}
	ginCtx.JSON(200, blockData)
}

func (c *Client) GetBlockDataFromPeer(ctx context.Context, req model.GetBlockDataFromPeerRequest) (chain.DataType, error) {
	resp, err := c.resty.R().SetBody(req).Post(fmt.Sprintf("http://%s/get-data-by-index", req.Peer.Address))
	if err != nil {
		return -1, err
	}
	data, err := strconv.Atoi(string(resp.Body()))
	if err != nil {
		return -1, err
	}
	return chain.DataType(data), nil
}

func (c *Client) Peers() ([]*model.Peer, error) {
	resp, err := c.resty.R().Get(fmt.Sprintf("http://%s/peers", c.discovery.Address))
	if err != nil {
		return nil, err
	}
	var response struct {
		Peers []*model.Peer `json:"peers"`
	}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response.Peers, nil
}
