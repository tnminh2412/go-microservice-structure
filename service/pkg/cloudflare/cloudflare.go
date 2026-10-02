package cloudflare

import (
	"net/http"
	"sync"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Manager struct {
	clients sync.Map
}

func NewManager() *Manager {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	c := &http.Client{
		Transport: t,
		Timeout:   30 * time.Second,
	}
	cloudflare.HTTPClient(c)
	return &Manager{}
}

func (m *Manager) Get(token string) (*cloudflare.API, error) {
	v, ok := m.clients.Load(token)
	if ok {
		return v.(*cloudflare.API), nil
	}

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}

	v, _ = m.clients.LoadOrStore(token, api)
	return v.(*cloudflare.API), nil
}
