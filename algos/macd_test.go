package stockfighter

import (
	"testing"
)

func TestNewMacDWithDefaults(t *testing.T) {
	macd := NewMacd(5.0, 10.0)
	actual := macd.Signal()
	if actual != SignalUnknown {
		t.Error("bad default signal", actual)
	}
}

func TestMacdSameRate(t *testing.T) {
	macd := NewMacd(5.0, 10.0)

	for _, v := range []int{10, 10, 10} {
		macd.Put(v)
	}

	actual := macd.Signal()
	if actual != SignalUnknown {
		t.Error("bad default signal", actual)
	}
}

func TestMacdSameRateWithinBoundary(t *testing.T) {
	macd := NewMacd(5.0, 10.0)

	for _, v := range []int{11, 12, 11} {
		macd.Put(v)
	}

	actual := macd.Signal()
	if actual != SignalUnknown {
		t.Error("bad default signal", actual)
	}
}

func TestMacdBuy(t *testing.T) {
	macd := NewMacd(5.0, 10.0)

	for _, v := range []int{34, 36, 37, 41, 40, 39} {
		macd.Put(v)
	}

	actual := macd.Signal()
	if actual != SignalBuy {
		t.Error("bad default signal", actual)
	}
}

func TestMacdSell(t *testing.T) {
	macd := NewMacd(5.0, 10.0)

	for _, v := range []int{30, 15, 10, 5} {
		macd.Put(v)
	}

	actual := macd.Signal()
	if actual != SignalSell {
		t.Error("bad default signal", actual)
	}
}