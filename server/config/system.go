package config

type System struct {
	RouterPrefix       string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr               int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	LimitCountIP       int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP        int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseStrictAuth      bool   `mapstructure:"use-strict-auth" json:"use-strict-auth" yaml:"use-strict-auth"`
	DisableAutoMigrate bool   `mapstructure:"disable-auto-migrate" json:"disable-auto-migrate" yaml:"disable-auto-migrate"`
}
