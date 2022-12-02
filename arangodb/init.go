package arangodb

import (
	"context"
	"ecdsa/config"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/arangodb/go-driver"
	arangodbHttp "github.com/arangodb/go-driver/http"
)

var once sync.Once
var handlers []Handler
var handlerIdx int
var mu sync.Mutex
var conf *config.ArangoDB

type Handler struct {
	db  driver.Database
	ctx context.Context
}

func GetConn() Handler {
	mu.Lock()
	defer mu.Unlock()

	if len(handlers) == 0 {
		log.Fatalf("ArangoDB uninitialized connection\n")
	}

	handlerIdx++
	if handlerIdx == len(handlers) {
		handlerIdx = 0
	}

	return handlers[handlerIdx]
}

func Initialize(c *config.ArangoDB) {
	once.Do(func() {
		conf = c

		// Init handler by http protocol
		switch conf.HttpProtocol {
		case "1.1":
			handlers = make([]Handler, 1)
		case "2":
			handlers = make([]Handler, conf.Connlimit)
		}

		for i := 0; i < len(handlers); i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			db, err := connect(ctx)
			if err != nil {
				log.Fatalf("ArangoDB connect error: %v\n", err)
			}
			handlers[i].db = db
			handlers[i].ctx = ctx
		}

		log.Printf("ArangoDB initialize done\n")
	})
}

func connect(ctx context.Context) (driver.Database, error) {
	var conn driver.Connection

	urls := strings.Split(conf.URLs, ",")
	for _, u := range urls {
		_, err := url.Parse(u)
		if err != nil {
			log.Fatalf("ArangoDB parse url error: %v\n", err)
		}
	}

	switch conf.HttpProtocol {
	case "1.1":
		var err error
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}),
			MaxIdleConns:          0,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		conn, err = arangodbHttp.NewConnection(arangodbHttp.ConnectionConfig{
			Endpoints: urls,
			Transport: transport,
			ConnLimit: conf.Connlimit,
		})

		if err != nil {
			log.Fatalf("ArangoDB http new connection error: %v", err)
			return nil, err
		}
	case "2":
		// TODO
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:                   conn,
		Authentication:               driver.BasicAuthentication(conf.Username, conf.Password),
		SynchronizeEndpointsInterval: 0, // Cluster sync interval, if value > 0 then sync, otherwise no sync.
	})
	if err != nil {
		log.Fatalf("ArangoDB driver new client error: %v", err)
		return nil, err
	}

	if _, err := client.Version(ctx); err != nil {
		log.Fatalf("ArangoDB check version error: %v", err)
		return nil, err
	}

	db, err := client.Database(ctx, conf.Database)
	if err != nil {
		log.Fatalf("ArangoDB get database error: %v", err)
		return nil, err
	}

	return db, nil
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
