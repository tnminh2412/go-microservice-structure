package cmd

import (
	"context"
	basegrpc "go-pkg/grpc"
	"go-pkg/helper/id_helper"
	"go-pkg/mongo"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/handler"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize/permission"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/spf13/cobra"

	"fmt"
	"go-pkg/conf"
	"go-pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run BackofficeService service",
	Run: func(cmd *cobra.Command, args []string) {
		StartService()
	},
}

const (
	serviceGroupFlag = "Service"
	serviceNameFlag  = "Name"
	servicePortFlag  = "Port"
)

var (
	serviceName, _ = conf.GetString(serviceGroupFlag, serviceNameFlag, "Backoffice")
	servicePort, _ = conf.GetInt(serviceGroupFlag, servicePortFlag, 9000)
)

func StartService() {
	globalService := global_conf.New(context.Background())
	sharedService := shared_conf.New(context.Background(), globalService)
	confSrv := local_conf.New(sharedService)

	// Methods callable without any token.
	anonymousMethods := []string{
		"/backoffice.Service/Login",
		"/backoffice.Service/RequestOtpReset",
		"/backoffice.Service/ConfirmOtpReset",
	}

	// Methods that require a valid token plus the given permission task.
	requireTaskMethods := map[string]string{
		"/backoffice.Service/Logout":   permission.Task.None,
		"/backoffice.Service/ListRole": permission.Task.ViewRole,

		"/backoffice.Service/GetSystemUser":         permission.Task.GetSystemUser,
		"/backoffice.Service/ListSystemUser":        permission.Task.ViewSystemUser,
		"/backoffice.Service/CreateSystemUser":      permission.Task.CreateSystemUser,
		"/backoffice.Service/UpdateSystemUser":      permission.Task.UpdateSystemUser,
		"/backoffice.Service/SetSystemUserPassword": permission.Task.SetSystemUserPassword,
		"/backoffice.Service/SetSystemUserStatus":   permission.Task.SetSystemUserStatus,
		"/backoffice.Service/ResetSystemUserOtp":    permission.Task.ResetSystemUserOtp,
	}

	// Methods reachable before the user has completed OTP verification.
	skipOtpMethods := []string{}

	// Methods that do not depend on a specific tenant.
	globalMethods := []string{}

	userAuthorizeService := user_authorize.New(anonymousMethods, requireTaskMethods, skipOtpMethods, globalMethods)

	service := handler.New(globalService, sharedService, confSrv, userAuthorizeService)
	server := basegrpc.NewServer(
		func(s *grpc.Server) {
			edge_backoffice.RegisterServiceServer(s, service)
		},
		serviceName,
		servicePort)

	server.SetAuthenticate(userAuthorizeService.Authenticate)

	server.RegisterGW(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return edge_backoffice.RegisterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})
	service.SharedService.AddRoute(server, service.LocalService.Module(), &local_conf.Serialize{})
	service.GlobalService.AddRoute(server)

	server.AddShutdownHook(func() {
		_ = service.SharedService.Shutdown()
		_ = service.GlobalService.Shutdown()
		_ = mongo.Cleanup(context.Background())
		id_helper.ClearUp()
	})

	log.Info(func() string {
		return fmt.Sprintf("Service %s is running at: %d\n", serviceName, servicePort)
	})

	server.Run()
}
