package adapter

import (
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/system_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/validator"
)

func ToUsernameRules(cfgUsernameRules []local_conf.UsernameRule) []validator.UsernameRule {
	usernameRules := make([]validator.UsernameRule, len(cfgUsernameRules))
	for i, u := range cfgUsernameRules {
		usernameRules[i] = validator.UsernameRule{
			Regex: u.Regex,
			Locales: func() []validator.Locale {
				var reply []validator.Locale
				for _, locale := range u.Locales {
					reply = append(reply, validator.Locale{
						Language:    locale.Language,
						Description: locale.Description,
					})
				}
				return reply
			}(),
		}
	}
	return usernameRules
}

func ToPasswordRules(cfgPasswordRules []local_conf.PasswordRule) []validator.PasswordRule {
	passwordRules := make([]validator.PasswordRule, len(cfgPasswordRules))
	for i, u := range cfgPasswordRules {
		passwordRules[i] = validator.PasswordRule{
			Regex: u.Regex,
			Locales: func() []validator.Locale {
				var reply []validator.Locale
				for _, locale := range u.Locales {
					reply = append(reply, validator.Locale{
						Language:    locale.Language,
						Description: locale.Description,
					})
				}
				return reply
			}(),
		}
	}
	return passwordRules
}

func ToDefaultChatConfigSystemUser(chatConfig *local_conf.Chat) *system_user.Chat {
	if chatConfig == nil {
		return nil
	}
	chat := &system_user.Chat{
		FreshWork: &system_user.FreshWork{
			PlayerApp: func() *system_user.App {
				if chatConfig.FreshWork == nil {
					return nil
				}
				return toConfigSystemApp(chatConfig.FreshWork.PlayerApp)
			}(),
			AgentID: func() string {
				if chatConfig.FreshWork == nil {
					return ""
				}
				return chatConfig.FreshWork.AgentID
			}(),
			Enable: func() bool {
				if chatConfig.FreshWork == nil {
					return false
				}
				return chatConfig.FreshWork.Enable
			}(),
		},
	}
	return chat
}

func toConfigSystemApp(app *local_conf.App) *system_user.App {
	if app == nil {
		return nil
	}
	return &system_user.App{
		Ticket:      toConfigSystemFreshWorkApp(app.Ticket),
		ChatSupport: toConfigSystemFreshWorkApp(app.ChatSupport),
	}
}

func toConfigSystemFreshWorkApp(freshWorkApp *local_conf.FreshWorkApp) *system_user.FreshWorkApp {
	if freshWorkApp == nil {
		return nil
	}
	return &system_user.FreshWorkApp{
		ClientID:     freshWorkApp.ClientID,
		LoginPageUri: freshWorkApp.LoginPageUri,
		Name:         freshWorkApp.Name,
		Description:  freshWorkApp.Description,
	}
}
