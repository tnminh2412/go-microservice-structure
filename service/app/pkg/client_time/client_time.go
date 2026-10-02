package client_time

import (
	"context"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
)

func ToClientTime(ctx context.Context, globalService global_conf.IService, clientID string, t time.Time) (time.Time, error) {
	cf, err := globalService.GetClientConf(ctx, clientID)
	if err != nil {
		return t, err
	}
	loc, err := cf.Location()
	if err != nil {
		return t, err
	}

	return t.In(loc), nil
}

func GetClientLocation(ctx context.Context, globalService global_conf.IService, clientID string) (*time.Location, error) {
	cf, err := globalService.GetClientConf(ctx, clientID)
	if err != nil {
		return nil, err
	}
	loc, err := cf.Location()
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func ToLocationTime(t time.Time, location *time.Location) time.Time {
	return t.In(location)
}
