package handler

import (
	"context"
	"encoding/json"
	"go-pkg/helper/context_helper"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/module/bot_manager"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"
)

const (
	telegramProvider = "telegram"
)

type Service struct {
	GlobalService global_conf.IService
	SharedService shared_conf.IService
	LocalService  local_conf.IService
	botManager    *bot_manager.BotManager
}

func New() (*Service, error) {
	ctx := context.Background()
	globalService := global_conf.New(ctx)
	sharedService := shared_conf.New(ctx, globalService)
	localService := local_conf.New(sharedService)

	s := &Service{
		GlobalService: globalService,
		SharedService: sharedService,
		LocalService:  localService,
		botManager:    bot_manager.New(),
	}

	err := s.init(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) init(ctx context.Context) error {
	allClients, err := s.GlobalService.GetAllClientConf(ctx)
	if err != nil {
		return err
	}

	for _, client := range allClients {
		for _, socialProvider := range client.SocialProviders {
			if socialProvider == nil || socialProvider.Config == nil || len(socialProvider.Domains) == 0 {
				continue
			}

			botCfg, ok := socialProvider.Config[telegramProvider]
			if ok == false {
				context_helper.Logger(ctx).Warn(func() string {
					return "tele bot config is empty"
				})
				return nil
			}
			webhookHost := socialProvider.Domains[0]
			err = s.botManager.AddBot(botCfg.Token, botCfg.SecretToken, webhookHost)
			if err != nil {
				return err
			}

			for k, cmdCfg := range botCfg.Commands {
				cmdOpt := bot_manager.CommandOtp{}
				cmdOpt, err = toCommandOtp(cmdCfg)
				if err != nil {
					return err
				}
				err = s.botManager.NewCommand(botCfg.Token, k, cmdOpt)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func toCommandOtp(cmd interface{}) (bot_manager.CommandOtp, error) {
	opt := bot_manager.CommandOtp{}
	bytes, err := json.Marshal(cmd)
	if err != nil {
		return opt, err
	}

	err = json.Unmarshal(bytes, &opt)
	return opt, err
}
