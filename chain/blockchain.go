package chain

import (
	"avalanche-consensus/p2pnetworking"
	"context"
	"sync"
)

type Chain interface {
}

type Blockchain struct {
	Blocks     []*Block
	PeerClient *p2pnetworking.Client
	P2pConfig  p2pnetworking.Config
	isRunning  bool
	mu         sync.Mutex
}

func InitBlockchain(ctx context.Context, p2pConfig p2pnetworking.Config, discovery *p2pnetworking.Discovery) (*Blockchain, error) {
	blockchain := &Blockchain{}
	blockchain.Blocks = make([]*Block, 0)

	return blockchain, nil
}
