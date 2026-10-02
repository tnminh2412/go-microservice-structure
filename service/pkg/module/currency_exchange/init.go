package currency_exchange

import (
	"go-pkg/conf"
	"time"
)

const (
	groupFlag   = "Exchange"
	baseUrlFlag = "BaseUrl"
	timeoutFlag = "Timeout"

	baseCurrencyCode = "USD"
)

var (
	baseUrl string
	timeout time.Duration
)

func init() {
	baseUrl, _ = conf.GetString(groupFlag, baseUrlFlag, "")
	timeout, _ = conf.GetDuration(groupFlag, timeoutFlag, time.Second*3)
}
