package utils

import (
	"sync"
)

const (
	//WithNanos time format str
	WithNanos = "2006-01-02"
	//TimeFormatStr time format str
	TimeFormatStr = "2006-01-02 15:04:05"
	//DateTimeMarksFormatStr time for str
	DateTimeMarksFormatStr = "20060102150405"
)

//Pool manage task
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

//NewPool init pool
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

//Add add new task
func (p *Pool) Add(delta int) {
	//for i := 0; i < delta; i++ {
	//	p.queue <- 1
	//}
	//for i := 0; i > delta; i-- {
	//	<-p.queue
	//}
	p.queue <- delta
	p.wg.Add(delta)
}

//Done one task
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

//Wait for go task
func (p *Pool) Wait() {
	p.wg.Wait()
}
