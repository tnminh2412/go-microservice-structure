package queue

import (
	"go-pkg/rabbitmq/shard/publisher"
	"go-pkg/rabbitmq/shard/subscriber"
	"testing"

	uuid2 "github.com/google/uuid"
)

func Test_queueManager_writeToFile(t *testing.T) {
	type fields struct {
		subs map[string]*subscriber.Subscriber
		pubs map[string]*publisher.Publisher
	}
	type args struct {
		logFolder string
		data      interface{}
	}

	uuid, _ := uuid2.NewUUID()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "write-to-file",
			fields: fields{
				subs: nil,
				pubs: nil,
			},
			args: args{
				logFolder: uuid.String(),
				data: &struct {
					Test string
				}{
					Test: "Name",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qm := &queueManager{
				subs: tt.fields.subs,
				pubs: tt.fields.pubs,
			}
			if err := qm.WriteToFile(tt.args.logFolder, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("writeToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
