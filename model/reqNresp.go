package model

type RegisterPeerRequest struct {
	Peer *Peer `json:"peer"`
}

type RegisterPeerResponse struct {
	Peers []*Peer `json:"peers"`
}

type GetBlockDataByIndexRequest struct {
	Index int `json:"index"`
}
