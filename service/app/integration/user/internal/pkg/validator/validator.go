package validator

import (
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	normalizer "github.com/dimuska139/go-email-normalizer"
)

type Validator struct {
	GlobalConf      global_conf.IService
	confSrv         local_conf.IService
	emailNormalizer *normalizer.Normalizer
}

func New(globalConf global_conf.IService, confSrv local_conf.IService) Validator {
	return Validator{
		GlobalConf:      globalConf,
		confSrv:         confSrv,
		emailNormalizer: normalizer.NewNormalizer(),
	}
}
