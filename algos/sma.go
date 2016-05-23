package stockfighter

import (
	"container/list"
	"sync"
)

// http://www.day-trading-stocks.org/moving-average-crossover.html

type SimpleMovingAvg struct {
	data   list.List
	period int
}

func InitSimpleMovingAvg(period int) *SimpleMovingAvg {
	return &SimpleMovingAvg{period: period}
}

func (o *SimpleMovingAvg) Push(x int) {
	if o.data.Len() >= o.period {
		o.data.Remove(o.data.Back())
	}
	o.data.PushFront(x)
}

func (o *SimpleMovingAvg) slice() []int {
	var s []int
	for e := o.data.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(int))
	}
	return s
}

func (o *SimpleMovingAvg) Value() int {
	var sum int = 0
	var currPeriod = 0
	for e := o.data.Front(); e != nil && currPeriod <= o.period; e = e.Next() {
		sum += e.Value.(int)
		currPeriod += 1
	}
	return sum / o.period
}

func (o *SimpleMovingAvg) Trend() int {
	data := o.slice()
	var leftSum int = 0
	for _, v := range data[0:len(data)] {
		leftSum += v
	}
	var rightSum int = 0
	for _, v := range data[len(data):] {
		rightSum += v
	}
	return leftSum - rightSum
}

type SmaTriple struct {
	fast  *SimpleMovingAvg
	mid   *SimpleMovingAvg
	slow  *SimpleMovingAvg
	mutex *sync.RWMutex
}

func InitSmaTriple(fast int, mid int, slow int) *SmaTriple {
	return &SmaTriple{
		fast: InitSimpleMovingAvg(fast),
		mid: InitSimpleMovingAvg(mid),
		slow: InitSimpleMovingAvg(slow),
		mutex: &sync.RWMutex{},
	}
}

func (o *SmaTriple) Push(x int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.fast.Push(x)
	o.mid.Push(x)
	o.slow.Push(x)
}

func (o *SmaTriple) Signal() {
	o.mutex.RLock()
	defer o.mutex.Unlock()
	if o.slow.Trend() > 0 {
		if o.fast.Value() > o.mid.Value() {
			return SignalBuy
		} else {
			return SignalSell
		}
	}
}