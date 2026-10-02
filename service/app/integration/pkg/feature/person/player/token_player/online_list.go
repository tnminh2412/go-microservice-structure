package token_player

import (
	"context"
	"errors"
	"time"

	gored "github.com/go-redis/redis/v8"
)

type Parent struct {
	Username string
	Level    int32
}
type ListOnlineParams struct {
	Skip              int64
	Limit             int64
	CurrencyCodes     []string
	GameCodes         []string
	Usernames         []string
	MinActiveGameTime time.Time
	Active            bool
	Parent            *Parent
}

func (s *Server) ListOnlinePlayer(ctx context.Context, clientID string, params ListOnlineParams) ([]*Data, error) {
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

	usernames, err := s.redis.ZRevRange(ctx, onlinePlayersKey(clientID), 0, -1).Result()

	if err != nil {
		return nil, err
	}

	if len(usernames) == 0 {
		return make([]*Data, 0), nil
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
		return make([]*Data, 0), nil
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

	players := make([]*Data, 0)
	for _, c := range cmds {
		err := c.Err()
		if err != nil {
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

		if d.LastActive.Before(params.MinActiveGameTime) {
			d.GameCode = ""
		}
		if d.GameCode == "" && params.Active {
			continue
		}

		players = append(players, d)
	}
	from := params.Skip
	to := params.Skip + params.Limit
	if from > int64(len(players))-1 {
		return make([]*Data, 0), nil
	}
	if to > int64(len(players)) {
		to = int64(len(players))
	}
	if from >= to {
		return make([]*Data, 0), nil
	}
	page := players[from:to]

	return page, nil
}

func getCheckParentFn(parent *Parent) (func(*Data) bool, bool) {
	if parent == nil || parent.Level > 15 || parent.Level <= 0 {
		return nil, false
	}
	switch parent.Level {
	case 1:
		{
			return func(d *Data) bool {
				return d.ParentL1 == parent.Username
			}, true
		}
	case 2:
		{
			return func(d *Data) bool {
				return d.ParentL2 == parent.Username
			}, true
		}
	case 3:
		{
			return func(d *Data) bool {
				return d.ParentL3 == parent.Username
			}, true
		}
	case 4:
		{
			return func(d *Data) bool {
				return d.ParentL4 == parent.Username
			}, true
		}
	case 5:
		{
			return func(d *Data) bool {
				return d.ParentL5 == parent.Username
			}, true
		}
	case 6:
		{
			return func(d *Data) bool {
				return d.ParentL6 == parent.Username
			}, true
		}
	case 7:
		{
			return func(d *Data) bool {
				return d.ParentL7 == parent.Username
			}, true
		}
	case 8:
		{
			return func(d *Data) bool {
				return d.ParentL8 == parent.Username
			}, true
		}
	case 9:
		{
			return func(d *Data) bool {
				return d.ParentL9 == parent.Username
			}, true
		}
	case 10:
		{
			return func(d *Data) bool {
				return d.ParentL10 == parent.Username
			}, true
		}
	case 11:
		{
			return func(d *Data) bool {
				return d.ParentL11 == parent.Username
			}, true
		}
	case 12:
		{
			return func(d *Data) bool {
				return d.ParentL12 == parent.Username
			}, true
		}
	case 13:
		{
			return func(d *Data) bool {
				return d.ParentL13 == parent.Username
			}, true
		}
	case 14:
		{
			return func(d *Data) bool {
				return d.ParentL14 == parent.Username
			}, true
		}
	case 15:
		{
			return func(d *Data) bool {
				return d.ParentL15 == parent.Username
			}, true
		}
	}

	return nil, false
}
