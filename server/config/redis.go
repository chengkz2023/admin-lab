package config

type Redis struct {
	Enable       bool     `mapstructure:"enable" json:"enable" yaml:"enable"`
	Name         string   `mapstructure:"name" json:"name" yaml:"name"`
	Addr         string   `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password     string   `mapstructure:"password" json:"password" yaml:"password"`
	DB           int      `mapstructure:"db" json:"db" yaml:"db"`
	UseCluster   bool     `mapstructure:"useCluster" json:"useCluster" yaml:"useCluster"`
	ClusterAddrs []string `mapstructure:"clusterAddrs" json:"clusterAddrs" yaml:"clusterAddrs"`
}
