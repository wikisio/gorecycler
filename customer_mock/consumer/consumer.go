package main

import (
	"github.com/wikisio/gorecycler/customer_mock/producer"
	"github.com/wikisio/gorecycler/recycle"
	"log"
	"os"
	"runtime/trace"
)

type node struct {
	i     int
	b     [32]byte
	left  producer.Node
	right producer.Node
}

func newNode(n *node, layer int) {
	if layer <= 0 {
		return
	}

	l := recycle.Get[producer.Node, node]()
	r := recycle.Get[producer.Node, node]()
	l.Assign(func(bt producer.Node) {
		bt.(*node).i = layer + 1
		n.left = bt
		newNode(bt.(*node), layer-1)
	})
	r.Assign(func(bt producer.Node) {
		bt.(*node).i = layer + 1
		n.right = bt
		newNode(bt.(*node), layer-1)
	})
}

func NewNode() recycle.Recycler[producer.Node] {
	root := recycle.Get[producer.Node, node]()
	root.Assign(func(bt producer.Node) {
		i := 3
		newNode(bt.(*node), i)
	})

	return root
}
func (i *node) Sum() int {
	return sum(i)
}

func sum(j *node) int {
	var i, l, r int
	if j == nil {
		return 0
	}

	i = j.i
	if j.left != nil {
		l = sum(j.left.(*node))
	}

	if j.right != nil {
		r = sum(j.right.(*node))
	}
	return i + l + r
}

func cleanF(bt producer.Node) {
	if bt == nil {
		return
	}

	l := bt.(*node).left
	r := bt.(*node).right

	if l != nil {
		lc := recycle.FindPool[producer.Node](l)
		lc.HandleAndRecycle(func(bt producer.Node) error {
			return nil
		})
	}

	if r != nil {
		rc := recycle.FindPool[producer.Node](r)
		rc.HandleAndRecycle(func(bt producer.Node) error {
			return nil
		})
	}
}

func main() {
	recycle.RegisterPoolWithCleaner[producer.Node, node](cleanF)

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()

	for i := 0; i < 1000*10; i++ {
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
}
