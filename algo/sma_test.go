package algo

import (
	"testing"
	"reflect"
	"github.com/cevaris/stockfighter"
)

var testPeriod int = 5

func TestSimpleMovingAvg(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	if sma.Value() != 0 {
		t.Error("bad value", sma.Value())
	}
}

func TestSimpleMovingAvg_Push(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for i := 1; i <= testPeriod; i++ {
		sma.Push(i)
	}
	if !reflect.DeepEqual(sma.slice(), []int{1, 2, 3, 4, 5}) {
		t.Error("bad slice", sma.slice())
	}
}

func TestSimpleMovingAvgPushOverflow(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for i := 1; i <= 5 * testPeriod; i++ {
		sma.Push(i)
	}
	if !reflect.DeepEqual(sma.slice(), []int{21, 22, 23, 24, 25}) {
		t.Error("bad slice", sma.slice())
	}
	if sma.Value() != 23 {
		t.Error("bad value", sma.Value())
	}
}

func TestSimpleMovingAvg_Value(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for i := 1; i <= testPeriod; i++ {
		sma.Push(i)
	}
	if sma.Value() != 3 {
		t.Error("bad value", sma.Value())
	}
}

func TestSimpleMovingAvg_Trend_Positive(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for _, v := range []int{10, 4, 8, 12, 15} {
		sma.Push(v)
	}
	if sma.Trend() != stockfighter.TrendUp {
		t.Error("bad trend", sma.Trend())
	}
}

func TestSimpleMovingAvg_Trend_Negative(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for _, v := range []int{22, 20, 15, 17, 10} {
		sma.Push(v)
	}
	if sma.Trend() != stockfighter.TrendDown {
		t.Error("bad trend", sma.Trend())
	}
}

func TestInitSmaTriple(t *testing.T) {
	sma := InitSmaTriple(5, 8, 13)
	for i := 1; i <= 2 * testPeriod; i++ {
		sma.Push(i)
	}
}

func TestSmaTriple_Signal_Buy(t *testing.T) {
	sma := InitSmaTriple(1, 2, 3)
	for _, v := range []int{1, 2, 4, 8} {
		sma.Push(v)
	}
	if sma.slow.Trend() != stockfighter.TrendUp {
		t.Error("bad signal", sma.slow.Trend())
	}
	if sma.Signal() != stockfighter.SignalBuy {
		t.Error("bad signal", sma.Signal())
	}
}

func TestSmaTriple_Signal_Sell(t *testing.T) {
	sma := InitSmaTriple(1, 2, 3)
	for _, v := range []int{5, 7, 5, 5, 1} {
		sma.Push(v)
	}
	if sma.Signal() != stockfighter.SignalSell {
		t.Error("bad signal", sma.Signal())
	}
}

func TestSmaTriple_Signal_Unknown(t *testing.T) {
	sma := InitSmaTriple(1, 2, 3)
	for _, v := range []int{1, 3, 1, 3} {
		sma.Push(v)
	}
	if sma.Signal() != stockfighter.SignalUnknown {
		t.Error("bad signal", sma.Signal())
	}
}