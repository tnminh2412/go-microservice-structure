package global_conf

import (
	"context"
	"testing"
)

func TestGetClientConf(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *Conf
		wantErr bool
	}{
		{
			name: "get-conf",
			args: args{
				ctx:      context.Background(),
				clientID: "demo",
			},
			want:    nil,
			wantErr: false,
		},
	}
	service := &Service{
		subscribers: make(map[string]struct{}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetClientConf(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetClientTemplate(t *testing.T) {
	type args struct {
		ctx    context.Context
		Domain string
	}
	tests := []struct {
		name    string
		args    args
		want    *Conf
		wantErr bool
	}{
		{
			name: "get-template",
			args: args{
				ctx:    context.Background(),
				Domain: "localhost:25000",
			},
			want:    nil,
			wantErr: false,
		},
	}
	service := &Service{
		subscribers: make(map[string]struct{}),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetTemplate(tt.args.ctx, tt.args.Domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
