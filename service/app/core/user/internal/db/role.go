package db

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Role struct {
	Name         string    `bson:"_id"`
	Tasks        []string  `bson:"tasks"`
	Locale       []*Locale `bson:"locale"`
	CreatedBy    string    `bson:"createdBy"`
	CreatedTime  time.Time `bson:"createdTime"`
	UpdatedBy    string    `bson:"updatedBy"`
	UpdatedTime  time.Time `bson:"updatedTime"`
	ExcludeTasks []string  `bson:"excludeTasks"`
}

type Locale struct {
	Language    string `bson:"language"`
	Description string `bson:"description"`
}

const roleCollectionName = "role"

func CreateRole(ctx context.Context, mongodb *mongo.Database, role *Role) error {
	_, err := mongodb.Collection(roleCollectionName).InsertOne(ctx, role, &options.InsertOneOptions{
		BypassDocumentValidation: nil,
		Comment:                  nil,
	})

	return err
}

func UpdateRole(ctx context.Context, mongodb *mongo.Database, role *Role) (matched int64, err error) {
	set := bson.M{}

	if len(role.Tasks) > 0 {
		set["tasks"] = role.Tasks
	}

	if len(role.Locale) > 0 {
		set["locale"] = role.Locale
	}

	if role.UpdatedBy != "" {
		set["updatedBy"] = role.UpdatedBy
	}

	if !role.UpdatedTime.IsZero() {
		set["updatedTime"] = role.UpdatedTime
	}

	updated, err := mongodb.Collection(roleCollectionName).UpdateOne(ctx, bson.M{"_id": role.Name}, bson.M{"$set": set})

	if err != nil {
		return 0, err
	}

	return updated.MatchedCount, err
}

type ListRoleParams struct {
	ExcludeSystem bool
}

func ListRole(ctx context.Context, mongodb *mongo.Database, p ListRoleParams) ([]*Role, error) {

	filter := bson.M{}
	if p.ExcludeSystem {
		filter["isSystem"] = false
	}
	cursor, err := mongodb.Collection(roleCollectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var reply []*Role

	for cursor.Next(ctx) {
		var role Role
		if err := cursor.Decode(&role); err != nil {
			context_helper.Logger(ctx).Error(func() string {
				return fmt.Sprintf("fetch role error %v", err)
			})
			return nil, err
		} else {
			reply = append(reply, &role)
		}
	}

	return reply, nil
}

func GetRole(ctx context.Context, mongodb *mongo.Database, roleName string) (*Role, error) {
	var role Role
	err := mongodb.Collection(roleCollectionName).FindOne(ctx, bson.M{"_id": roleName}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}
