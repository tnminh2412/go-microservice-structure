package exchange_token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/helper/json_helper"
	"go-pkg/jwt"
	"go-pkg/redis"
	"time"

	gored "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	cacheServer                          = "Cache"
	exchangeTokenExpireTimeInMinus int64 = 5
)

var (
	cacheClient gored.Cmdable
)

func init() {
	cacheClient = redis.GetUniversal(cacheServer)
}

type jwtData struct {
}

type exchangeDbData struct {
	SecretKey string
	Data      string
}

func NewToken(ctx context.Context, clientID string, keyPrefix string, data interface{}) (string, error) {
	_secretKey, err := uuid.NewRandom()
	if err != nil {
		return "", Error(service_error.Internal, "uuid.NewRandom error", err)
	}
	secretKey := _secretKey.String()

	expiration := time.Duration(exchangeTokenExpireTimeInMinus) * time.Minute
	expiryAt := time.Now().Add(expiration)
	token, err := jwt.NewToken(secretKey, &jwtData{}, expiryAt)
	if err != nil {
		return "", Error(service_error.Internal, "jwt.NewToken error", err)
	}

	key := exchangeTokenKey(clientID, keyPrefix, token)
	setData := &exchangeDbData{
		SecretKey: secretKey,
		Data:      json_helper.SafeJsonToString(data),
	}
	_, err = cacheClient.Set(ctx, key, json_helper.SafeJsonToString(setData), expiration).Result()
	if err != nil {
		return "", Error(service_error.Internal, "cacheClient.Set error", err)
	}

	return token, nil
}

func Verify(ctx context.Context, clientID string, keyPrefix string, token string, target interface{}) error {
	key := exchangeTokenKey(clientID, keyPrefix, token)
	val, err := cacheClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, gored.Nil) {
			return Error(service_error.Unauthenticated, "Token not found", err)
		}
		return Error(service_error.Internal, "cacheClient.Get error", err)
	}

	if len(val) == 0 {
		return Error(service_error.Unauthenticated, "Token not found", err)
	}

	db := exchangeDbData{}
	err = json.Unmarshal([]byte(val), &db)
	if err != nil {
		return Error(service_error.Internal, "json.Unmarshal error", err)
	}
	var jData *jwtData
	err = jwt.Validate(token, db.SecretKey, &jData)
	if err != nil {
		if errors.Is(err, jwt.TokenExpired) || errors.Is(err, jwt.TokenInvalid) {
			return Error(service_error.Unauthenticated, "Token not found", err)
		}
		return Error(service_error.Internal, "jwt.Validate error", err)
	}

	err = json.Unmarshal([]byte(db.Data), target)
	if err != nil {
		return Error(service_error.Internal, "json.Unmarshal error", err)
	}

	_, err = cacheClient.Del(ctx, key).Result()
	if err != nil {
		return Error(service_error.Internal, "cacheClient.Del error", err)
	}

	return nil
}

func exchangeTokenKey(clientID string, keyPrefix string, token string) string {
	return fmt.Sprintf("_%v:cache:exchange_token:%v:%v", clientID, keyPrefix, token)
}
