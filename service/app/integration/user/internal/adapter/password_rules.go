package adapter

import (
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
)

func ToUserPasswordRules(list []local_conf.PasswordRule) []*integration_user.Locale {
	arr := make([]*integration_user.Locale, 0)

	for _, rule := range list {
		for _, r := range rule.Locales {
			arr = append(arr, toLocaleRule(r.Language, r.Description))
		}
	}

	return arr
}
