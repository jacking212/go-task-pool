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

	gp := NewGoPool(31)

	wg := sync.WaitGroup{}

	delta := 15
	wg.Add(delta)
	for i := 0; i < delta; i++ {
		go func() {
			kid := util.GenRandomDigitLowerLetter(32)
			idx := 0
			lock := sync.Mutex{}
			for j := 0; j < 10; j++ {
				wg.Add(1)
				gp.DoWorker(context.TODO(), kid, func(ctx context.Context, workerInfo WorkerInfo) {
					defer wg.Done()
					t.Log(fmt.Sprintf("%s,%d,%d,%s", kid, idx, workerInfo.IntId, workerInfo.WorkerId))
					lock.Lock()
					idx++
					lock.Unlock()
					time.Sleep(time.Second / 10)
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
