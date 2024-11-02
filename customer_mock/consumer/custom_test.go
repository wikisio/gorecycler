package main

import (
	"fmt"
	"github.com/wikisio/gorecycler/customer_mock/producer"
	"github.com/wikisio/gorecycler/recycle"
	"testing"
)

func TestRegisterPool(t *testing.T) {
	recycle.RegisterPool[producer.Node, node]()

	i := NewNode()

	i.Assign(func(t producer.Node) {
		pi := recycle.FindPool[producer.Node](t)
		fmt.Println(i, pi)
	})

	i.HandleAndRecycle(func(o producer.Node) error {
		fmt.Println(o.Sum())
		if o.Sum() == 0 {
			t.Fail()
		}
		return nil
	})

	i = recycle.Get[producer.Node, node]()
	i.HandleAndRecycle(func(o producer.Node) error {
		fmt.Println(o.Sum())
		if o.Sum() == 0 {
			t.Fail()
		}
		return nil
	})
}
func TestPool(t *testing.T) {
	recycle.RegisterPoolWithCleaner[producer.Node, node](cleanF)

	x := NewNode()
	x.HandleAndRecycle(func(bt producer.Node) error {
		bt.Sum()
		return nil
	})

}

func BenchmarkPool(b *testing.B) {
	recycle.RegisterPoolWithCleaner[producer.Node, node](cleanF)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := NewNode()
			x.HandleAndRecycle(func(bt producer.Node) error {
				bt.Sum()
				return nil
			})
		}
	})
}
