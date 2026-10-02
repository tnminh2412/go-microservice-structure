package adapter

import (
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
)

func ToUserUsernameRules(list []local_conf.UsernameRule) []*integration_user.Locale {
	arr := make([]*integration_user.Locale, 0)
	for _, rule := range list {
		for _, r := range rule.Locales {
			arr = append(arr, toLocaleRule(r.Language, r.Description))
		}
	}

	return arr
}

func toLocaleRule(language string, description string) *integration_user.Locale {
	return &integration_user.Locale{
		Language:    language,
		Description: description,
	}
}
