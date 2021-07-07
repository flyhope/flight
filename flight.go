package flight

import (
	"sync"
	"sync/atomic"
)

type result struct {
	val interface{}
	err error
}

type call struct {
	num     uint32
	waitNum uint32
	channel chan *result
}

type Group struct {
	m sync.Map
}

func (g *Group) Do(max uint32, key string, callback func() (interface{}, error)) (interface{}, error) {
	// load map value
	var c *call
	if storeCall, ok := g.m.Load(key); ok {
		c = storeCall.(*call)
	} else {
		// set map value
		c = &call{channel: make(chan *result)}
		if storeCall, ok := g.m.LoadOrStore(key, c); ok {
			c = storeCall.(*call)
		}
	}

	// wait value
	num := atomic.AddUint32(&c.num, 1)
	if num > max {
		atomic.AddUint32(&c.waitNum, 1)
		if ret, ok := <-c.channel; ok {
			return ret.val, ret.err
		}
		return callback()
	}

	val, err := callback()
	g.m.Delete(key)

	ret := &result{val: val, err: err}
	for i := uint32(0); i < c.waitNum; i++ {
		select {
		case c.channel <- ret:
		default:
		}
	}
	close(c.channel)

	return val, err
}
