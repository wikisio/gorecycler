package main

import (
	"fmt"
	"github.com/wikisio/gorecycler/customer_mock/producer"
	"github.com/wikisio/gorecycler/recycle"
	"testing"
)

func TestPool(t *testing.T) {
	recycle.RegisterPoolWithCleaner[producer.Node, node](cleanF)

	x := NewNode()
	x.HandleAndRecycle(func(bt producer.Node) error {
		bt.Sum()
		return nil
	})

	fmt.Println("From: ", recycle.From.Load(), " to: ", recycle.To.Load(), " Clean: ", Clean.Load())
}

func BenchmarkPool(b *testing.B) {
	recycle.RegisterPoolWithCleaner[producer.Node, node](cleanF)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for m := 0; m < 100; m++ {
				for n := 0; n < 2; n++ {
					x := NewNode()
					x.HandleAndRecycle(func(bt producer.Node) error {
						bt.Sum()
						return nil
					})
				}
			}
		}
	})
	fmt.Println("From: ", recycle.From.Load(), " to: ", recycle.To.Load(), " Clean: ", Clean.Load())
}
