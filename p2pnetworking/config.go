package p2pnetworking

type Config struct {
	Name       string `yaml:"name" mapstructure:"name" default:"avalanche-consensus"`
	ProtocolID string `yaml:"protocolId" mapstructure:"protocolId" default:"avalanche-consensus/1.0.0"`
	Host       string `yaml:"host" mapstructure:"host" default:"127.0.0.1"`
	Port       int    `yaml:"port" mapstructure:"port"`
}
