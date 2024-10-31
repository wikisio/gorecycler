package recycle

import "reflect"

type Recycle[BT any] interface {
	HandleAndRecycle(h func(t BT) error) error
	Assign(a func(t BT))
}

func Get[BT any, T any]() Recycle[BT] {
	p, _ := pools.Load(reflect.TypeFor[T]())
	return p.(*pool[BT]).get()
}

func RegisterPool[BT any, T any]() {
	RegisterPoolWithCleaner[BT, T](nil)
}

// RegisterPoolWithCleaner registers a pool to allocates instance of T
// T is the target struct to be used. *T should implement the BT interface.
// if T needs to do close or clean operations, do it in cleanFunc, otherwise nil
func RegisterPoolWithCleaner[BT any, T any](cleanFunc func(bt BT)) {
	_, ok := pools.Load(reflect.TypeFor[T]())
	if ok {
		return
	}

	p := &pool[BT]{}
	p.freeF = func(bt BT) {
		empty := new(T)
		if cleanFunc != nil {
			cleanFunc(bt)
		}
		t := any(bt).(*T)
		*t = *empty
	}
	p.p.New = func() any {
		t := new(baseRecycler[BT])
		t.p = p
		t.b = any(new(T)).(BT)
		return t
	}

	pools.Store(reflect.TypeFor[T](), p)
}
