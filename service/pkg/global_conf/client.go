package global_conf

import (
	"context"
	"errors"
	"fmt"
	"go-pkg/helper/context_helper"
	"go-pkg/log/base"
	m "go-pkg/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/time_function"
)

const clientCollection = "client"

type Row struct {
	ClientID    string    `bson:"_id"`
	Conf        *Conf     `bson:"conf"`
	CreatedTime time.Time `bson:"createdTime"`
	CreatedBy   string    `bson:"createdBy"`
	UpdatedTime time.Time `bson:"updatedTime"`
	UpdatedBy   string    `bson:"updatedBy"`
}

type Conf struct {
	Providers       []*Provider       `bson:"providers"`
	Apps            []*App            `bson:"apps"`
	Timezone        string            `bson:"timezone"`
	Db              Db                `bson:"db"`
	SSODb           Db                `bson:"ssoDb"`
	SocialProviders []*SocialProvider `bson:"socialProviders"`
}

type Db struct {
	DbName              string `bson:"dbName"`
	Connection          string `bson:"connection"`
	EncryptedConnection string `bson:"encryptedConnection"`
}

type Provider struct {
	Provider string   `bson:"provider"`
	AppIDs   []string `bson:"appIDs"`
}

type App struct {
	LoginUsername string `bson:"loginUsername"`
	LoginPassword string `bson:"loginPassword"`

	SecretKey    string `bson:"secretKey"`
	Username     string `bson:"username"`
	AppID        string `bson:"appID"`
	CurrencyCode string `bson:"currencyCode"`
	Method       string `bson:"method"`

	ExternalUrl       string   `bson:"externalUrl"`
	ExternalTimeout   int32    `bson:"externalTimeout"`
	WhiteListIPs      []string `bson:"whiteListIPs"`
	Status            bool     `bson:"status"`
	IsMaintaince      bool     `bson:"isMaintaince"`
	IsSupportFreePlay bool     `bson:"isSupportFreePlay"`
}

type SocialProvider_Config struct {
	Token       string                 `bson:"token"`
	SecretToken string                 `bson:"secretToken"`
	Commands    map[string]interface{} `bson:"commands"`
}

type SocialProvider struct {
	Domains []string                          `bson:"domains"`
	Config  map[string]*SocialProvider_Config `bson:"config"`
}

func (service *Service) GetClientConf(ctx context.Context, clientID string) (*Conf, error) {
	if err := service.subscribeClient(); err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "subscribeClient error"
		}, base.Error(err))
	}

	if localConf, ok := service.getConfig(clientID); ok {
		return localConf.Clone(), nil
	}

	//don't lock here so that warm up requests not waiting for each others
	configs, err := service.getConfigs(ctx)
	if err != nil {
		return nil, err
	}

	service.assignConfig(configs)

	localConf, ok := service.getConfig(clientID)
	if !ok {
		return nil, errors.New("no configuration")
	}
	return localConf.Clone(), nil
}

func (service *Service) GetAllClientConf(ctx context.Context) (map[string]*Conf, error) {
	ok := service.hasConfig()
	if !ok {
		configs, err := service.getConfigs(ctx)
		if err != nil {
			return nil, err
		}
		service.assignConfig(configs)
	}
	return service.cloneConf(), nil
}

func (service *Service) subscribeClient() error {
	err := service.subscribe(clientCollection, func(ctx context.Context, data *m.WatchEvent) {
		service.removeConfig()
	})
	return err
}

func (service *Service) getConfigs(ctx context.Context) (map[string]*Conf, error) {
	db, err := m.Client()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "get global db error"
		}, base.Error(err))
		return nil, err
	}

	cursor, err := db.Database("app_global").Collection("client").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	rows := make(map[string]*Conf)
	for cursor.Next(ctx) {
		var row Row
		if err := cursor.Decode(&row); err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return fmt.Sprintf("fetch configuration error %v", err)
			})
			return nil, err
		} else {
			rows[row.ClientID] = row.Conf
		}
	}
	return rows, nil
}

func (service *Service) getConfig(clientID string) (*Conf, bool) {
	service.configLock.RLock()
	defer service.configLock.RUnlock()

	if c, ok := service.configs[clientID]; ok {
		return c, true
	}

	return nil, false
}

func (service *Service) hasConfig() bool {
	service.configLock.RLock()
	defer service.configLock.RUnlock()

	return len(service.configs) > 0
}

func (service *Service) assignConfig(values map[string]*Conf) {
	service.configLock.Lock()
	defer service.configLock.Unlock()

	service.configs = values
}

func (service *Service) removeConfig() {
	service.configLock.Lock()
	defer service.configLock.Unlock()

	service.configs = make(map[string]*Conf)

}

func (service *Service) cloneConf() map[string]*Conf {
	service.configLock.RLock()
	defer service.configLock.RUnlock()

	configs := make(map[string]*Conf)

	for k, v := range service.configs {
		configs[k] = v.Clone()
	}

	return configs
}

func (c *Conf) Clone() *Conf {
	return &Conf{
		Providers: func() []*Provider {
			providers := make([]*Provider, len(c.Providers))

			for i, p := range c.Providers {
				providers[i] = &Provider{
					Provider: p.Provider,
					AppIDs:   p.AppIDs,
				}
			}

			return providers

		}(),
		Apps: func() []*App {
			apps := make([]*App, len(c.Apps))

			for i, p := range c.Apps {
				apps[i] = &App{
					LoginUsername:     p.LoginUsername,
					LoginPassword:     p.LoginPassword,
					SecretKey:         p.SecretKey,
					Username:          p.Username,
					AppID:             p.AppID,
					CurrencyCode:      p.CurrencyCode,
					Method:            p.Method,
					ExternalUrl:       p.ExternalUrl,
					ExternalTimeout:   p.ExternalTimeout,
					WhiteListIPs:      p.WhiteListIPs,
					Status:            p.Status,
					IsMaintaince:      p.IsMaintaince,
					IsSupportFreePlay: p.IsSupportFreePlay,
				}
			}

			return apps

		}(),
		SocialProviders: func() []*SocialProvider {
			sps := make([]*SocialProvider, len(c.SocialProviders))
			for i, p := range c.SocialProviders {
				sps[i] = &SocialProvider{
					Domains: p.Domains,
					Config: func() map[string]*SocialProvider_Config {
						cfg := make(map[string]*SocialProvider_Config, len(p.Config))
						for k, v := range p.Config {
							cfg[k] = &SocialProvider_Config{
								Token:       v.Token,
								SecretToken: v.SecretToken,
								Commands:    v.Commands,
							}
						}
						return cfg
					}(),
				}
			}
			return sps
		}(),
		Timezone: c.Timezone,
		Db: Db{
			DbName:              c.Db.DbName,
			Connection:          c.Db.Connection,
			EncryptedConnection: c.Db.EncryptedConnection,
		},
		SSODb: Db{
			DbName:              c.SSODb.DbName,
			Connection:          c.SSODb.Connection,
			EncryptedConnection: c.SSODb.EncryptedConnection,
		},
	}
}

func (c *Conf) Client() (*mongo.Client, error) {
	return m.Client(c.Db.DbName)
}

func (c *Conf) SSOClient() (*mongo.Client, error) {
	return m.Client(c.SSODb.DbName)
}

func (c *Conf) Location() (*time.Location, error) {
	return time_function.LoadLocation(c.Timezone)
}
