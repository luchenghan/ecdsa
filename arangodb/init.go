package arangodb

import (
	"context"
	"ecdsa/config"
	"fmt"
	"sync"

	"github.com/arangodb/go-driver"
)

var once sync.Once
var handlers []*Handler
var handlerIdx int
var mu sync.Mutex

type Handler struct {
	db  driver.Database
	ctx context.Context
}

func GetConn() *Handler {
	mu.Lock()
	defer mu.Unlock()

	if len(handlers) == 0 {
		fmt.Println("ArangoDB uninitialized connection")
		return nil
	}

	handlerIdx++
	if handlerIdx == len(handlers) {
		handlerIdx = 0
	}

	return handlers[handlerIdx]
}

func Initialize(conf *config.ArangoDB) {
	once.Do(func() {
		switch conf.HttpProtocol {
		case "1.1":
			handlers = make([]*Handler, 1)
		case "2":
			handlers = make([]*Handler, conf.Connlimit)
		}

		for i := 0; i < len(handlers); i++ {
			client, err := connect(context.Background(), conf)
			if err != nil {
				panic(err)
			}
			handlers[i].db = client
		}

		fmt.Println("ArangoDB Initialize Done")
	})
}

func connect(ctx context.Context, conf *config.ArangoDB) (driver.Database, error) {
	// var conn driver.Database
}
