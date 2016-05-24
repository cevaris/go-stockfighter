package algo

import (
	"container/list"
	"sync"
	"github.com/cevaris/stockfighter"
	"fmt"
)

/*
http://www.investopedia.com/articles/active-trading/012815/top-technical-indicators-scalping-trading-strategy.asp
http://www.investopedia.com/articles/active-trading/010116/perfect-moving-averages-day-trading.asp
http://www.investopedia.com/articles/active-trading/012815/top-technical-indicators-scalping-trading-strategy.asp
http://www.day-trading-stocks.org/moving-average-crossover.html
 */

type SimpleMovingAvg struct {
	data   list.List
	period int
}

func InitSimpleMovingAvg(period int) *SimpleMovingAvg {
	return &SimpleMovingAvg{period: period}
}

func (o *SimpleMovingAvg) Push(x int) {
	if o.data.Len() >= o.period {
		o.data.Remove(o.data.Front())
	}
	o.data.PushBack(x)
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

func (o *SimpleMovingAvg) String() string {
	return fmt.Sprintf(
		"SimpleMovingAvg(Period:%d Value:%d Trend:%d Slice:%+v)",
		o.period, o.Value(), o.Trend(), o.slice(),
	)
}

func (o *SimpleMovingAvg) slice() []int {
	var s []int
	for e := o.data.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(int))
	}
	return s
}

func (o *SimpleMovingAvg) Trend() int {
	data := o.slice()
	var leftSum int = 0
	var leftCount int = 0
	for _, v := range data[0:len(data) / 2] {
		leftSum += v
		leftCount += 1
	}
	var rightSum int = 0
	var rightCount int = 0
	for _, v := range data[len(data) / 2:] {
		rightSum += v
		rightCount += 1
	}

	var leftAvg = 0
	if leftCount != 0 {
		leftAvg = leftSum / leftCount
	}

	var rightAvg = 0
	if rightCount != 0 {
		rightAvg = rightSum / rightCount
	}

	x1, y1 := 0, leftAvg
	x2, y2 := len(data), rightAvg

	trend := float64(y2 - y1) / float64(x2 - x1)
	if trend < -0.5 {
		return stockfighter.TrendDown
	}
	if trend > 0.5 {
		return stockfighter.TrendUp
	}
	return stockfighter.TrendUnknown
}

type SmaTriple struct {
	fast  *SimpleMovingAvg
	mid   *SimpleMovingAvg
	slow  *SimpleMovingAvg
	mutex *sync.RWMutex
}

func (o *SmaTriple) String() string {
	return fmt.Sprintf(
		"SmaTriple(Fast: %+v Mid: %+v, Slow: %+v Signal:%d)",
		o.fast, o.mid, o.slow, o.Signal(),
	)
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

func (o *SmaTriple) Signal() int {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	if o.slow.Trend() == stockfighter.TrendUp {
		if o.fast.Value() > o.mid.Value() {
			return stockfighter.SignalBuy
		} else {
			return stockfighter.SignalSell
		}
	}
	if o.slow.Trend() == stockfighter.TrendDown {
		return stockfighter.SignalSell
	}
	return stockfighter.SignalUnknown
}