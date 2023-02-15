package p2pnetworking

type Config struct {
	Name       string `yaml:"name" mapstructure:"name"`
	ProtocolID string `yaml:"protocolId" mapstructure:"protocolId"`
	Host       string `yaml:"host" mapstructure:"host"`
	Port       int    `yaml:"port" mapstructure:"port"`
}
