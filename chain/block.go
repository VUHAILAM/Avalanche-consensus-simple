package chain

type DataType int

type Block struct {
	Data      DataType `json:"data"`
	Hash      []byte   `json:"hash"`
	Timestamp int64    `json:"timestamp"`
}

func (b *Block) SetData(data DataType) error {
	b.Data = data
	return nil
}
