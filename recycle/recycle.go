package recycle

import (
	"sync"
)

var pools = sync.Map{}

type pool[BT any] struct {
	p     sync.Pool
	freeF func(bt BT)
}

func (p *pool[BT]) get() Recycle[BT] {
	return p.p.Get().(Recycle[BT])
}

func (p *pool[BT]) put(r Recycle[BT]) {
	p.p.Put(r)
}

func (p *pool[BT]) free(bt BT) {
	p.free(bt)
}

type baseRecycler[BT any] struct {
	b BT
	p *pool[BT]
}

func (b *baseRecycler[BT]) HandleAndRecycle(h func(t BT) error) error {
	defer func() {
		b.p.freeF(b.b)
		b.p.put(b)
	}()

	return h(b.b)
}

func (b *baseRecycler[BT]) Assign(h func(t BT)) {
	h(b.b)
}
