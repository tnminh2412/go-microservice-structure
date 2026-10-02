package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	db2 "github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/db"
	"github.com/tnminh2412/go-microservice-structure/service/app/pkg/shared_conf"
	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUpdateRole(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	name := uuid.New().String()
	type args struct {
		ctx     context.Context
		mongodb *mongo.Database
		role    *db2.Role
	}
	db2.CreateRole(context.Background(), db, &db2.Role{
		Name:  name,
		Tasks: []string{"task2"},
		Locale: []*db2.Locale{
			&db2.Locale{
				Language:    "en",
				Description: "hello",
			},
		},
		CreatedBy:   "create-by",
		CreatedTime: time.Now(),
		UpdatedBy:   "update-by",
		UpdatedTime: time.Now(),
	})

	tests := []struct {
		name        string
		args        args
		wantMatched int64
		wantErr     bool
	}{
		{
			name: "update-existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				role: &db2.Role{
					Name:  name,
					Tasks: []string{"task-1"},
					Locale: []*db2.Locale{
						&db2.Locale{
							Language:    "en",
							Description: "hello",
						},
					},
					CreatedBy:   "create-by",
					CreatedTime: time.Now(),
					UpdatedBy:   "update-by",
					UpdatedTime: time.Now(),
				},
			},
			wantMatched: 1,
			wantErr:     false,
		}, {
			name: "update-not-existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				role: &db2.Role{
					Name:  uuid.New().String(),
					Tasks: []string{"task-1"},
					Locale: []*db2.Locale{
						&db2.Locale{
							Language:    "en",
							Description: "hello",
						},
					},
					CreatedBy:   "create-by",
					CreatedTime: time.Now(),
					UpdatedBy:   "update-by",
					UpdatedTime: time.Now(),
				},
			},
			wantMatched: 0,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatched, err := db2.UpdateRole(tt.args.ctx, tt.args.mongodb, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil {
				if gotMatched != tt.wantMatched {
					t.Errorf("UpdateRole() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
				}
			}
		})
	}
}

func TestCreateRole(t *testing.T) {
	globalService := global_conf.New(context.Background())
	confService := shared_conf.New(context.Background(), globalService)

	db, _ := confService.GetDatabase(context.Background(), "demo")

	name := uuid.New().String()
	type args struct {
		ctx     context.Context
		mongodb *mongo.Database
		role    *db2.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create-new",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				role: &db2.Role{
					Name:  name,
					Tasks: []string{"task-1"},
					Locale: []*db2.Locale{
						&db2.Locale{
							Language:    "en",
							Description: "hello",
						},
					},
					CreatedBy:   "create-by",
					CreatedTime: time.Now(),
					UpdatedBy:   "update-by",
					UpdatedTime: time.Now(),
				},
			},
			wantErr: false,
		}, {
			name: "create-existing",
			args: args{
				ctx:     context.Background(),
				mongodb: db,
				role: &db2.Role{
					Name:  name,
					Tasks: []string{"task-1"},
					Locale: []*db2.Locale{
						&db2.Locale{
							Language:    "en",
							Description: "hello",
						},
					},
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
			if err := db2.CreateRole(tt.args.ctx, tt.args.mongodb, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("CreateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRole(t *testing.T) {
	type args struct {
		ctx      context.Context
		mongodb  *mongo.Database
		roleName string
	}
	tests := []struct {
		name    string
		args    args
		want    *db2.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db2.GetRole(tt.args.ctx, tt.args.mongodb, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListRole(t *testing.T) {
	type args struct {
		ctx     context.Context
		mongodb *mongo.Database
	}
	tests := []struct {
		name    string
		args    args
		want    []*db2.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db2.ListRole(tt.args.ctx, tt.args.mongodb, db2.ListRoleParams{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}
