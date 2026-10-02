package token_player

import (
	gored "github.com/go-redis/redis/v8"
)

var luaSetActive = gored.NewScript(`      
	local key = KEYS[1]
	local now = ARGV[1]
	local gameCode = ARGV[2]
	local currencyCode = ARGV[3]

	if redis.call('exists', key) == 1 then 
		redis.call('hset', key, 'LastActive', now) 
		redis.call('hset', key, 'GameCode', gameCode) 
		redis.call('hset', key, 'ActiveCurrencyCode', currencyCode) 
		return {1}
	end
	return {0}
`)
