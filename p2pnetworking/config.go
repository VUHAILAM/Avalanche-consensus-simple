package p2pnetworking

type Config struct {
	Name       string `json:"name" mapstructure:"NAME"`
	ProtocolID string `json:"protocol_id" mapstructure:"PROTOCOL_ID"`
	Host       string `json:"host" mapstructure:"HOST"`
	Port       int    `json:"port" mapstructure:"PORT"`
}
