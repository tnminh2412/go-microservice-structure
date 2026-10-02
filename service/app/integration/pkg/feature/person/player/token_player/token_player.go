package token_player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-pkg/helper/error_helper"
	"go-pkg/helper/time_helper"
	"go-pkg/jwt"
	"go-pkg/redis"
	"go-pkg/service"
	"time"

	"github.com/google/uuid"

	gored "github.com/go-redis/redis/v8"
)

const (
	jwtTTL = time.Duration(24*999) * time.Hour
)

type Server struct {
	redis redis.Client
}

type Data struct {
	Token              string
	SessionID          string
	Username           string
	ClientID           string
	IP                 string
	CountryCode        string
	CurrencyCodes      []string
	ActiveCurrencyCode string
	DeviceID           string
	DevicePlatform     string
	LastActive         time.Time
	GameCode           string
	ParentL1           string
	ParentL2           string
	ParentL3           string
	ParentL4           string
	ParentL5           string
	ParentL6           string
	ParentL7           string
	ParentL8           string
	ParentL9           string
	ParentL10          string
	ParentL11          string
	ParentL12          string
	ParentL13          string
	ParentL14          string
	ParentL15          string
	DirectParent       string
	IsTour             bool
	Tour               *Tour
	IsTest             bool
	IsExternal         bool
	ExternalInfo       map[string]string
	Referral           string
	SecretKey          string
	Time               time.Time
}

type Tour struct {
	ID           string `json:"id"`
	MainUsername string `json:"mainUsername"`
	Index        int32  `json:"index"`
}

type dbData struct {
	Token              string `redis:"Token"`
	SessionID          string `redis:"SessionID"`
	Username           string `redis:"Username"`
	ClientID           string `redis:"ClientID"`
	IP                 string `redis:"IP"`
	CountryCode        string `redis:"CountryCode"`
	CurrencyCodes      string `redis:"CurrencyCodes"`
	ActiveCurrencyCode string `redis:"ActiveCurrencyCode"`
	DeviceID           string `redis:"DeviceID"`
	DevicePlatform     string `redis:"DevicePlatform"`
	Time               string `redis:"Time"`
	LastActive         string `redis:"LastActive"`
	GameCode           string `redis:"GameCode"`
	ParentL1           string `redis:"ParentL1"`
	ParentL2           string `redis:"ParentL2"`
	ParentL3           string `redis:"ParentL3"`
	ParentL4           string `redis:"ParentL4"`
	ParentL5           string `redis:"ParentL5"`
	ParentL6           string `redis:"ParentL6"`
	ParentL7           string `redis:"ParentL7"`
	ParentL8           string `redis:"ParentL8"`
	ParentL9           string `redis:"ParentL9"`
	ParentL10          string `redis:"ParentL10"`
	ParentL11          string `redis:"ParentL11"`
	ParentL12          string `redis:"ParentL12"`
	ParentL13          string `redis:"ParentL13"`
	ParentL14          string `redis:"ParentL14"`
	ParentL15          string `redis:"ParentL15"`
	DirectParent       string `redis:"DirectParent"`
	IsTour             bool   `redis:"IsTour"`
	Tour               string `redis:"Tour"`
	IsTest             bool   `redis:"IsTest"`
	IsExternal         bool   `redis:"IsExternal"`
	ExternalInfo       string `redis:"ExternalInfo"`
	Referral           string `redis:"Referral"`

	SecretKey string `redis:"SecretKey"`
}

func toData(d dbData) (*Data, error) {
	loginTime, _ := time.Parse(time_helper.Format_ISO8601, d.Time)
	lastActive, _ := time.Parse(time_helper.Format_ISO8601, d.LastActive)
	externalInfo := make(map[string]string)
	_ = json.Unmarshal([]byte(d.ExternalInfo), &externalInfo)
	var tour Tour
	_ = json.Unmarshal([]byte(d.Tour), &tour)
	currencyCodes := make([]string, 0)
	_ = json.Unmarshal([]byte(d.CurrencyCodes), &currencyCodes)

	data := &Data{
		Token:              d.Token,
		SessionID:          d.SessionID,
		Username:           d.Username,
		ClientID:           d.ClientID,
		IP:                 d.IP,
		CountryCode:        d.CountryCode,
		CurrencyCodes:      currencyCodes,
		ActiveCurrencyCode: d.ActiveCurrencyCode,
		DeviceID:           d.DeviceID,
		DevicePlatform:     d.DevicePlatform,
		Time:               loginTime,
		LastActive:         lastActive,
		GameCode:           d.GameCode,
		ParentL1:           d.ParentL1,
		ParentL2:           d.ParentL2,
		ParentL3:           d.ParentL3,
		ParentL4:           d.ParentL4,
		ParentL5:           d.ParentL5,
		ParentL6:           d.ParentL6,
		ParentL7:           d.ParentL7,
		ParentL8:           d.ParentL8,
		ParentL9:           d.ParentL9,
		ParentL10:          d.ParentL10,
		ParentL11:          d.ParentL11,
		ParentL12:          d.ParentL12,
		ParentL13:          d.ParentL13,
		ParentL14:          d.ParentL14,
		ParentL15:          d.ParentL15,
		DirectParent:       d.DirectParent,
		IsTour:             d.IsTour,
		Tour:               &tour,
		IsTest:             d.IsTest,
		IsExternal:         d.IsExternal,
		ExternalInfo:       externalInfo,
		Referral:           d.Referral,

		SecretKey: d.SecretKey,
	}

	return data, nil
}

