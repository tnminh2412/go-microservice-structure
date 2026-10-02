package test

import (
	"context"
	"testing"
	"time"

	db2 "github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateUser(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	playerId := uuid.New().String()

	type args struct {
		ctx     context.Context
		mongodb *mongo.Database
		user    *db2.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "insert-new",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				user: &db2.User{
					Username:    playerId,
					Roles:       []string{"role1"},
					Password:    "pass",
					Status:      "active",
					CreatedBy:   "create-by",
					CreatedTime: time.Now(),
					UpdatedBy:   "update-by",
					UpdatedTime: time.Now(),
				},
			},
			wantErr: false,
		}, {
			name: "insert-existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				user: &db2.User{
					Username:    playerId,
					Roles:       []string{"role1"},
					Password:    "pass",
					Status:      "active",
					CreatedBy:   "create-by",
					CreatedTime: time.Now(),
					UpdatedBy:   "update-by",
					UpdatedTime: time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db2.CreateUser(tt.args.ctx, tt.args.mongodb, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	playerId := uuid.New().String()
	db2.CreateUser(context.Background(), db, &db2.User{
		Username:    playerId,
		Roles:       []string{"role1"},
		Password:    "password",
		Status:      "active",
		CreatedBy:   "create-by",
		CreatedTime: time.Now(),
		UpdatedBy:   "update-by",
		UpdatedTime: time.Time{},
	})
	type args struct {
		ctx      context.Context
		mongodb  *mongo.Database
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *db2.User
		wantErr bool
	}{
		{
			name: "get-existing",
			args: args{
				ctx:      context.Background(),
				mongodb:  db,
				username: playerId,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "get-not-existing",
			args: args{
				ctx:      context.Background(),
				mongodb:  db,
				username: uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db2.GetUser(tt.args.ctx, tt.args.mongodb, db2.GetUserParams{ID: tt.args.username})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil {
				if got.Username != tt.args.username {
					t.Errorf("GetUser() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestListUser(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	type args struct {
		ctx     context.Context
		mongodb *mongo.Database
	}
	tests := []struct {
		name    string
		args    args
		want    []*db2.User
		wantErr bool
	}{
		{
			name: "list-existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db2.ListUser(tt.args.ctx, tt.args.mongodb, db2.ListUserParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUpdateUser(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	playerId := uuid.New().String()

	db2.CreateUser(context.Background(), db, &db2.User{
		Username:    playerId,
		Roles:       []string{"role1"},
		Password:    "password",
		Status:      "active",
		CreatedBy:   "create-by",
		CreatedTime: time.Now(),
		UpdatedBy:   "update-by",
		UpdatedTime: time.Time{},
	})

	type args struct {
		ctx      context.Context
		mongodb  *mongo.Database
		user     *db2.User
		modified int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				user: &db2.User{
					Username:    playerId,
					Roles:       []string{"test"},
					Password:    "pwd-",
					Status:      "active",
					UpdatedBy:   "update",
					UpdatedTime: time.Now(),
				},
				modified: 1,
			},
			wantErr: false,
		},
		{
			name: "update not existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				user: &db2.User{
					Username:    uuid.New().String(),
					Roles:       []string{"test"},
					Password:    "pwd-",
					Status:      "active",
					UpdatedBy:   "update",
					UpdatedTime: time.Now(),
				},
				modified: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if modified, err := db2.UpdateUser(tt.args.ctx, tt.args.mongodb, &db2.UpdateUserParams{User: tt.args.user}); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				if modified != tt.args.modified {
					t.Errorf("UpdateUser() modified %v", modified)
				}
			}
		})
	}
}
