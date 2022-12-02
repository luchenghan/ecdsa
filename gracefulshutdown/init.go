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
var wgs []*sync.WaitGroup
var osShutdownChan *chan os.Signal

type srvShutdownChan chan srvSignal
type srvSignal struct{}

// Initialize gracefulshutdown
func Init(srvNum int) {
	osShutdownChan = new(chan os.Signal)
	// kill no param => syscall.SIGTERM
	// kill -2 => syscall.SIGINT, ctrl + c, os.Interrupt
	// kill -9 => syscall.SIGKILL
	signal.Notify(*osShutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	cancels := make([]context.CancelFunc, srvNum)
	ctxs = make([]context.Context, srvNum)
	wgs = make([]*sync.WaitGroup, srvNum)

	for i := 0; i < srvNum; i++ {
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
		wgs[i] = new(sync.WaitGroup)
	}

	// Listen os signal & wait group done
	go func() {
		<-*osShutdownChan
		signal.Stop(*osShutdownChan)
		for i := 0; i < srvNum; i++ {
			cancels[i]()
			wgs[i].Wait()
		}
		log.Fatal("GracefulShutdown Done")
	}()
}

func GetContext(servicelevel uint8) (context.Context, srvShutdownChan) {
	c := make(srvShutdownChan)
	wgs[servicelevel].Add(1)
	go func(c srvShutdownChan, wg *sync.WaitGroup) {
		<-c
		wg.Done()
	}(c, wgs[servicelevel])
	return ctxs[servicelevel], c
}

func Shutdown() {
	*osShutdownChan <- os.Interrupt
}
