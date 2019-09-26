package utils

import (
	"github.com/emirpasic/gods/utils"
	"sync"

	"github.com/emirpasic/gods/sets/treeset"
)

type Comparator func(a, b interface{}) int

type OrderSet struct {
	s *treeset.Set // 这个有序集合，不是协程安全的
	m sync.Mutex
	l int //大小限制, only top n
}

func NewWith(comparator utils.Comparator, size int) *OrderSet {
	return &OrderSet{s: treeset.NewWith(comparator), l: size}
}

func (ors *OrderSet) Add(items ...interface{}) {
	ors.m.Lock()
	ors.s.Add(items...)
	ors.m.Unlock()
}

func (ors *OrderSet) AddTopN(items ...interface{}) {
	ors.m.Lock()
	ors.s.Add(items...)
	if ors.s.Size() > ors.l {
		ors.s.Remove(ors.s.Values()[ors.s.Size()-1])
	}
	ors.m.Unlock()
}

func (ors *OrderSet) Remove(items ...interface{}) {
	ors.m.Lock()
	ors.s.Remove(items...)
	ors.m.Unlock()
}

func (ors *OrderSet) Contains(items ...interface{}) bool {
	ors.m.Lock()
	val := ors.s.Contains(items...)
	ors.m.Unlock()
	return val

}

func (ors *OrderSet) Empty() bool {
	ors.m.Lock()
	val := ors.s.Empty()
	ors.m.Unlock()
	return val
}

func (ors *OrderSet) Size() int {
	ors.m.Lock()
	val := ors.s.Size()
	ors.m.Unlock()
	return val
}

func (ors *OrderSet) Clear() {
	ors.m.Lock()
	ors.s.Clear()
	ors.m.Unlock()
}

func (ors *OrderSet) Values() []interface{} {
	ors.m.Lock()
	val := ors.s.Values()
	ors.m.Unlock()
	return val
}

func (ors *OrderSet) String() string {
	ors.m.Lock()
	val := ors.s.String()
	ors.m.Unlock()
	return val
}
