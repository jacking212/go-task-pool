package gopool

import (
	"context"
	"sync"
	"testing"
	"time"
)

func dealTaskQ(tq TaskQ) {
	go func() {
		for {
			f, shutdown := tq.Get()
			if shutdown {
				return
			}
			f()
		}
	}()
}

func TestTaskQ_Size(t *testing.T) {
	testData := [][]int{
		{-1, 1},
		{0, 1},
		{1, 1},
		{2, 2},
	}
	for _, datum := range testData {
		tq := newTaskQ(datum[0])
		if tq.Size() == datum[1] {
			continue
		}
		t.Errorf("taskqsizetest input(%d),expect(%d),but get(%d)",
			datum[0], datum[1], tq.Size())
	}

}

// testTaskQAdd
func TestTaskQ_Add(t *testing.T) {
	tq := newTaskQ(10)

	f := func() {
		time.Sleep(time.Second)
	}
	// test f==nil
	tq.Add(nil)
	if tq.QLen() != 0 {
		t.Errorf("testQAdd :input(nil),expect(0),but get %d", tq.QLen())
	}
	// f !=nil
	tq.Add(f)
	if tq.QLen() != 1 {
		t.Errorf("testQAdd :input(nil),expect(1),but get %d", tq.QLen())
	}
	tq.Shutdown()
	//
	tq.Get()
	// retest
	tq.Add(nil)
	if tq.QLen() != 0 {
		t.Errorf("testQAdd :input(nil),expect(0),but get %d", tq.QLen())
	}
	// f !=nil // closed so cant add into q
	tq.Add(f)
	if tq.QLen() != 0 {
		t.Errorf("testQAdd :input(nil),expect(0),but get %d", tq.QLen())
	}
}

func TestTaskQ_Get(t *testing.T) {
	tq := newTaskQ(10)
	f := func() {}
	testData := []struct {
		f        func()
		isNil    bool
		shutdown bool
	}{
		// case 1
		{func() {
			tq.Add(f)
		}, false, false},
		// case 2
		{func() {
			tq.Clear()
			// 1s 后加数据，加之前tq.Get阻塞
			go func() {
				time.Sleep(time.Second)
				tq.Add(f)
			}()
		}, false, false},
		// case 3
		{func() {
			tq.Clear()
			tq.Add(f)
			tq.Shutdown()
		}, false, false},
		// case 4
		{func() {
			tq.Clear()
			tq.Shutdown()
		}, true, true},
	}
	// 1. 正常获取队列已有数据
	// 2. 正常获取，队列无数据，但是后续加数据
	// 3. 获取关闭的q, 但队列未空的数据
	// 4. 获取关闭的q, 队列无数据

	for i, item := range testData {
		if item.f != nil {
			item.f()
		}
		ff, shutdown := tq.Get()
		fIsNil := ff == nil
		if fIsNil != item.isNil || shutdown != item.shutdown {
			t.Errorf("case %d :want(%v %v),but get(%v %v)",
				i+1, item.isNil, item.shutdown, fIsNil, shutdown)
		}
	}
}

func TestTaskQ(t *testing.T) {
	tq := newTaskQ(10)

	dealTaskQ(tq)

	k := 0
	wg := sync.WaitGroup{}
	it := 10
	wg.Add(it)
	f := func(ctx context.Context, wi WorkerInfo) {
		defer wg.Done()
		k++
		t.Log(k)
		time.Sleep(time.Second * 1)
	}
	for i := 0; i < it; i++ {
		tq.Add(func() {
			f(context.TODO(), WorkerInfo{WorkerId: "11"})
		})
	}
	wg.Wait()
}
