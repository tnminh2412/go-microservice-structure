package test

import (
	"context"
	"go-pkg/redis"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/player/token_player"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func BenchmarkSet(b *testing.B) {

	ctx := context.Background()
	tk := token_player.New(redis.GetUniversal("Test"))
	errCount := int64(0)
	successCount := int64(0)
	for i := 0; i < b.N; i++ {
		err := tk.Set(ctx, &token_player.SetData{
			ClientID: "demo",
			Username: uuid.NewString(),
			Token:    uuid.NewString(),
			ParentL1: "a",
			ParentL2: "b",
		})
		if err != nil {
			atomic.AddInt64(&errCount, 1)
		} else {
			atomic.AddInt64(&successCount, 1)
		}
	}

	b.Log("err count ", errCount)
	b.Log("success count ", successCount)

}

func BenchmarkSum(b *testing.B) {

	ctx := context.Background()
	tk := token_player.New(redis.GetUniversal("Test"))
	errCount := int64(0)
	successCount := int64(0)
	for i := 0; i < b.N; i++ {
		r, err := tk.SumOnlinePlayer(ctx, "demo", token_player.SumOnlineParams{
			MinActiveGameTime: time.Now().Add(-1 * time.Hour),
		})
		if err != nil {
			b.Log(err)
			atomic.AddInt64(&errCount, 1)
		} else {
			b.Log(r)
			atomic.AddInt64(&successCount, 1)
		}
	}

	b.Log("err count ", errCount)
	b.Log("success count ", successCount)

}

func BenchmarkList(b *testing.B) {

	ctx := context.Background()
	tk := token_player.New(redis.GetUniversal("Test"))
	errCount := int64(0)
	successCount := int64(0)
	for i := 0; i < b.N; i++ {
		r, err := tk.ListOnlinePlayer(ctx, "demo", token_player.ListOnlineParams{
			Limit: 10,
		})
		if err != nil {
			b.Log(err)
			atomic.AddInt64(&errCount, 1)
		} else {
			b.Log(r)
			atomic.AddInt64(&successCount, 1)
		}
	}

	b.Log("err count ", errCount)
	b.Log("success count ", successCount)

}

func TestManySet(t *testing.T) {

	setCount := 50_000
	ctx := context.Background()
	tk := token_player.New(redis.GetUniversal("Test"))
	errCount := int64(0)
	successCount := int64(0)
	successRunDurationMillies := int64(0)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(200)
	for i := 0; i < setCount; i++ {
		g.Go(func() error {
			start := time.Now()
			token := uuid.NewString()
			err := tk.Set(gCtx, &token_player.SetData{
				ClientID: "demo",
				Username: uuid.NewString(),
				Token:    token,
				ParentL1: "a",
				ParentL2: "f",
			})
			if err != nil {
				atomic.AddInt64(&errCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
				dur := time.Since(start)
				atomic.AddInt64(&successRunDurationMillies, int64(dur.Milliseconds()))
			}
			tk.SetActive(gCtx, "demo", token, "123", "ABC")
			return nil
		})

	}

	t.Log("err count ", errCount)
	t.Log("success count ", successCount)
	t.Log("avg success millies ", float64(successRunDurationMillies)/float64(successCount))

}
