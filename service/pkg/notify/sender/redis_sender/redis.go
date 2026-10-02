package redis_sender

import "fmt"

type pubNotifyMessage struct {
	ClientID  string
	Topics    map[string]struct{}
	DirectIDs map[string]struct{}
	Data      string
	Provider  string
}

type pubNotifyManagement struct {
	ClientID     string
	Event        string
	NotifyID     string
	Username     string
	UserType     string
	CurrencyCode string
	Data         map[string]string
}

type redisTopicType struct{}

var redisTopic redisTopicType

func (t redisTopicType) NotifyDirectMessage(clientID string) string {
	return fmt.Sprintf("_%s.notify.direct", clientID)
}

func (t redisTopicType) NotifyMessage(clientID string) string {
	return fmt.Sprintf("_%s.notify", clientID)
}

func (t redisTopicType) NotifyManagement(clientID string) string {
	return fmt.Sprintf("_%s.notify.management", clientID)
}

const (
	notifyManagementMaxLen = 10000
	notifyMessageMaxLen    = 10000
)

var event = struct {
	Register    string
	DeRegister  string
	Subscribe   string
	UnSubscribe string
}{
	Register:    "Register",
	DeRegister:  "DeRegister",
	Subscribe:   "Subscribe",
	UnSubscribe: "UnSubscribe",
}
