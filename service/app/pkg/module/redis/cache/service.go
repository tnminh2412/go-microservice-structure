package cache

import (
	"context"
	"encoding/json"
	"fmt"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/helper/json_helper"
	"go-pkg/redis"
	"strconv"
	"time"

	gored "github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

const (
	cacheServer = "Cache"
)

var (
	cacheClient gored.Cmdable
)

func init() {
	cacheClient = redis.GetUniversal(cacheServer)
}

func Get(ctx context.Context, clientID string, key string, target interface{}) (bool, error) {
	cacheKey := getCacheKey(clientID, key)
	val, err := cacheClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == gored.Nil {
			return false, nil
		}
		return false, Error(service_error.Internal, "cacheClient.Get error", err)
	}

	if len(val) == 0 {
		return false, nil
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		return false, Error(service_error.Internal, "json.Unmarshal error", err)
	}

	return true, nil
}

func Del(ctx context.Context, clientID string, key string) error {
	cacheKey := getCacheKey(clientID, key)
	_, err := cacheClient.Del(ctx, cacheKey).Result()
	if err != nil {
		return err
	}

	return nil
}

func Set(ctx context.Context, clientID string, key string, data interface{}, expiration time.Duration) error {
	cacheKey := getCacheKey(clientID, key)
	_, err := cacheClient.Set(ctx, cacheKey, json_helper.SafeJsonToString(data), expiration).Result()
	if err != nil {
		return Error(service_error.Internal, "cacheClient.Set error", err)
	}

	return nil
}

func IncrByDecimal(ctx context.Context, clientID string, key string, incValue decimal.Decimal) (decimal.Decimal, error) {
	cacheKey := getCacheKey(clientID, key)
	v, err := strconv.ParseFloat(incValue.StringFixed(16), 64)
	if err != nil {
		return decimal.Zero, Error(service_error.Internal, "strconv.ParseFloat error", err)
	}

	result, err := cacheClient.IncrByFloat(ctx, cacheKey, v).Result()
	if err != nil {
		return decimal.Zero, Error(service_error.Internal, "s.redis.IncrByFloat error", err)
	}

	return decimal.NewFromFloat(result), nil
}

func Expire(ctx context.Context, clientID string, key string, expiration time.Duration) error {
	cacheKey := getCacheKey(clientID, key)
	_, err := cacheClient.Expire(ctx, cacheKey, expiration).Result()
	if err != nil {
		return Error(service_error.Internal, "cacheClient.Expire error", nil)
	}

	return err
}

func HGet(ctx context.Context, clientID string, key, field string, target interface{}) (bool, error) {
	cacheKey := getCacheKey(clientID, key)
	val, err := cacheClient.HGet(ctx, cacheKey, field).Result()
	if err != nil {
		if err == gored.Nil {
			return false, nil
		}
		return false, Error(service_error.Internal, "cacheClient.HGet error", err)
	}

	if len(val) == 0 {
		return false, nil
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		return false, Error(service_error.Internal, "json.Unmarshal error", err)
	}

	return true, nil
}

func HGetAll(ctx context.Context, clientID string, key string) (map[string]string, error) {
	cacheKey := getCacheKey(clientID, key)
	val, err := cacheClient.HGetAll(ctx, cacheKey).Result()
	if err != nil {
		if err == gored.Nil {
			return nil, nil
		}
		return nil, Error(service_error.Internal, "cacheClient.HGet error", err)
	}

	return val, nil
}

func HSet(ctx context.Context, clientID string, key string, data interface{}, expiration time.Duration) error {
	cacheKey := getCacheKey(clientID, key)
	_, err := cacheClient.TxPipelined(ctx, func(p gored.Pipeliner) error {
		p.HSet(ctx, cacheKey, data)
		if expiration != -1 {
			p.Expire(ctx, cacheKey, expiration)
		}
		return nil
	})
	if err != nil {
		return Error(service_error.Internal, "cacheClient.Set error", err)
	}

	return nil
}

func getCacheKey(clientID, key string) string {
	return fmt.Sprintf("_%s:cache:%s", clientID, key)
}
