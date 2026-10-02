package shared_conf

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	"github.com/tnminh2412/go-microservice-structure/service/pkg/global_conf"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Value struct {
	Name string
}

type Serialize struct {
}

func (s *Serialize) Decode(_ context.Context, result *mongo.SingleResult) (interface{}, error) {
	var _conf struct {
		Value *Value
	}
	err := result.Decode(&_conf)
	return _conf.Value, err
}

func (s *Serialize) Deserialize(_ context.Context, bytes []byte) (interface{}, error) {
	var _conf struct {
		Value *Value
	}
	err := json.Unmarshal(bytes, &_conf)
	return _conf.Value, err
}

func TestService_GetConf(t *testing.T) {
	ctx := context.Background()
	configKey := uuid.New().String()
	name := uuid.New().String()

	value := &Value{
		Name: name,
	}

	globalService := global_conf.New(ctx)
	sharedService := New(ctx, globalService)

	db, _ := globalService.Client(context.Background(), "demo")
	db.Database("app_client_demo").Collection("conf").InsertOne(context.Background(), bson.M{"_id": configKey, "value": value})

	type fields struct {
		GlobalService global_conf.IService
	}
	type args struct {
		ctx      context.Context
		clientID string
		key      string
		s        SerializeInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test-db-config",
			fields: fields{
				GlobalService: globalService,
			},
			args: args{
				ctx:      context.Background(),
				clientID: "demo",
				key:      configKey,
				s:        &Serialize{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sharedService.GetConf(tt.args.ctx, tt.args.clientID, tt.args.key, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err == nil {
				if v, ok := got.(*Value); !ok {
					t.Errorf("GetConf() error = %v, wantErr %v", err, tt.wantErr)
					return
				} else if v.Name != name {
					t.Errorf("GetConf() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func TestService_GetDatabase(t *testing.T) {
	ctx := context.Background()
	type fields struct {
	}
	type args struct {
		ctx      context.Context
		clientID string
	}

	globalService := global_conf.New(ctx)
	sharedService := New(ctx, globalService)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.Database
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sharedService.GetDatabase(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDatabase() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_confKey(t *testing.T) {
	type args struct {
		clientID string
		key      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := confKey(tt.args.clientID, tt.args.key); got != tt.want {
				t.Errorf("confKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UpdateContact(t *testing.T) {
	type fields struct {
		globalService   global_conf.IService
		context         context.Context
		stopFunc        context.CancelFunc
		subscribers     map[string]struct{}
		subscribersLock sync.RWMutex
		configs         map[string]map[string]interface{}
		configLock      sync.RWMutex
		connections     map[string]struct{}
	}
	type args struct {
		ctx         context.Context
		clientID    string
		contactInfo *ContactInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				globalService:   nil,
				context:         context.Background(),
				stopFunc:        nil,
				subscribers:     nil,
				subscribersLock: sync.RWMutex{},
				configs:         nil,
				configLock:      sync.RWMutex{},
				connections:     nil,
			},
			args: args{
				ctx:      context.Background(),
				clientID: "demo",
				contactInfo: &ContactInfo{
					UseDefault: map[string]bool{
						"VN": false,
						"MY": true,
					},
					Contacts: []*Contact{
						{
							Type:        "type",
							Value:       "value",
							ID:          "id",
							Icon:        "icon",
							Url:         "url",
							IsShowQR:    false,
							Enabled:     false,
							Order:       0,
							CountryCode: "countryCode",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(context.Background(), global_conf.New(context.Background()))
			if err := service.UpdateContact(tt.args.ctx, tt.args.clientID, tt.args.contactInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
