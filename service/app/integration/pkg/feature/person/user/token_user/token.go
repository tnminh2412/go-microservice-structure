package token_user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-pkg/helper/error_helper"
	"go-pkg/jwt"
	"go-pkg/service"
	"time"

	"github.com/google/uuid"

	gored "github.com/go-redis/redis/v8"
)

const (
	jwtTTL = time.Hour * 24 * 999
)

type Server struct {
	redis gored.Cmdable
}

type Data struct {
	Token       string
	SubUsername string
	Username    string
	ClientID    string
	Roles       []string
	Tasks       []string
	AccessLines []string
	ParentL1    string
	ParentL2    string
	ParentL3    string
	ParentL4    string
	ParentL5    string
	ParentL6    string
	ParentL7    string
	ParentL8    string
	ParentL9    string
	ParentL10   string
	ParentL11   string
	ParentL12   string
	ParentL13   string
	ParentL14   string
	ParentL15   string
	Level       int32
	IsSystem    bool
	IsMain      bool
	RequireOtp  bool

	SecretKey string
}

type dbData struct {
	Token       string `redis:"Token"`
	SubUsername string `redis:"SubUsername"`
	Username    string `redis:"Username"`
	ClientID    string `redis:"ClientID"`
	Roles       string `redis:"Roles"`
	Tasks       string `redis:"Tasks"`
	AccessLines string `redis:"AccessLines"`
	ParentL1    string `redis:"ParentL1"`
	ParentL2    string `redis:"ParentL2"`
	ParentL3    string `redis:"ParentL3"`
	ParentL4    string `redis:"ParentL4"`
	ParentL5    string `redis:"ParentL5"`
	ParentL6    string `redis:"ParentL6"`
	ParentL7    string `redis:"ParentL7"`
	ParentL8    string `redis:"ParentL8"`
	ParentL9    string `redis:"ParentL9"`
	ParentL10   string `redis:"ParentL10"`
	ParentL11   string `redis:"ParentL11"`
	ParentL12   string `redis:"ParentL12"`
	ParentL13   string `redis:"ParentL13"`
	ParentL14   string `redis:"ParentL14"`
	ParentL15   string `redis:"ParentL15"`
	Level       int32  `redis:"Level"`
	IsSystem    bool   `redis:"IsSystem"`
	IsMain      bool   `redis:"IsMain"`
	RequireOtp  bool   `redis:"RequireOtp"`

	SecretKey string `redis:"SecretKey"`
}

func toData(d dbData) (*Data, error) {
	roles := make([]string, 0)
	err := json.Unmarshal([]byte(d.Roles), &roles)
	if err != nil {
		return nil, err
	}

	tasks := make([]string, 0)
	err = json.Unmarshal([]byte(d.Tasks), &tasks)
	if err != nil {
		return nil, err
	}

	accessLines := make([]string, 0)
	if len(d.AccessLines) > 0 {
		err = json.Unmarshal([]byte(d.AccessLines), &accessLines)
		if err != nil {
			return nil, err
		}
	}

	data := &Data{
		Token:       d.Token,
		SubUsername: d.SubUsername,
		Username:    d.Username,
		ClientID:    d.ClientID,
		Roles:       roles,
		Tasks:       tasks,
		AccessLines: accessLines,
		ParentL1:    d.ParentL1,
		ParentL2:    d.ParentL2,
		ParentL3:    d.ParentL3,
		ParentL4:    d.ParentL4,
		ParentL5:    d.ParentL5,
		ParentL6:    d.ParentL6,
		ParentL7:    d.ParentL7,
		ParentL8:    d.ParentL8,
		ParentL9:    d.ParentL9,
		ParentL10:   d.ParentL10,
		ParentL11:   d.ParentL11,
		ParentL12:   d.ParentL12,
		ParentL13:   d.ParentL13,
		ParentL14:   d.ParentL14,
		ParentL15:   d.ParentL15,
		Level:       d.Level,
		IsSystem:    d.IsSystem,
		IsMain:      d.IsMain,
		RequireOtp:  d.RequireOtp,

		SecretKey: d.SecretKey,
	}
	return data, nil
}

type SetData struct {
	Token       string
	SubUsername string
	Username    string
	AccessLines []string
	ClientID    string
	Roles       []string
	Tasks       []string
	ParentL1    string
	ParentL2    string
	ParentL3    string
	ParentL4    string
	ParentL5    string
	ParentL6    string
	ParentL7    string
	ParentL8    string
	ParentL9    string
	ParentL10   string
	ParentL11   string
	ParentL12   string
	ParentL13   string
	ParentL14   string
	ParentL15   string
	Level       int32
	IsSystem    bool
	IsMain      bool
	RequireOtp  bool

	SecretKey string
}

