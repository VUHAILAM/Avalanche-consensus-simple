package p2pnetworking

type Config struct {
	Name       string `json:"name"`
	ProtocolID string `json:"protocol_id"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
}
