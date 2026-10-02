package currency_exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log/base"
	"time"

	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type ExchangeRequest struct {
	ClientID         string
	FromCurrencyCode string
	ToCurrencyCode   string
	Amount           decimal.Decimal
}

type ExchangeResult struct {
	Amount decimal.Decimal
	Rate   decimal.Decimal
}

func (s *Service) Exchange(ctx context.Context, req *ExchangeRequest) (*ExchangeResult, error) {
	if req.FromCurrencyCode == req.ToCurrencyCode {
		return &ExchangeResult{
			Amount: req.Amount,
			Rate:   decimal.New(1, 0),
		}, nil
	}

	rates, err := s.getRates(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}
	fromRate, ok := rates[req.FromCurrencyCode]
	if !ok {
		return nil, Error(InvalidCurrency, fmt.Sprintf("currency %s is invalid", req.FromCurrencyCode), nil)
	}
	toRate, ok := rates[req.ToCurrencyCode]
	if !ok {
		return nil, Error(InvalidCurrency, fmt.Sprintf("currency %s is invalid", req.ToCurrencyCode), nil)
	}

	rate := toRate.Div(fromRate)
	amount := req.Amount.Mul(rate)
	return &ExchangeResult{
		Amount: amount,
		Rate:   rate,
	}, nil
}

type exchangeApiResponse struct {
	Status  int32              `json:"status"`
	Message string             `json:"message"`
	Rates   []*exchangeApiRate `json:"rates"`
}

type exchangeApiRate struct {
	CurrencyCode string          `json:"currencycode"`
	Rate         decimal.Decimal `json:"rate"`
}

const (
	exchange_api_status_success = 200
)

func (s *Service) getRates(ctx context.Context, clientID string) (map[string]decimal.Decimal, error) {
	cacheRates, ok := s.cache.GetRates(clientID)
	if ok {
		return cacheRates, nil
	}

	rv, err, _ := s.group.Do(clientID, func() (interface{}, error) {

		url := fmt.Sprintf("%s/exchange-rates?from=%s", baseUrl, baseCurrencyCode)
		deadline := time.Now().Add(timeout)

		context_helper.Logger(ctx).Info(func() string {
			return fmt.Sprintf("start call exchange rate")
		}, base.Field("clientID", clientID), base.Field("baseCurrencyCode", baseCurrencyCode))

		bodyBytes, statusCode, err := doDeadline(s.fastHttp, url, fasthttp.MethodGet, deadline)
		if err != nil {
			return nil, Error(service_error.Internal, "doDeadline error", err)
		}
		if statusCode != fasthttp.StatusOK {
			return nil, Error(service_error.ErrorExternal, fmt.Sprintf("get exchange rates failed, statusCode: %v", statusCode), nil)
		}

		var exchangeApiRes exchangeApiResponse
		err = json.Unmarshal(bodyBytes, &exchangeApiRes)
		if err != nil {
			return nil, Error(service_error.Internal, "json.Unmarshal error", err)
		}

		if exchangeApiRes.Status != exchange_api_status_success {
			return nil, Error(service_error.ErrorExternal, fmt.Sprintf("get exchange rates failed, exchangeApiRes.Status: %v", exchangeApiRes.Status), nil)
		}

		rates := make(map[string]decimal.Decimal)
		for _, rate := range exchangeApiRes.Rates {
			rates[rate.CurrencyCode] = rate.Rate
		}

		s.cache.SetRates(clientID, rates, time.Minute*5)
		return rates, nil
	})
	if err != nil {
		return nil, err
	}

	return rv.(map[string]decimal.Decimal), nil
}

func doDeadline(fastHttp *fasthttp.Client, url string, method string, deadline time.Time) (bodyBytes []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(method)

	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	err = fastHttp.DoDeadline(req, resp, deadline)
	bodyBytes = resp.Body()
	var result = make([]byte, len(bodyBytes))
	copy(result, bodyBytes)
	return result, resp.StatusCode(), err
}
