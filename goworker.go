package gopool

import (
	"context"
	"github.com/geedchin/go-task-pool/util"
	"sync"
	"time"
)

type GoWorker interface {
	start()
	stop()
	Running() bool
	WorkerInfo() WorkerInfo
	doWork(ctx context.Context, f func(ctx context.Context, wi WorkerInfo))
}

type WorkerInfo struct {
	WorkerId string
	IntId    int
}

func newGoWorker(id int) GoWorker {
	wi := WorkerInfo{
		WorkerId: util.GenRandomDigitLowerLetter(32),
		IntId:    id,
	}
	return &goWorker{
		wi:     wi,
		lock:   sync.Mutex{},
		queue:  newTaskQ(100),
		stopCh: make(chan struct{}),
	}
}

type goWorker struct {
	wi      WorkerInfo
	lock    sync.Mutex
	queue   TaskQ
	running bool

	stopLocker sync.Mutex
	stopCh     chan struct{}
}

func (gw *goWorker) Running() bool {
	gw.lock.Lock()
	defer gw.lock.Unlock()
	return gw.running
}

func (gw *goWorker) stop() {
	gw.stopLocker.Lock()
	defer gw.stopLocker.Unlock()
	util.CloseChannel(gw.stopCh)
	gw.queue.Shutdown()
}

func (gw *goWorker) start() {
	gw.lock.Lock()
	defer gw.lock.Unlock()
	if gw.running {
		return
	}
	go func() {
		for {
			f, shutdown := gw.queue.Get()
			if shutdown {
				return
			}
			if f != nil {
				f()
			}
		}
	}()
	gw.running = true
	return
}

func (gw *goWorker) WorkerInfo() WorkerInfo {
	return gw.wi
}

func (gw *goWorker) doWork(ctx context.Context, f func(ctx context.Context, wi WorkerInfo)) {
	gw.queue.Add(func() {
		done := make(chan struct{})
		go func() {
			f(ctx, gw.WorkerInfo())
			done <- struct{}{}
		}()
		select {
		case <-time.NewTimer(time.Minute * 1).C:
			// do nothing
		case <-done:
			// do nothing
		}
	})
}
