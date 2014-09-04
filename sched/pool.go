package sched

import (
	"container/heap"
)

type item struct {
	id uint16
	pr float64
}

type pool []item

func newPool(cap uint16) pool {
	return make(pool, 0, cap)
}

func (p pool) Len() int {
	return len(p)
}

func (p pool) Less(i, j int) bool {
	return p[i].pr < p[j].pr
}

func (p pool) Swap(i, j int) {
	p[i].id, p[i].pr, p[j].id, p[j].pr = p[j].id, p[j].pr, p[i].id, p[i].pr
}

func (p *pool) Push(i interface{}) {
	*p = append(*p, i.(item))
}

func (p *pool) Pop() interface{} {
	size := len(*p)

	item := (*p)[size-1]
	*p = (*p)[:size-1]

	return item
}

func (p *pool) push(id uint16, pr float64) {
	heap.Push(p, item{id, pr})
}

func (p *pool) pop() uint16 {
	return heap.Pop(p).(item).id
}
