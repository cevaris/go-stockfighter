package stockfighter

import (
	"testing"
)

var testVenue = "TestVenue"

func TestInitSession(t *testing.T) {
	session := InitSession(&config{}, testVenue)
	if session.Position != 0 {
		t.Error("invalid position", session.Position)
	}
	if session.Cash != 0 {
		t.Error("invalid cash", session.Cash)
	}
	if session.NAV != 0 {
		t.Error("invalid nav", session.NAV)
	}
}

func TestSessionUpdate(t *testing.T) {
	session := InitSession(&config{}, testVenue)
	session.LatestQuote = &StockQuote{Last: 100}
	session.Update(&StockOrderAccountStatus{
		Ok: true,
		Venue: testVenue,
		Orders: []*StockOrder{
			&StockOrder{
				Ok: true,
				Direction: DirectionBuy,
				Fills: []*Fill{
					&Fill{Qty:10, Price:100},
					&Fill{Qty:5, Price:200},
				},
			},
			&StockOrder{
				Ok: true,
				Direction: DirectionSell,
				Fills: []*Fill{
					&Fill{Qty:10, Price:200},
				},
			},
		},
	})
	if session.Position != (10 + 5 - 10) {
		t.Error("invalid position", session.Position)
	}
	if session.Cash != 2000 - 2000 {
		t.Error("invalid cash", session.Cash)
	}
	if session.NAV != 5 * 100 {
		t.Error("invalid nav", session.NAV)
	}
}

func TestSessionUpdateOkOnlyStockOrders(t *testing.T) {
	session := InitSession(&config{}, testVenue)

	session.Update(&StockOrderAccountStatus{
		Ok: true,
		Venue: testVenue,
		Orders: []*StockOrder{
			&StockOrder{
				Ok: true,
				Direction: DirectionBuy,
				Fills: []*Fill{
					&Fill{Qty:10, Price:100},
					&Fill{Qty:5, Price:200},
				},
			},
			&StockOrder{
				Ok: false,
				Direction: DirectionBuy,
				Fills: []*Fill{
					&Fill{Qty:100, Price:100},
				},
			},
		},
	})
	if session.Position != 15 {
		t.Error("invalid position", session.Position)
	}
	if session.Cash != -2000 {
		t.Error("invalid cash", session.Cash)
	}
	if session.NAV != 0 {
		t.Error("invalid nav", session.NAV)
	}
}

func TestSessionUpdateOkOnlyStatus(t *testing.T) {
	session := InitSession(&config{}, testVenue)

	session.Update(&StockOrderAccountStatus{
		Ok: false,
		Venue: testVenue,
		Orders: []*StockOrder{
			&StockOrder{
				Ok: true,
				Direction: DirectionBuy,
				Fills: []*Fill{
					&Fill{Qty:10, Price:100},
					&Fill{Qty:5, Price:200},
				},
			},
		},
	})
	if session.Position != 0 {
		t.Error("invalid position", session.Position)
	}
	if session.Cash != 0 {
		t.Error("invalid cash", session.Cash)
	}
	if session.NAV != 0 {
		t.Error("invalid nav", session.NAV)
	}
}