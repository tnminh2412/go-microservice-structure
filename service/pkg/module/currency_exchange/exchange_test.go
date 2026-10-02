package currency_exchange

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExchange(t *testing.T) {
	currencyPkg := New()
	ctx := context.TODO()
	clientID := "demo"
	Convey("convert thb to myr", t, func() {
		res, err := currencyPkg.Exchange(ctx, &ExchangeRequest{
			ClientID:         clientID,
			FromCurrencyCode: "THB",
			ToCurrencyCode:   "MYR",
			Amount:           decimal.New(1, 0),
		})
		So(err, ShouldBeNil)
		So(res.Amount.GreaterThan(decimal.Zero), ShouldBeTrue)
		t.Log(res.Amount)
	})
	Convey("convert myr to thb", t, func() {
		res, err := currencyPkg.Exchange(ctx, &ExchangeRequest{
			ClientID:         clientID,
			FromCurrencyCode: "MYR",
			ToCurrencyCode:   "THB",
			Amount:           decimal.New(1, 0),
		})
		So(err, ShouldBeNil)
		So(res.Amount.GreaterThan(decimal.Zero), ShouldBeTrue)
		t.Log(res.Amount)
	})
}
