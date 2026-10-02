package token_player

import (
	"context"
	"errors"
	"fmt"
	"time"

	gored "github.com/go-redis/redis/v8"
)

type ClearInactivePlayerParams struct {
	Limit         int64
	MinActiveTime time.Time
}

type ClearInactivePlayerResult struct {
	Count     int64
	Usernames []string
}

func (s *Server) ClearInactivePlayer(ctx context.Context, clientID string, params ClearInactivePlayerParams) (*ClearInactivePlayerResult, error) {
	opt := &gored.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprint(params.MinActiveTime.Unix()),
		Count: params.Limit,
	}
	onlKey := onlinePlayersKey(clientID)
	usernames, err := s.redis.ZRangeByScore(ctx, onlKey, opt).Result()
	if err != nil {
		return nil, err
	}

	if len(usernames) == 0 {
		return &ClearInactivePlayerResult{}, nil
	}

	cmds, err := s.redis.Pipelined(ctx, func(p gored.Pipeliner) error {
		for _, u := range usernames {
			p.Get(ctx, usernameKey(clientID, u))
		}
		return nil
	})
	if err != nil {
		if err != gored.Nil {
			return nil, err
		}
	}

	tokens := make([]string, 0)
	validUsernames := make([]string, 0)
	for i, c := range cmds {
		err := c.Err()
		if err != nil {
			if err == gored.Nil {
				continue
			}
			return nil, err
		}
		cmd, ok := c.(*gored.StringCmd)
		if !ok {
			return nil, errors.New("pipeline command cast error")
		}
		tokens = append(tokens, cmd.Val())
		validUsernames = append(validUsernames, usernames[i])
	}

	if len(tokens) == 0 {
		return &ClearInactivePlayerResult{}, nil
	}

	_, err = s.redis.Pipelined(ctx, func(p gored.Pipeliner) error {
		for _, u := range validUsernames {
			p.Del(ctx, usernameKey(clientID, u))
			p.ZRem(ctx, onlKey, u)
		}
		for _, t := range tokens {
			p.Del(ctx, tokenKey(clientID, t))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ClearInactivePlayerResult{
		Count:     int64(len(usernames)),
		Usernames: usernames,
	}, nil
}
