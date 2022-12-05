package gracefulshutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ctxs []context.Context
var cancels []context.CancelFunc
var wgs []*sync.WaitGroup
var osShutdownChan chan os.Signal

type SrvShutdownChan chan srvSignal
type srvSignal struct{}

// Initialize gracefulshutdown
func Init(srvNum int) {
	osShutdownChan = make(chan os.Signal, 1)
	// kill no param => syscall.SIGTERM
	// kill -2 => syscall.SIGINT, ctrl + c, os.Interrupt
	// kill -9 => syscall.SIGKILL
	signal.Notify(osShutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	cancels = make([]context.CancelFunc, srvNum)
	ctxs = make([]context.Context, srvNum)
	wgs = make([]*sync.WaitGroup, srvNum)

	for i := 0; i < srvNum; i++ {
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
		wgs[i] = new(sync.WaitGroup)
	}

	// Listen os signal & wait group done
	go func() {
		<-osShutdownChan
		signal.Stop(osShutdownChan)
		for i := 0; i < srvNum; i++ {
			cancels[i]()
			wgs[i].Wait()
		}
		log.Fatal("GracefulShutdown Done")
	}()
}

func SetContext(servicelevel uint8, ctx context.Context, cancel context.CancelFunc) SrvShutdownChan {
	srvChan := make(SrvShutdownChan)
	ctxs = append(ctxs, ctx)
	cancels = append(cancels, cancel)
	wgs[servicelevel].Add(1)

	go func(c SrvShutdownChan, wg *sync.WaitGroup) {
		<-c
		wg.Done()
	}(srvChan, wgs[servicelevel])

	return srvChan
}

func GetContext(servicelevel uint8) (context.Context, SrvShutdownChan) {
	c := make(SrvShutdownChan)
	wgs[servicelevel].Add(1)

	go func(c SrvShutdownChan, wg *sync.WaitGroup) {
		<-c
		wg.Done()
	}(c, wgs[servicelevel])

	return ctxs[servicelevel], c
}

func Shutdown() {
	osShutdownChan <- os.Interrupt
}
