package cmd

import (
	"context"
	basegrpc "go-pkg/grpc"
	"go-pkg/mongo"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/handler"
	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/local_conf"

	"github.com/spf13/cobra"

	"fmt"
	"go-pkg/conf"
	"go-pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run UserService service",
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
	serviceName, _ = conf.GetString(serviceGroupFlag, serviceNameFlag, "User")
	servicePort, _ = conf.GetInt(serviceGroupFlag, servicePortFlag, 9000)
)

func StartService() {
	service := handler.New()
	server := basegrpc.NewServer(
		func(s *grpc.Server) {
			core_user.RegisterUserServiceServer(s, service)
		},
		serviceName,
		servicePort)

	server.RegisterGW(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return core_user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})
	service.SharedService.AddRoute(server, service.LocalService.Module(), &local_conf.Serialize{})
	service.GlobalService.AddRoute(server)

	server.AddShutdownHook(func() {
		_ = service.SharedService.Shutdown()
		_ = service.GlobalService.Shutdown()
		_ = mongo.Cleanup(context.Background())
	})

	log.Info(func() string {
		return fmt.Sprintf("Service %s is running at: %d\n", serviceName, servicePort)
	})

	server.Run()
}
