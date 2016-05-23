package stockfighter

import (
	"testing"
	"reflect"
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
	if !reflect.DeepEqual(sma.slice(), []int{5, 4, 3, 2, 1}) {
		t.Error("bad slice", sma.slice())
	}
}

func TestSimpleMovingAvgPushOverflow(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for i := 1; i <= 5 * testPeriod; i++ {
		sma.Push(i)
	}
	if !reflect.DeepEqual(sma.slice(), []int{25, 24, 23, 22, 21}) {
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
	for i := 1; i <= testPeriod; i++ {
		sma.Push(i)
	}
	if sma.Trend() < 0 {
		t.Error("bad trend", sma.Trend())
	}
}

func TestSimpleMovingAvg_Trend_Negative(t *testing.T) {
	sma := InitSimpleMovingAvg(testPeriod)
	for i := testPeriod; i <=1 ; i-- {
		sma.Push(i)
	}
	if sma.Trend() > 0 {
		t.Error("bad trend", sma.Trend())
	}
}

func TestInitSmaTriple(t *testing.T) {
	sma := InitSmaTriple(5, 8, 13)
	for i := 1; i <= 2 * testPeriod; i++ {
		sma.Push(i)
	}
}