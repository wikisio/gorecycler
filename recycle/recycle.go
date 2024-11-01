package recycle

import (
	"sync"
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
	return p.p.Get().(Recycler[BT])
}

func (p *pool[BT]) put(r *recycle[BT]) {
	p.p.Put(r)
}

type recycle[BT any] struct {
	b BT
	p *pool[BT]
}

func (b *recycle[BT]) HandleAndRecycle(cleanFunc func(bt BT) error) error {
	defer func() {
		b.p.cleanFunc(b.b)
		b.p.put(b)
	}()

	return cleanFunc(b.b)
}

func (b *recycle[BT]) Assign(h func(t BT)) {
	h(b.b)
}