type SetData struct {
	Token          string            `json:"token"`
	SessionID      string            `json:"sessionID"`
	Username       string            `json:"username"`
	ClientID       string            `json:"clientID"`
	IP             string            `json:"ip"`
	CountryCode    string            `json:"countryCode"`
	CurrencyCodes  []string          `json:"currencyCodes"`
	DeviceID       string            `json:"deviceID"`
	DevicePlatform string            `json:"devicePlatform"`
	Time           time.Time         `json:"time"`
	LastActive     time.Time         `json:"lastActive"`
	ParentL1       string            `json:"parentL1"`
	ParentL2       string            `json:"parentL2"`
	ParentL3       string            `json:"parentL3"`
	ParentL4       string            `json:"parentL4"`
	ParentL5       string            `json:"parentL5"`
	ParentL6       string            `json:"parentL6"`
	ParentL7       string            `json:"parentL7"`
	ParentL8       string            `json:"parentL8"`
	ParentL9       string            `json:"parentL9"`
	ParentL10      string            `json:"parentL10"`
	ParentL11      string            `json:"parentL11"`
	ParentL12      string            `json:"parentL12"`
	ParentL13      string            `json:"parentL13"`
	ParentL14      string            `json:"parentL14"`
	ParentL15      string            `json:"parentL15"`
	DirectParent   string            `json:"directParent"`
	IsTour         bool              `json:"isTour"`
	Tour           *Tour             `json:"tour"`
	IsTest         bool              `json:"isTest"`
	IsExternal     bool              `json:"isExternal"`
	ExternalInfo   map[string]string `json:"externalInfo"`
	SecretKey      string            `json:"secretKey"`
	Referral       string            `json:"referral"`
}

type JwtData struct {
	Username  string `json:"sid"`
	Timestamp int64  `json:"auth_time"`
}

func toHashValues(data *SetData) (map[string]interface{}, error) {
	externalInfo, _ := json.Marshal(data.ExternalInfo)
	tour, _ := json.Marshal(data.Tour)
	currencyCodes, _ := json.Marshal(data.CurrencyCodes)
	res := map[string]interface{}{
		"Token":          data.Token,
		"SessionID":      data.SessionID,
		"Username":       data.Username,
		"ClientID":       data.ClientID,
		"IP":             data.IP,
		"CountryCode":    data.CountryCode,
		"CurrencyCodes":  string(currencyCodes),
		"DeviceID":       data.DeviceID,
		"DevicePlatform": data.DevicePlatform,
		"Time":           data.Time.Format(time_helper.Format_ISO8601),
		"LastActive":     data.LastActive.Format(time_helper.Format_ISO8601),
		"GameCode":       "",
		"ParentL1":       data.ParentL1,
		"ParentL2":       data.ParentL2,
		"ParentL3":       data.ParentL3,
		"ParentL4":       data.ParentL4,
		"ParentL5":       data.ParentL5,
		"ParentL6":       data.ParentL6,
		"ParentL7":       data.ParentL7,
		"ParentL8":       data.ParentL8,
		"ParentL9":       data.ParentL9,
		"ParentL10":      data.ParentL10,
		"ParentL11":      data.ParentL11,
		"ParentL12":      data.ParentL12,
		"ParentL13":      data.ParentL13,
		"ParentL14":      data.ParentL14,
		"ParentL15":      data.ParentL15,
		"DirectParent":   data.DirectParent,
		"IsTour":         data.IsTour,
		"IsExternal":     data.IsExternal,
		"Tour":           tour,
		"IsTest":         data.IsTest,
		"ExternalInfo":   externalInfo,
		"SecretKey":      data.SecretKey,
		"Referral":       data.Referral,
	}
	return res, nil
}

