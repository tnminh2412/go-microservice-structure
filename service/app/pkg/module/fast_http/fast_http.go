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

type IService interface {
	Send(ctx context.Context, url, method string, bodyReq interface{}, headers map[string]string, target interface{}) error
}

type service struct {
	baseUrl string
	client  *fasthttp.Client
	timeout time.Duration
}

func New(baseUrl string, timeout time.Duration) IService {
	client := &fasthttp.Client{
		MaxConnsPerHost:    maxCurrentConnections,
		MaxConnWaitTimeout: maxConnWaitTimeout,
	}

	return &service{
		baseUrl: baseUrl,
		client:  client,
		timeout: timeout,
	}
}

func (s *service) NewRequest(url string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	return req
}

func (s *service) Send(ctx context.Context, url, method string, bodyReq interface{}, headers map[string]string, target interface{}) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url = fmt.Sprintf("%s%s", s.baseUrl, url)
	req.SetRequestURI(url)

	if bodyReq != nil {
		reqJson, err := json.Marshal(bodyReq)
		if err != nil {
			return Error(service_error.Internal, "json.Marshal error", err)
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

	err := s.client.DoTimeout(req, resp, s.timeout)
	if err != nil {
		return err
	}

	bodyBytes := getBodyResponse(resp)

	err = json.Unmarshal(bodyBytes, target)
	if err != nil {
		context_helper.Logger(ctx).Debug(func() string {
			return "Body"
		}, base.Field("Body", string(bodyBytes)))
		return Error(service_error.Internal, "json.Unmarshal error", err)
	}

	return nil
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
