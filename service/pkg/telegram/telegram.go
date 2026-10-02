package telegram

import (
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Manager struct {
	clients    sync.Map
	httpClient *http.Client
}

func NewManager() *Manager {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Manager{
		httpClient: &http.Client{
			Transport: t,
			Timeout:   30 * time.Second,
		},
	}
}

func NewManagerWithTimeout(timeout time.Duration) *Manager {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Manager{
		httpClient: &http.Client{
			Transport: t,
			Timeout:   timeout,
		},
	}
}

func (m *Manager) Get(token string) (*tgbotapi.BotAPI, error) {
	v, ok := m.clients.Load(token)
	if ok {
		return v.(*tgbotapi.BotAPI), nil
	}

	telegramApi, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, m.httpClient)
	if err != nil {
		return nil, err
	}

	v, _ = m.clients.LoadOrStore(token, telegramApi)
	return v.(*tgbotapi.BotAPI), nil
}
