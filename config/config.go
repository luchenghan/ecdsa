package config

import "time"

type ArangoDB struct {
	Address       string        `yaml:"address,omitempty"`
	Database      string        `yaml:"database,omitempty"`
	Connlimit     int           `yaml:"connlimit,omitempty"`
	Username      string        `yaml:"username,omitempty"`
	Password      string        `yaml:"password,omitempty"`
	RetryCount    int           `yaml:"retryCount,omitempty"`
	RetryInterval time.Duration `yaml:"retryInterval,omitempty"`
	HttpProtocol  string        `yaml:"httpProtocol,omitempty"`
}