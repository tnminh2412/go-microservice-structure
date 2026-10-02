package token_player

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	gored "github.com/go-redis/redis/v8"
)

type SumOnlinePlayerResult struct {
	GameSummaries    []GameSummary
	TotalPlayerCount int64
}

type GameSummary struct {
	GameCode    string
	PlayerCount int64
}

type SumOnlineParams struct {
	CurrencyCodes     []string
	GameCodes         []string
	Usernames         []string
	MinActiveGameTime time.Time
	Parent            *Parent
}

func (s *Server) SumOnlinePlayer(ctx context.Context, clientID string, params SumOnlineParams) (*SumOnlinePlayerResult, error) {

	byCurrency := len(params.CurrencyCodes) > 0
	byGame := len(params.GameCodes) > 0
	byUsername := len(params.Usernames) > 0
	checkParent, byParent := getCheckParentFn(params.Parent)
	currencies := func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, v := range params.CurrencyCodes {
			m[v] = struct{}{}
		}
		return m
	}()
	games := func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, v := range params.GameCodes {
			m[v] = struct{}{}
		}
		return m
	}()
	users := func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, v := range params.Usernames {
			m[v] = struct{}{}
		}
		return m
	}()

	opt := &gored.ZRangeBy{
		Min: fmt.Sprint(params.MinActiveGameTime.Unix()),
		Max: "+inf",
	}
	usernames, err := s.redis.ZRevRangeByScore(ctx, onlinePlayersKey(clientID), opt).Result()
	if err != nil {
		return nil, err
	}
	if len(usernames) == 0 {
		return &SumOnlinePlayerResult{}, nil
	}

	cmds, err := s.redis.Pipelined(ctx, func(p gored.Pipeliner) error {
		for _, u := range usernames {
			if byUsername {
				_, ok := users[u]
				if !ok {
					continue
				}
			}
			p.Get(ctx, usernameKey(clientID, u))
		}
		return nil
	})
	if err != nil && err != gored.Nil {
		return nil, err
	}

	tokens := make([]string, 0)
	for _, c := range cmds {
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
	}

	if len(tokens) == 0 {
		return &SumOnlinePlayerResult{}, nil
	}

	cmds, err = s.redis.Pipelined(ctx, func(p gored.Pipeliner) error {
		for _, t := range tokens {
			p.HGetAll(ctx, tokenKey(clientID, t))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	summaryBuckets := make(map[string]*GameSummary)
	playerCount := 0

	for _, c := range cmds {
		err := c.Err()
		if err != nil {
			if err == gored.Nil {
				continue
			}
			return nil, err
		}
		cmd, ok := c.(*gored.StringStringMapCmd)
		if !ok {
			return nil, errors.New("pipeline command cast error")
		}

		var p dbData
		err = cmd.Scan(&p)
		if err != nil {
			return nil, err
		}

		if p.IsTour == true {
			continue
		}

		d, err := toData(p)
		if err != nil {
			return nil, err
		}
		if byCurrency {
			_, ok := currencies[d.ActiveCurrencyCode]
			if !ok {
				continue
			}
		}
		if byGame {
			_, ok := games[d.GameCode]
			if !ok {
				continue
			}
		}
		if byParent {
			ok := checkParent(d)
			if !ok {
				continue
			}
		}

		gameCode := d.GameCode
		if len(gameCode) == 0 {
			continue
		}
		playerCount += 1
		gameSummary, ok := summaryBuckets[gameCode]
		if !ok {
			summaryBuckets[gameCode] = &GameSummary{
				GameCode:    gameCode,
				PlayerCount: 0,
			}
			gameSummary = summaryBuckets[gameCode]
		}
		gameSummary.PlayerCount += 1
	}

	gameSummaries := make([]GameSummary, 0)
	for _, sum := range summaryBuckets {
		gameSummaries = append(gameSummaries, *sum)
	}
	sort.Slice(gameSummaries, func(i, j int) bool {
		return gameSummaries[i].GameCode < gameSummaries[j].GameCode
	})
	rep := &SumOnlinePlayerResult{
		GameSummaries:    gameSummaries,
		TotalPlayerCount: int64(playerCount),
	}

	return rep, nil
}
