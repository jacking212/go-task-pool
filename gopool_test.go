package gopool

import (
	"context"
	"fmt"
	"github.com/geedchin/go-task-pool/util"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

func TestGoPool_DoWorker(t *testing.T) {

	go func() {
		//http.ListenAndServe(":8080", nil)
	}()
	delta := 20
	deltaBatchSize := 30
	poolSize := 31

	gp := NewGoPool(poolSize)

	wg := sync.WaitGroup{}
	tt := NewTaskTotal()
	wg.Add(delta)
	for i := 0; i < delta; i++ {
		go func() {
			kid := util.GenRandomDigitLowerLetter(32)
			idx := 0
			lock := sync.Mutex{}
			for j := 0; j < deltaBatchSize; j++ {
				wg.Add(1)
				gp.DoWorker(context.TODO(), kid, func(ctx context.Context, workerInfo WorkerInfo) {
					defer wg.Done()
					//t.Log(fmt.Sprintf("%s,%d,%d,%s", kid, idx, workerInfo.IntId, workerInfo.WorkerId))
					tt.DoWithLock(func() {
						item := tt.tasks[kid]
						item.idxs = append(item.idxs, idx)
						item.workerIntIds = append(item.workerIntIds, workerInfo.IntId)
						item.workerIds = append(item.workerIds, workerInfo.WorkerId)
						tt.tasks[kid] = item
					})
					lock.Lock()
					idx++
					lock.Unlock()
					time.Sleep(time.Second)
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// check
	tt.DoWithLock(func() {
		for key, task := range tt.tasks {
			l1 := len(task.idxs)
			l2 := len(task.workerIntIds)
			l3 := len(task.workerIds)
			if !(l1 == l2 && l2 == l3 && l1 == deltaBatchSize) {
				t.Fatalf("!(l1 == l2 && l2 == l3 && l1 == deltaBatchSize) ")
			}
			//sort.Ints(task.idxs)
			if task.idxs[0] != 0 {
				t.Fatalf("idx[0] is not 0")
			}
			for i := 1; i < len(task.idxs); i++ {
				if task.idxs[i]-task.idxs[i-1] != 1 {
					t.Fatalf("task.idxs[i]-task.idxs[i-1] != 1")
				}
			}
			for i := 1; i < len(task.workerIds); i++ {
				if task.workerIds[i] != task.workerIds[i-1] {
					t.Fatalf("task.workerIds[i]!=task.workerIds[i-1]")
				}
			}
			for i := 1; i < len(task.workerIntIds); i++ {
				if task.workerIntIds[i] != task.workerIntIds[i-1] {
					t.Fatalf("task.workerIntIds[i]!=task.workerIntIds[i-1]")
				}
			}
			t.Log(fmt.Sprintf("%s\t%v\t%d\t%s", key, task.idxs, task.workerIntIds[0], task.workerIds[0]))
		}
	})
	// total task dis
	tt.DoWithLock(func() {
		total := map[string]int{}
		for _, task := range tt.tasks {
			total[task.workerIds[0]] += 1
		}
		for key, val := range total {
			t.Logf("%s\t%d", key, val)
		}
	})
}

func NewTaskTotal() *taskTotal {
	return &taskTotal{
		tasks: map[string]taskItem{},
	}
}

type taskTotal struct {
	lock  sync.Mutex
	tasks map[string]taskItem
}
type taskItem struct {
	idxs         []int
	workerIntIds []int
	workerIds    []string
}

func (t *taskTotal) DoWithLock(f func()) {
	t.lock.Lock()
	defer t.lock.Unlock()
	f()
}
