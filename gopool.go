package gopool

import (
	"context"
	"github.com/geedchin/go-task-pool/util"
	"sync"
)

type GoPool interface {
	Running() bool
	StopAll()
	DoWorker(ctx context.Context, key string, f func(ctx context.Context, workerInfo WorkerInfo))
}

func NewGoPool(poolSize int) GoPool {
	if poolSize <= 0 {
		poolSize = 1
	} else {
		if poolSize > 200 {
			poolSize = 200
		}
		poolSize = nextPrim(poolSize)
	}
	return &goPool{
		size:   poolSize,
		wkMap:  make(map[int]GoWorker, poolSize),
		stopCh: make(chan struct{}),
	}
}

type goPool struct {
	size     int
	wkMap    map[int]GoWorker
	wkLocker sync.Mutex

	stopCh     chan struct{}
	stopLocker sync.Mutex
}

func (g *goPool) Running() bool {
	g.stopLocker.Lock()
	defer g.stopLocker.Unlock()
	return !util.ChannelIsClosed(g.stopCh)
}

func (g *goPool) closeAllWorker() {
	g.wkLocker.Lock()
	defer g.wkLocker.Unlock()
	for _, worker := range g.wkMap {
		worker.stop()
	}
}

// 重入关闭通道
func (g *goPool) StopAll() {
	g.stopLocker.Lock()
	defer g.stopLocker.Unlock()
	util.CloseChannel(g.stopCh)
	// close worker
	g.closeAllWorker()
}

func (g *goPool) DoWorker(ctx context.Context, key string, f func(ctx context.Context, workerInfo WorkerInfo)) {
	wk := g.getWorker(hash(key))
	if g.Running() {
		wk.doWork(ctx, f)
	}
}

func (g *goPool) getWorker(hashId int) GoWorker {
	g.wkLocker.Lock()
	defer g.wkLocker.Unlock()
	id := hashId % g.size
	wk, ok := g.wkMap[id]
	if ok {
		return wk
	}
	wk = newGoWorker(id)
	g.wkMap[id] = wk
	wk.start()
	return wk
}
