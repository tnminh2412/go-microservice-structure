package global

import (
	"go-pkg/redis"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/internal/config"
)

var (
	tokenClient         redis.Client
	schedulerLockClient redis.Client
)

func SchedulerLockClient() redis.Client {
	return schedulerLockClient
}

func TokenClient() redis.Client {
	return tokenClient
}

func init() {
	tokenClient = redis.GetUniversal(config.TokenServer)
	schedulerLockClient = redis.GetUniversal(config.SchedulerLockServer)

}
