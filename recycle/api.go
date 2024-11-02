package recycle

import (
	"reflect"
)

type Recycler[BT any] interface {
	HandleAndRecycle(handleBtFunc func(bt BT) error) error
	Assign(h func(bt BT))
}

func Get[BT any, T any]() Recycler[BT] {
	p, _ := pools.Load(reflect.TypeFor[T]())
	return p.(*pool[BT]).get()

}

func RegisterPool[BT any, T any]() {
	RegisterPoolWithCleaner[BT, T](nil)
}

// RegisterPoolWithCleaner registers a pool To allocates instance of T
// T is the target struct To be used. *T should implement the BT interface.
// if T needs To do close or clean operations, do it in cleanFunc, otherwise nil
func RegisterPoolWithCleaner[BT any, T any](cleanFunc func(bt BT)) {
	_, ok := pools.Load(reflect.TypeFor[T]())
	if ok {
		return
	}

	p := &pool[BT]{}
	p.empty = any(new(T)).(BT)
	p.cleanFunc = func(bt BT) {
		if cleanFunc != nil {
			cleanFunc(bt)
		}
		*(any(bt).(*T)) = *any(p.empty).(*T)

	}
	p.p.New = func() any {
		t := new(recycle[BT])
		t.p = p
		b := new(T)
		t.b = any(b).(BT)
		//bts.Store(t.b, any(t))
		return t
	}

	pools.Store(reflect.TypeFor[T](), p)
}

func FindPool[BT any](bt BT) Recycler[BT] {
	return nil
	b, ok := bts.Load(bt)
	if ok {
		return b.(Recycler[BT])
	}
	panic("FindPool failed")
	return nil
}
