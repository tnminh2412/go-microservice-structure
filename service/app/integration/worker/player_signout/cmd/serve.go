package cmd

import (
	"context"
	basegrpc "go-pkg/grpc"
	"go-pkg/helper/id_helper"
	"go-pkg/mongo"
	"proto-repo/generated/cmm/service/cmm/integration/worker/player_signout/worker_player_signout"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/handler"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/pkg/local_conf"

	"github.com/spf13/cobra"

	"fmt"
	"go-pkg/conf"
	"go-pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run PlayerSignoutService service",
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
	serviceName, _ = conf.GetString(serviceGroupFlag, serviceNameFlag, "PlayerSignoutService")
	servicePort, _ = conf.GetInt(serviceGroupFlag, servicePortFlag, 9000)
)

func StartService() {

	worker := handler.New()
	server := basegrpc.NewServer(
		func(s *grpc.Server) {
			worker_player_signout.RegisterServiceServer(s, worker)
		},
		serviceName,
		servicePort)

	server.RegisterGW(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return worker_player_signout.RegisterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})

	worker.SharedService.AddRoute(server, worker.LocalService.Module(), &local_conf.Serialize{})
	worker.GlobalService.AddRoute(server)

	worker.Start()
	server.AddShutdownHook(func() {
		worker.Stop()
		_ = worker.SharedService.Shutdown()
		_ = worker.GlobalService.Shutdown()
		id_helper.ClearUp()
		_ = mongo.Cleanup(context.Background())
	})

	log.Info(func() string {
		return fmt.Sprintf("Service %s is running at: %d\n", serviceName, servicePort)
	})

	server.Run()
}
