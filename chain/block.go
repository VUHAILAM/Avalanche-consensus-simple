package chain

import "avalanche-consensus/model"

type Block struct {
	Data      model.DataType `json:"data"`
	Hash      []byte         `json:"hash"`
	Timestamp int64          `json:"timestamp"`
}

func (b *Block) SetData(data model.DataType) error {
	b.Data = data
	return nil
}
