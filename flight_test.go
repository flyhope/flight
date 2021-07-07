package flight

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {

	sleepTime := time.Millisecond * 100
	g := &Group{}

	var callNum uint32
	timeStart := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, _ := g.Do(1, "test", func() (interface{}, error) {
				atomic.AddUint32(&callNum, 1)
				time.Sleep(sleepTime)
				return "test", nil
			})

			t.Log(val)
		}()
	}
	wg.Wait()

	timeUse := time.Now().Sub(timeStart)
	if timeUse > sleepTime * 2 {
		t.Fatal("time use so long", timeUse.String())
	}

	if callNum != 1 {
		t.Fatal("call func num is not 1", callNum)
	}
}