type JwtData struct {
	Username  string `json:"sid"`
	Timestamp int64  `json:"auth_time"`
}

func toHashValues(data *SetData) (map[string]interface{}, error) {
	roles, err := json.Marshal(data.Roles)
	if err != nil {
		return nil, err
	}
	tasks, err := json.Marshal(data.Tasks)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"Token":       data.Token,
		"SubUsername": data.SubUsername,
		"Username":    data.Username,
		"ClientID":    data.ClientID,
		"Roles":       roles,
		"Tasks":       tasks,
		"ParentL1":    data.ParentL1,
		"ParentL2":    data.ParentL2,
		"ParentL3":    data.ParentL3,
		"ParentL4":    data.ParentL4,
		"ParentL5":    data.ParentL5,
		"ParentL6":    data.ParentL6,
		"ParentL7":    data.ParentL7,
		"ParentL8":    data.ParentL8,
		"ParentL9":    data.ParentL9,
		"ParentL10":   data.ParentL10,
		"ParentL11":   data.ParentL11,
		"ParentL12":   data.ParentL12,
		"ParentL13":   data.ParentL13,
		"ParentL14":   data.ParentL14,
		"ParentL15":   data.ParentL15,
		"Level":       data.Level,
		"IsSystem":    data.IsSystem,
		"IsMain":      data.IsMain,
		"RequireOtp":  data.RequireOtp,

		"SecretKey": data.SecretKey,
	}
	if data.IsSystem {
		accessLines, err := json.Marshal(data.AccessLines)
		if err != nil {
			return nil, err
		}
		res["AccessLines"] = accessLines
	}
	return res, nil
}

func New(redisClient gored.Cmdable) *Server {
	return &Server{
		redis: redisClient,
	}
}

func tokenKey(clientID string, token string) string {
	return fmt.Sprintf("_%v:user:token:%v", clientID, token)
}

func usernameKey(clientID string, username string) string {
	return fmt.Sprintf("_%v:user:%v:tokens", clientID, username)
}

func subUsernameKey(clientID string, username string, subUsername string) string {
	return fmt.Sprintf("_%v:user:%v:%v:tokens", clientID, username, subUsername)
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

func (s *Server) Generate(ctx context.Context, data *SetData, ttl time.Duration) (string, error) {
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

	err = s.Set(ctx, data, ttl)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Server) Set(ctx context.Context, data *SetData, ttl time.Duration) error {
	now := time.Now()

	fields, err := toHashValues(data)
	if err != nil {
		return err
	}

	key := tokenKey(data.ClientID, data.Token)
	userKey := usernameKey(data.ClientID, data.Username)
	subUserKey := subUsernameKey(data.ClientID, data.Username, data.SubUsername)
	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		p.HMSet(ctx, key, fields)
		p.ExpireAt(ctx, key, now.Add(ttl))
		p.SAdd(ctx, userKey, data.Token)
		p.ExpireAt(ctx, userKey, now.Add(ttl))
		p.SAdd(ctx, subUserKey, data.Token)
		p.ExpireAt(ctx, subUserKey, now.Add(ttl))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Remove(ctx context.Context, clientID string, token string) error {
	key := tokenKey(clientID, token)
	vals, err := s.redis.HMGet(ctx, key, "Username", "SubUsername").Result()
	if err != nil {
		if err == gored.Nil {
			return error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
		}
		return err
	}

	username, ok := vals[0].(string)
	if !ok {
		return error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
	}
	subUsername, ok := vals[1].(string)
	if !ok {
		return error_helper.New(error_helper.NotFound, service.Domain, service.Name, "token not found", nil, "")
	}

	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		p.Del(ctx, key)
		p.SRem(ctx, usernameKey(clientID, username), token)
		p.SRem(ctx, subUsernameKey(clientID, username, subUsername), token)
		return nil
	})
	return err
}

func (s *Server) RemoveByUsername(ctx context.Context, clientID string, username string) error {
	key := usernameKey(clientID, username)
	tokens, err := s.redis.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		for _, t := range tokens {
			p.Del(ctx, tokenKey(clientID, t))
		}
		p.Del(ctx, key)

		return nil
	})
	return err
}

func (s *Server) RemoveBySubUsername(ctx context.Context, clientID string, username string, subUsername string) error {
	key := subUsernameKey(clientID, username, subUsername)
	tokens, err := s.redis.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = s.redis.TxPipelined(ctx, func(p gored.Pipeliner) error {
		for _, t := range tokens {
			p.Del(ctx, tokenKey(clientID, t))
			p.SRem(ctx, usernameKey(clientID, username), t)
		}
		p.Del(ctx, key)

		return nil
	})
	return err
}
