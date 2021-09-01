package dao

import (
	"context"
	mgo "coolcar/shared/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

// Mongo defines a mongo dao
type Mongo struct {
	col      *mongo.Collection
	newObjID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{col: db.Collection("account"),
		newObjID: primitive.NewObjectID}
}

func (m *Mongo) ResolveAccountID(ctx context.Context, openID string) (string, error) {
	//m.col.InsertOne(ctx, bson.M{
	//	mgo.IDField: m.newObjID(),
	//	openIDField: openID,
	//})

	result := m.col.FindOneAndUpdate(ctx,
		bson.M{
			openIDField: openID,
		},
		mgo.Set(bson.M{
			"open_id": openID,
		}),
		options.FindOneAndUpdate().
			SetUpsert(true).SetReturnDocument(options.After))
	if err := result.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}
	var row struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := result.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}
	return row.ID.Hex(), nil
}
