package flight

import (
	"sync"
	"sync/atomic"
)

type call struct {
	num  uint32
	val  interface{}
	err  error
	cond sync.Cond
}

type Group struct {
	m sync.Map
}

func (g *Group) Do (max uint32, key string, callback func() (interface{}, error)) (interface{}, error)  {
	// load map value
	var c *call
	if storeCall, ok := g.m.Load(key); ok {
		c = storeCall.(*call)
	} else {
		// set map value
		c = &call{}
		if storeCall, ok := g.m.LoadOrStore(key, c); ok {
			c = storeCall.(*call)
		}
	}

	// wait value
	num := atomic.AddUint32(&c.num, 1)
	if num > max {
		c.cond.Wait()
		return c.val, c.err
	}

	val, err := callback()
	c.val = val
	c.err = err
	g.m.Delete(key)

	// todo map after Broadcast for wait
	c.cond.Broadcast()

	return val, err
}
