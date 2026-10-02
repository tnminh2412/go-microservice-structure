package user

import (
	"context"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/token_user"
)

var (
	tokenServer *token_user.Server
)

type ServiceInfo struct {
	TokenUserService *token_user.Server
}

func Init(info *ServiceInfo) {
	tokenServer = info.TokenUserService
}

func RemoveToken(ctx context.Context, clientID string, token string) error {
	return tokenServer.Remove(ctx, clientID, token)
}

func GetToken(ctx context.Context, clientID string, token string) (*token_user.Data, error) {
	return tokenServer.Get(ctx, clientID, token)
}
