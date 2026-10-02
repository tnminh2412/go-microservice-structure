package handler

import (
	"context"
	"go-pkg/service"
	"proto-repo/generated/cmm/service/cmm/edge/backoffice/edge_backoffice"
	"proto-repo/generated/cmm/service/cmm/integration/user/integration_user"
	"strings"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/global"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/access_check"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity/background_log"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/activity/logger_log"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/local_conf"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/backoffice/internal/pkg/request_context"
	"github.com/tnminh2412/go-microservice-structure/service/app/edge/pkg/feature/user_authorize"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	// Returns codes.Unimplemented for business RPCs stripped from this sample.
	edge_backoffice.UnimplementedServiceServer

	integrationUserServiceClient integration_user.ServiceClient
	userAuthorize                user_authorize.IService

	LocalService  local_conf.IService
	GlobalService global_conf.IService
	SharedService shared_conf.IService
	activityLog   activity.Log
	accessCheck   *access_check.Server
}

func New(globalService global_conf.IService,
	sharedService shared_conf.IService,
	confSrv local_conf.Service,
	userAuthorize user_authorize.IService) *Service {

	accessCheck := access_check.New(global.IntegrationUserClient(), confSrv, userAuthorize)

	// Activity logging is fire-and-forget so it never slows down the request path.
	activityLog := background_log.New(logger_log.New())

	return &Service{
		LocalService:  confSrv,
		GlobalService: globalService,
		SharedService: sharedService,

		integrationUserServiceClient: global.IntegrationUserClient(),

		activityLog:   activityLog,
		accessCheck:   accessCheck,
		userAuthorize: userAuthorize,
	}
}

func (s *Service) getSkip(pageIndex int64, limit int64) int64 {
	return pageIndex * limit
}

func (s *Service) requestBy() string {
	return service.Name
}
func (s *Service) requestTime() *timestamp.Timestamp {
	return timestamppb.Now()
}

type activityInfo struct {
	Username   string
	ActivityID string
	Data       map[string]string
}

func (s *Service) logActivity(ctx context.Context, a activityInfo) error {
	now := time.Now()
	actv := activity.Activity{
		Username:    a.Username,
		ActivityID:  a.ActivityID,
		Time:        now,
		Data:        a.Data,
		Country:     request_context.Country(ctx),
		IP:          request_context.IP(ctx),
		CreatedBy:   s.requestBy(),
		CreatedTime: now,
	}
	err := s.activityLog.Append(ctx, request_context.ClientID(ctx), actv)

	if err != nil {
		return err
	}
	return nil
}

func extractUsername(ignoreCase bool, username string) string {
	username = strings.TrimSpace(username)
	if ignoreCase {
		return strings.ToLower(username)
	}

	return username
}
