package shared_conf

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/helper/context_helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service *Service) GetSharedConf(ctx context.Context, clientID string) (*SharedValue, error) {
	configKey := "integration_shared"
	err := service.subscribe(clientID, configKey)
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return fmt.Sprintf("subscribe %v %v error %v", clientID, configKey, err)
		})
	}

	if localConf, ok := service.getConfig(clientID, configKey); ok {
		return localConf.(*SharedValue), nil
	}

	db, err := service.GetDatabase(ctx, clientID)
	if err != nil {
		return nil, err
	}

	row := db.Collection("conf").FindOne(ctx, bson.D{{"_id", configKey}})
	s := &SharedSerialize{}
	result, err := s.Decode(ctx, row)

	if err != nil {
		return nil, err
	}

	service.assignConfig(result, clientID, configKey)
	return result.(*SharedValue), nil
}

// SharedValue holds settings shared by every service of a tenant
// (document "integration_shared" in the tenant's conf collection).
type SharedValue struct {
	ContactInfo *ContactInfo `bson:"contactInfo"`
}

type ContactInfo struct {
	UseDefault map[string]bool `bson:"useDefault"`
	Contacts   []*Contact      `bson:"contacts"`
}

type Contact struct {
	Type        string `bson:"type"`
	Value       string `bson:"value"`
	ID          string `bson:"id"`
	Icon        string `bson:"icon"`
	Url         string `bson:"url"`
	IsShowQR    bool   `bson:"isShowQR"`
	Enabled     bool   `bson:"enabled"`
	Order       int32  `bson:"order"`
	CountryCode string `bson:"countryCode"`
}

type SharedSerialize struct {
}

func (s *SharedSerialize) Decode(_ context.Context, result *mongo.SingleResult) (interface{}, error) {
	var _conf struct {
		Value *SharedValue
	}
	err := result.Decode(&_conf)
	return _conf.Value, err
}

func (s *SharedSerialize) Deserialize(_ context.Context, bytes []byte) (interface{}, error) {
	var _conf struct {
		Value *SharedValue
	}
	err := json.Unmarshal(bytes, &_conf)
	return _conf.Value, err
}
