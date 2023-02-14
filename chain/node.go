package chain

import (
	"avalanche-consensus/consensus"
	"avalanche-consensus/p2pnetworking"
	"context"

	"github.com/sirupsen/logrus"
)

type Node struct {
	*Blockchain
}

func InitNode(ctx context.Context, config p2pnetworking.Config, snowballConf consensus.Config, discovery *p2pnetworking.Discovery) (*Node, error) {
	node := &Node{}
	blockchain, err := InitBlockchain(ctx, config, snowballConf, discovery)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	node.Blockchain = blockchain
	return node, nil
}
