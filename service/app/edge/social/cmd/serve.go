package cmd

import (
	"context"
	basegrpc "go-pkg/grpc"
	"go-pkg/mongo"
	"proto-repo/generated/cmm/service/cmm/edge/social/edge_social"

	"github.com/spf13/cobra"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/handler"

	"fmt"
	"go-pkg/conf"
	"go-pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run EdgeSocial service",
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
	serviceName, _ = conf.GetString(serviceGroupFlag, serviceNameFlag, "GenerateID")
	servicePort, _ = conf.GetInt(serviceGroupFlag, servicePortFlag, 9000)
)

func StartService() {
	service, err := handler.New()
	if err != nil {
		panic(err)
	}
	server := basegrpc.NewServer(
		func(s *grpc.Server) {
			edge_social.RegisterServiceServer(s, service)
		},
		serviceName,
		servicePort)

	server.RegisterGW(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return edge_social.RegisterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})

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
