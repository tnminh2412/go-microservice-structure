package queue

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	service_error "go-pkg/helper/error_helper"
	"go-pkg/helper/json_helper"
	"go-pkg/helper/local_file_helper"
	"go-pkg/log"
	"go-pkg/log/base"
	"go-pkg/log/file"
	"go-pkg/rabbitmq/shard"
	"go-pkg/rabbitmq/shard/publisher"
	"go-pkg/rabbitmq/shard/subscriber"
	sr "go-pkg/rabbitmq/subscriber"
	"go-pkg/service"
	"path/filepath"
)

type Message struct {
	*sr.Payload
}

type Info struct {
	Subscribers []SubscriberInfo
	Publishers  []PublisherInfo
}

type SubscriberInfo struct {
	Key       string
	DataType  interface{}
	ConsumeFn ConsumeWithBatchFn
}

type PublisherInfo struct {
	Key string
}

type ConsumeWithBatchFn func(ctx context.Context, processMessages []Message) (rejectedMessages []Message)

type IQueueManager interface {
	Init(Info)
	Subscribe(key string) error
	Publish(key string, data interface{}, headers map[string]interface{}) error
	SafePublish(key string, data interface{}, headers map[string]interface{}, logDirectory string) error
	Stop()
	WriteToFile(logFolder string, data interface{}) error
}

type queueManager struct {
	subs map[string]*subscriber.Subscriber
	pubs map[string]*publisher.Publisher
}

func New() IQueueManager {
	return &queueManager{
		subs: make(map[string]*subscriber.Subscriber, 0),
		pubs: make(map[string]*publisher.Publisher, 0),
	}
}

func (qm *queueManager) Init(info Info) {
	if info.Subscribers != nil {
		qm.InitSubscribers(info.Subscribers)
	}
	if info.Publishers != nil {
		qm.InitPublishers(info.Publishers)
	}
}

func (qm *queueManager) InitPublishers(infos []PublisherInfo) {
	for _, info := range infos {
		if _, ok := qm.pubs[info.Key]; ok == true {
			return
		}

		logger := file.NewLogger()
		ctx := context.Background()
		pubConnection, err := shard.NewPub(logger, ctx)
		if err != nil {
			log.Error(func() string {
				return fmt.Sprintf("shard.NewPub, create Pub connection failed, err: %v", err.Error())
			})
			panic(err)
		}
		pub, err := publisher.New(logger, ctx, pubConnection, info.Key)
		if err != nil {
			log.Error(func() string {
				return fmt.Sprintf("publisher.New, new publisher failed, err: %v", err.Error())
			})
			panic(err)
		}

		qm.pubs[info.Key] = pub
	}
}

func (qm *queueManager) InitSubscribers(infos []SubscriberInfo) {
	for _, info := range infos {
		if _, ok := qm.subs[info.Key]; ok == true {
			return
		}

		logger := file.NewLogger()
		ctx := context.Background()
		subConnection, err := shard.NewSub(logger, ctx)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return "shard.NewSub, create Sub connection failed"
			}, base.Error(err))
			panic(err)
		}
		sub, err := newSubscriberWithBatch(logger, ctx, subConnection, info.Key, info.DataType, info.ConsumeFn)
		if err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return fmt.Sprintf("newSubscribe, newSubscribe failed, err: %v", err.Error())
			})
			panic(err)
		}

		qm.subs[info.Key] = sub
	}
}

func (qm *queueManager) Subscribe(subscribeKey string) error {
	sub, ok := qm.subs[subscribeKey]
	if ok == false {
		return Error(service_error.NotFound, fmt.Sprintf("Not found subscriber: %s", subscribeKey), nil)
	}

	err := sub.Start()
	if err != nil {
		return Error(service_error.Internal, "sub.Start error", err)
	}

	return nil
}

func (qm *queueManager) Publish(key string, data interface{}, headers map[string]interface{}) error {
	pub, ok := qm.pubs[key]
	if ok == false {
		return Error(service_error.NotFound, fmt.Sprintf("Not found publisher: %s", key), nil)
	}

	_, err := pub.Publish(data, headers)
	if err != nil {
		return Error(service_error.Internal, "pub.Publish error", err)
	}

	return nil
}

const dataDirectory = "/data"

func (qm *queueManager) SafePublish(key string, data interface{}, headers map[string]interface{}, logFolder string) error {
	err := qm.Publish(key, data, headers)
	if err != nil {
		return qm.WriteToFile(logFolder, data)
	}
	return nil
}

func (qm *queueManager) WriteToFile(logFolder string, data interface{}) error {
	dirName := filepath.Join(dataDirectory, logFolder)
	return local_file_helper.WriteFile(dirName, service.HostName, json_helper.SafeJsonToString(data))
}

func (qm *queueManager) Stop() {
	for key, sub := range qm.subs {
		err := sub.Close()
		if err != nil {
			context_helper.Logger(context.Background()).Error(func() string {
				return fmt.Sprintf("sub.Close() error, subscribeKey: %s", key)
			}, base.Error(err))
		}
	}
}

func newSubscriberWithBatch(logger base.Logger, ctx context.Context, subConnection *shard.Shard, subscribeKey string, dataType interface{}, consumeFn ConsumeWithBatchFn) (*subscriber.Subscriber, error) {
	sub, err := subscriber.NewWithBatch(logger, ctx, subConnection, subscribeKey, dataType, func(datas []*sr.Payload) []uint64 {

		var msgs []Message
		for _, payload := range datas {
			msgs = append(msgs, Message{
				Payload: payload,
			})
		}

		rejectedMessages := consumeFn(ctx, msgs)
		rejectedTag := make([]uint64, len(rejectedMessages))
		for i, rejectedMsg := range rejectedMessages {
			rejectedTag[i] = rejectedMsg.GetTag()
		}

		return rejectedTag
	}, true)

	return sub, err
}

func (m Message) GetBodyData() interface{} {
	bd := m.GetBody()
	return bd
}
