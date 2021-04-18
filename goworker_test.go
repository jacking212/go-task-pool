package gopool

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGoWorkerStartStop(t *testing.T) {

	wk := newGoWorker(0)
	if wkR := wk.Running(); wkR != false {
		t.Errorf("want %v,but get %v", false, wkR)
	}
	wk.start()
	if wkR := wk.Running(); wkR != true {
		t.Errorf("want %v,but get %v", true, wkR)
	}
	wk.start()
	if wkR := wk.Running(); wkR != true {
		t.Errorf("want %v,but get %v", true, wkR)
	}
	startDo := make(chan struct{})
	done := make(chan struct{})
	go func() { <-done }()
	// stop
	f := func(ctx context.Context, wi WorkerInfo) {
		startDo <- struct{}{}
		time.Sleep(time.Second * 3)
		done <- struct{}{}
	}
	wk.doWork(context.TODO(), f)
	wk.doWork(context.TODO(), f)
	<-startDo
	// 第一个任务开始后，stop wk，构造wk停止，但是仍有task场景，
	// 此时仍会执行task，startDo同步判断
	wk.stop()
	timeout := time.NewTimer(time.Second * 8)
	select {
	case <-startDo:
	case <-timeout.C:
		t.Errorf("wk停止后仍有task在队列，但未执行")
	}
	<-done
}

func TestGoWorker_doWorker(t *testing.T) {

	wk := newGoWorker(0)

	wk.start()

	k := 0
	wg := sync.WaitGroup{}
	it := 10
	wg.Add(it)
	f := func(ctx context.Context, wi WorkerInfo) {
		defer wg.Done()
		k++
		t.Log(k)
		time.Sleep(time.Second)
	}
	for i := 0; i < it; i++ {
		wk.doWork(context.TODO(), f)
	}
	wg.Wait()
}
