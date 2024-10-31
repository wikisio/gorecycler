package consumer

import (
	"fmt"
	"github.com/wikisio/gorecycler/customer_mock/producer"
	"github.com/wikisio/gorecycler/recycle"
	"testing"
)

type Int struct {
	i int
	b [32]byte
}

func (i *Int) Display() string {
	return fmt.Sprintf("%d", i.i)
}

func TestRegisterPool(t *testing.T) {
	recycle.RegisterPool[producer.Object, Int]()

	i := recycle.Get[producer.Object, Int]()
	i.Assign(func(t producer.Object) {
		t.(*Int).i = 13
	})

	i.HandleAndRecycle(func(o producer.Object) error {
		fmt.Println(o.Display())
		if o.Display() != "13" {
			t.Fail()
		}
		return nil
	})

	i = recycle.Get[producer.Object, Int]()
	i.HandleAndRecycle(func(o producer.Object) error {
		fmt.Println(o.Display())
		if o.Display() != "0" {
			t.Fail()
		}
		return nil
	})

}
