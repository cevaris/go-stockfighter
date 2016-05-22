package stockfighter

import (
	"sync"
	"fmt"
)

const (
	SignalBuy = 1
	SignalUnknown = 0
	SignalSell = -1
)

type Macd struct {
	FastChan       chan float32
	fastCurr       float32
	fastFactor     float32
	SlowChan       chan float32
	slowCurr       float32
	slowFactor     float32
	mutex          *sync.Mutex
	signal         int
	signalBoundary float32
}

func NewMacd(slow float32, fast float32) *Macd {
	macd := &Macd{
		FastChan: make(chan float32),
		fastCurr: 0.0,
		fastFactor: 1.0 - (1.0 / fast),
		SlowChan: make(chan float32),
		slowCurr: 0.0,
		slowFactor: 1.0 - (1.0 / slow),
		mutex: &sync.Mutex{},
		signal: SignalUnknown,
		signalBoundary: 0.5,
	}

	return macd
}

func (macd *Macd) Signal() int {
	return macd.signal
}

func (macd *Macd) calcSignal() int {
	fmt.Println("Slow", macd.slowCurr, "Fast", macd.fastCurr)

	if (macd.slowCurr - macd.fastCurr) > macd.signalBoundary {
		return SignalSell
	} else if (macd.fastCurr - macd.slowCurr) > macd.signalBoundary {
		return SignalBuy
	} else {
		return SignalUnknown
	}
}

func (macd *Macd)Put(v int) {
	macd.mutex.Lock()
	macd.fastCurr = movingAvg(macd.fastFactor, macd.fastCurr, float32(v))
	macd.slowCurr = movingAvg(macd.slowFactor, macd.slowCurr, float32(v))
	macd.signal = macd.calcSignal()
	macd.mutex.Unlock()
}

func movingAvg(smooth float32, prevVal float32, newVal float32) float32 {
	return (newVal * smooth) + (prevVal * ( 1.0 - smooth))
}
