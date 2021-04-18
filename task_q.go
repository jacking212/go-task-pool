package gopool

import (
	"github.com/geedchin/go-task-pool/util"
	"sync"
)

type TaskQ interface {
	Add(f func())
	Shutdown()
	Get() (func(), bool)
	Size() int
	QLen() int
	Clear()
}

func newTaskQ(size int) TaskQ {
	if size < 1 {
		size = 1
	}
	return &taskQ{
		stopCh: make(chan struct{}),
		q:      make(chan func(), size),
		size:   size,
	}
}

type taskQ struct {
	stopCh chan struct{}
	q      chan func()
	l      sync.Mutex
	size   int
}

func (t *taskQ) Clear() {
	t.l.Lock()
	defer t.l.Unlock()
	t.q <- nil
	t.q = make(chan func(), t.size)
}

func (t *taskQ) QLen() int {
	t.l.Lock()
	defer t.l.Unlock()
	return len(t.q)
}

func (t *taskQ) Size() int {
	return t.size
}

func (t *taskQ) Add(f func()) {
	if f == nil {
		return
	}
	t.l.Lock()
	defer t.l.Unlock()
	select {
	case <-t.stopCh:
		return
	default:
		t.q <- f
	}
}
func (t *taskQ) Shutdown() {
	t.l.Lock()
	defer t.l.Unlock()
	util.CloseChannel(t.stopCh)
}

func (t *taskQ) Get() (func(), bool) {
	t.l.Lock()
	q := t.q
	t.l.Unlock()
	// 保证队列有数据时优先返回队列数据
	select {
	case f := <-q:
		return f, false
	default:
	}
	// 判断关闭
	select {
	case <-t.stopCh:
		return nil, true
	case f := <-q:
		return f, false
	}
}
