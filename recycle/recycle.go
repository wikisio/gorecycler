package recycle

import (
	"sync"
	"sync/atomic"
)

var pools = sync.Map{}
var bts = sync.Map{}

type pool[BT any] struct {
	p         sync.Pool
	cleanFunc func(bt BT)

	empty BT
	t     sync.Map
}

func (p *pool[BT]) get() Recycler[BT] {
	From.Add(1)
	return p.p.Get().(Recycler[BT])
}

func (p *pool[BT]) put(r *recycle[BT]) {
	To.Add(1)
	p.p.Put(r)
}

type recycle[BT any] struct {
	b BT
	p *pool[BT]
}

var To = atomic.Int64{}
var From = atomic.Int64{}

func (b *recycle[BT]) HandleAndRecycle(processBTFunc func(bt BT) error) error {
	defer func() {
		b.p.cleanFunc(b.b)
		b.p.put(b)
	}()

	return processBTFunc(b.b)
}

func (b *recycle[BT]) Assign(h func(t BT)) {
	h(b.b)
}