func New(redisClient redis.Client) *Server {
	return &Server{
		redis: redisClient,
	}
}

func tokenKey(clientID string, token string) string {
	return fmt.Sprintf("_%v:player:token:%v", clientID, token)
}

func usernameKey(clientID string, username string) string {
	return fmt.Sprintf("_%v:player:%v:token", clientID, username)
}

func onlinePlayersKey(clientID string) string {
	return fmt.Sprintf("_%v:players", clientID)
}

func (s *Server) Get(ctx context.Context, clientID string, token string) (*Data, error) {
	cmd := s.redis.HGetAll(ctx, tokenKey(clientID, token))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	if len(cmd.Val()) == 0 {
		return nil, error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
	}
	db := dbData{}
	err := cmd.Scan(&db)
	if err != nil {
		return nil, err
	}

	var jwtData *JwtData
	err = jwt.Validate(token, db.SecretKey, &jwtData)
	if err != nil {
		if errors.Is(err, jwt.TokenExpired) || errors.Is(err, jwt.TokenInvalid) {
			return nil, error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", err, "")
		}
		return nil, error_helper.New(error_helper.Internal, service.Domain, service.Name, "jwt.Validate error", err, "")
	}

	data, err := toData(db)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Server) GetByUsername(ctx context.Context, clientID string, username string) (*Data, error) {
	cmd := s.redis.Get(ctx, usernameKey(clientID, username))
	err := cmd.Err()
	if err != nil {
		if err == gored.Nil {
			return nil, error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
		}
		return nil, err
	}

	if len(cmd.Val()) == 0 {
		return nil, error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
	}
	return s.Get(ctx, clientID, cmd.Val())
}

func (s *Server) Generate(ctx context.Context, data *SetData) (string, error) {
	_tokenSecretKey, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	tokenSecretKey := _tokenSecretKey.String()
	token, err := jwt.NewToken(tokenSecretKey, &JwtData{
		Username:  data.Username,
		Timestamp: time.Now().UnixNano(),
	}, time.Now().Add(jwtTTL))
	if err != nil {
		return "", err
	}
	data.Token = token
	data.SecretKey = tokenSecretKey

	err = s.Set(ctx, data)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Server) Set(ctx context.Context, data *SetData) error {
	fields, err := toHashValues(data)
	if err != nil {
		return err
	}

	onlinePlayer := &gored.Z{
		Score:  float64(time.Now().Unix()),
		Member: data.Username,
	}
	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		p.HMSet(ctx, tokenKey(data.ClientID, data.Token), fields)
		p.Set(ctx, usernameKey(data.ClientID, data.Username), data.Token, 0)
		p.ZAdd(ctx, onlinePlayersKey(data.ClientID), onlinePlayer)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Remove(ctx context.Context, clientID string, token string) error {
	key := tokenKey(clientID, token)
	currentUsername, err := s.redis.HGet(ctx, key, "Username").Result()
	if err != nil {
		if err == gored.Nil {
			return error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
		}
		return err
	}
	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		p.Del(ctx, key)
		p.Del(ctx, usernameKey(clientID, currentUsername))
		p.ZRem(ctx, onlinePlayersKey(clientID), currentUsername)
		return nil
	})
	return err
}

func (s *Server) SetActive(ctx context.Context, clientID string, token string, gameCode string, currencyCode string) error {
	now := time.Now()
	r := luaSetActive.Run(ctx, s.redis, []string{tokenKey(clientID, token)}, now.Format(time_helper.Format_ISO8601), gameCode, currencyCode)
	var data interface{}
	var err error
	if data, err = r.Result(); err == nil {
		if arr, ok := data.([]interface{}); ok {
			success := arr[0].(int64) == 1
			if success == true {
				if username, err := s.redis.HGet(ctx, tokenKey(clientID, token), "Username").Result(); err != nil {
					return err
				} else {
					return s.setOnlinePlayer(ctx, clientID, username, now.Unix())
				}
			} else {
				return error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
			}
		}
	}
	return err
}

func (s *Server) setOnlinePlayer(ctx context.Context, clientID string, username string, score int64) error {
	onlinePlayer := &gored.Z{
		Score:  float64(score),
		Member: username,
	}
	cmd := s.redis.ZAdd(ctx, onlinePlayersKey(clientID), onlinePlayer)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
