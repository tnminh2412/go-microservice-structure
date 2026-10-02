package global

import (
	"fmt"
	"go-pkg/conf"
	"go-pkg/grpc"
	"go-pkg/log"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/config/service_config"
)

var (
	integrationUser integration_user.ServiceClient
)

func IntegrationUserClient() integration_user.ServiceClient {
	return integrationUser
}

func init() {
	integrationUserHost, _ := conf.GetString(service_config.Group, service_config.IntegrationUserHostFlag, "")
	integrationUserConn, err := grpc.NewClient(integrationUserHost)
	if err != nil {
		log.Error(func() string { return fmt.Sprintf("err: %s", err.Error()) })
		panic(err)
	}
	integrationUser = integration_user.NewServiceClient(integrationUserConn)
}
