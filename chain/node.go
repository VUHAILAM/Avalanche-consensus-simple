package chain

import (
	"avalanche-consensus/consensus"
	"avalanche-consensus/p2pnetworking"
	"context"
)

type Node struct {
	*Blockchain
}

func InitNode(ctx context.Context, config p2pnetworking.Config, snowballConf consensus.Config, discovery *p2pnetworking.Discovery) (*Node, error) {
	node := &Node{}
	blockchain, err := InitBlockchain(ctx, config, snowballConf, discovery)
	if err != nil {
		return nil, err
	}
	node.Blockchain = blockchain
	return node, nil
}
