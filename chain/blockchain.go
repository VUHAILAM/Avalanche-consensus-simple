package chain

import (
	"avalanche-consensus/consensus"
	"avalanche-consensus/model"
	"avalanche-consensus/p2pnetworking"
	"context"
	"errors"
	"math/rand"
	"sync"

	"github.com/sirupsen/logrus"
)

type Chain interface {
	Length() int
	Get(int) (*Block, error)
	Add(*Block) error
	Set(int, model.DataType) error
}

type Blockchain struct {
	Blocks        []*Block
	PeerClient    *p2pnetworking.Client
	P2pConfig     p2pnetworking.Config
	SnowbalConfig consensus.Config
	isRunning     bool
	mu            sync.Mutex
}

func InitBlockchain(ctx context.Context, p2pConfig p2pnetworking.Config, consensusConf consensus.Config, discovery *p2pnetworking.Discovery) (*Blockchain, error) {
	blockchain := &Blockchain{
		P2pConfig:     p2pConfig,
		SnowbalConfig: consensusConf,
	}
	blockchain.Blocks = make([]*Block, 0)

	getDataCallback := func(index int) (model.DataType, error) {
		return blockchain.GetBlockData(index)
	}

	client, err := p2pnetworking.InitPeerClient(ctx, p2pConfig, discovery, getDataCallback)
	if err != nil {
		return nil, err
	}
	blockchain.PeerClient = client
	return blockchain, nil
}

func (b *Blockchain) Add(block *Block) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Blocks = append(b.Blocks, block)
	return nil
}

func (b *Blockchain) Length() int {
	return len(b.Blocks)
}

func (b *Blockchain) Get(index int) (*Block, error) {
	if index < 0 || index >= b.Length() {
		return nil, errors.New("index out of range")
	}

	return b.Blocks[index], nil
}

func (b *Blockchain) Set(index int, data model.DataType) error {
	dataBlock, err := b.Get(index)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	dataBlock.Data = data

	return nil
}

func (b *Blockchain) GetBlockData(index int) (model.DataType, error) {
	dataBlock, err := b.Get(index)
	if err != nil {
		return -1, err
	}

	return dataBlock.Data, nil
}

func (b *Blockchain) GetBlockDataFromPeer(ctx context.Context, peer *model.Peer, index int) (model.DataType, error) {
	req := model.GetBlockDataFromPeerRequest{
		Peer:  peer,
		Index: index,
	}
	data, err := b.PeerClient.GetBlockDataFromPeer(ctx, req)
	if err != nil {
		return -1, err
	}

	return data, nil
}

func (b *Blockchain) GetDataFromKRandomBlock(ctx context.Context, index int, k int) ([]model.DataType, error) {
	peers, err := b.PeerClient.Peers()
	if err != nil {
		return nil, err
	}

	lenPeers := len(peers)
	preferencesFromPeers := make([]model.DataType, 0)

	cnt := 0
	for _, i := range rand.Perm(lenPeers) {
		if peers[i] == nil {
			continue
		}
		preference, err := b.GetBlockDataFromPeer(ctx, peers[i], index)
		if err != nil {
			logrus.Error(err)
			continue
		}

		preferencesFromPeers = append(preferencesFromPeers, preference)
		cnt++

		if cnt >= k {
			break
		}
	}
	return preferencesFromPeers, nil
}

func (b *Blockchain) Sync(ctx context.Context) error {
	if b.isRunning {
		return nil
	}
	b.isRunning = true
	for i, block := range b.Blocks {
		snowballConsensus, err := consensus.NewConcensus(b.SnowbalConfig, block.Data)
		if err != nil {
			return err
		}
		getKDataCb := func(k int) ([]model.DataType, error) {
			return b.GetDataFromKRandomBlock(ctx, i, k)
		}
		setDataCb := func(data model.DataType) error {
			return block.SetData(data)
		}
		err = snowballConsensus.Run(ctx, setDataCb, getKDataCb)
		if err != nil {
			return err
		}
	}
	b.isRunning = false
	return nil
}
