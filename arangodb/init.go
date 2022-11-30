package arangodb

import (
	"context"
	"sync"
	"time"

	"github.com/arangodb/go-driver"
)

var once sync.Once

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

type Handler struct {
	db  driver.Database
	ctx context.Context
}

func Initialize(config ArangoDB) {
	once.Do(func() {

	})
}
