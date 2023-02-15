package p2pnetworking

import (
	"avalanche-consensus/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

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
	getBlockDataCallback func(int) (model.DataType, error)
}

func InitPeerClient(
	ctx context.Context,
	cfg Config, discovery *Discovery,
	getBlockDataCb func(int) (model.DataType, error),
) (*Client, error) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	client := &Client{
		conf:                 cfg,
		ginR:                 r,
		discovery:            discovery,
		resty:                resty.New(),
		getBlockDataCallback: getBlockDataCb,
	}

	peer, err := client.InitP2P()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to start the client with host: %s, port: %d", cfg.Host, cfg.Port))
	}

	client.peer = peer
	logrus.Infof("Init P2P Client successfully, host: %s, port: %d", cfg.Host, cfg.Port)

	client.Router()

	err = client.RegisterDiscovery(ctx, peer)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Register done, host: %s, port: %d", cfg.Host, cfg.Port)
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
			logrus.Fatal(err)
		}
	}()
	return p, nil
}

func (c *Client) Router() {
	c.ginR.POST("/get-data-by-index", c.GetBlockDataByIndex)
	c.ginR.GET("/health", c.Health)
}

func (c *Client) Health(ginCtx *gin.Context) {
	ginCtx.JSON(200, fmt.Sprintf("Peer address: %s still ok", c.peer.Address))
}

func (c *Client) RegisterDiscovery(ctx context.Context, peer *model.Peer) error {
	_, err := c.resty.R().SetBody(map[string]interface{}{
		"peer": peer,
	}).Post(fmt.Sprintf("http://%s/register-peer", c.discovery.Address))
	if err != nil {
		logrus.Error(err)
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

func (c *Client) GetBlockDataFromPeer(ctx context.Context, req model.GetBlockDataFromPeerRequest) (model.DataType, error) {
	resp, err := c.resty.R().SetBody(req).Post(fmt.Sprintf("http://%s/get-data-by-index", req.Peer.Address))
	if err != nil {
		return -1, err
	}
	data, err := strconv.Atoi(string(resp.Body()))
	if err != nil {
		return -1, err
	}
	return model.DataType(data), nil
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
		logrus.Error(err)
		return nil, err
	}

	return response.Peers, nil
}
