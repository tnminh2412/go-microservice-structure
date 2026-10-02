package fast_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/log/base"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	maxCurrentConnections = 2000
	maxConnWaitTimeout    = time.Millisecond * 1000
)

type Option struct {
	FireAndForget bool
}

type IService interface {
	Send(ctx context.Context, url string, method string, bodyReq interface{}, headers map[string]string, target interface{}, timeout time.Duration, opts ...*Option) error
	Get(ctx context.Context, url string, headers map[string]string, target interface{}, timeout time.Duration, opts ...*Option) error
}

type Service struct {
	baseUrl string
	client  *fasthttp.Client
}

func New(baseUrl string) IService {
	client := &fasthttp.Client{
		MaxConnsPerHost:    maxCurrentConnections,
		MaxConnWaitTimeout: maxConnWaitTimeout,
	}

	return &Service{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (s Service) Send(ctx context.Context, url, method string, bodyReq interface{}, headers map[string]string, target interface{}, timeout time.Duration, opts ...*Option) error {
	o := mergeOptions(opts...)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url = fmt.Sprintf("%s%s", s.baseUrl, url)
	req.SetRequestURI(url)

	if bodyReq != nil {
		reqJson, err := json.Marshal(bodyReq)
		if err != nil {
			return Error(service_error.Internal, "json.Marshal error - marshal body request fail", err)
		}
		req.SetBody(reqJson)
	}

	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := s.client.DoTimeout(req, resp, timeout)

	if o.FireAndForget == false {
		if err != nil {
			return err
		}

		bodyBytes := getBodyResponse(resp)
		context_helper.Logger(ctx).Info(func() string {
			return string(bodyBytes)
		})
		err = json.Unmarshal(bodyBytes, target)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "fast_http pkg - unexpected response format"
			}, base.Field("StatusCode", resp.StatusCode()),
				base.Field("RawBody", string(bodyBytes)))
			return Error(service_error.Internal, "json.Unmarshal error - unexpected response format", err)
		}
	}

	return nil
}

func (s Service) Get(ctx context.Context, url string, headers map[string]string, target interface{}, timeout time.Duration, opts ...*Option) error {
	return s.Send(ctx, url, "GET", nil, headers, target, timeout, opts...)
}

func mergeOptions(opts ...*Option) *Option {
	o := &Option{}
	for _, opt := range opts {
		o.FireAndForget = opt.FireAndForget
	}
	return o
}

func getBodyResponse(resp *fasthttp.Response) []byte {
	contentEncoding := resp.Header.Peek(fasthttp.HeaderContentEncoding)
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		httpBody, _ := resp.BodyGunzip()
		body := make([]byte, len(httpBody))
		copy(body, httpBody)
		return body
	}

	httpBody := resp.Body()
	body := make([]byte, len(httpBody))
	copy(body, httpBody)
	return body
}
