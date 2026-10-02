package cmd

import (
	"context"
	basegrpc "go-pkg/grpc"
	"go-pkg/grpc/header"
	"go-pkg/helper/context_helper"
	"go-pkg/helper/header_helper"
	"go-pkg/helper/id_helper"
	"go-pkg/mongo"
	"net/http"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/handler"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/user/internal/pkg/local_conf"

	"github.com/spf13/cobra"

	"fmt"
	"go-pkg/conf"
	"go-pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	serviceName, _ = conf.GetString(serviceGroupFlag, serviceNameFlag, "IntegrationUser")
	servicePort, _ = conf.GetInt(serviceGroupFlag, servicePortFlag, 9000)
)

func StartService() {
	service := handler.New()
	server := basegrpc.NewServer(
		func(s *grpc.Server) {
			integration_user.RegisterServiceServer(s, service)
		},
		serviceName,
		servicePort)

	server.RegisterGW(func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		mux.Handle("POST",
			runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"insert-image"}, "", runtime.AssumeColonVerbOpt(true))),
			func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
				ctx, cancel := context.WithCancel(req.Context())
				defer cancel()
				_, outboundMarshaller := runtime.MarshalerForRequest(mux, req)
				rctx, err := runtime.AnnotateContext(ctx, mux, req)
				if err != nil {
					runtime.HTTPError(ctx, mux, outboundMarshaller, w, req, err)
					return
				}
				var md metadata.MD
				var ok bool
				if md, ok = metadata.FromOutgoingContext(rctx); !ok {
					md = metadata.MD{}
				}
				rctx = context_helper.NewRequestContext(rctx, header_helper.NewWithMD(header.GetRaw(md)))

				newReq := req.WithContext(rctx)
				service.CreateImage(w, newReq)
			})
		return integration_user.RegisterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
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
