package dao

import (
	"context"
	mgutil "coolcar/shared/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

// Mongo defines a mongo dao
type Mongo struct {
	col *mongo.Collection
	//newObjID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{col: db.Collection("account")}
}

func (m *Mongo) ResolveAccountID(ctx context.Context, openID string) (string, error) {
	//m.col.InsertOne(ctx, bson.M{
	//	mgutil.IDFieldName: m.newObjID(),
	//	openIDField: openID,
	//})
	insertdID := mgutil.NewObjID()
	result := m.col.FindOneAndUpdate(ctx,
		bson.M{
			openIDField: insertdID,
		},
		mgutil.SetOnInsert(bson.M{
			mgutil.IDFieldName: insertdID,
			openIDField:        openID,
		}),
		options.FindOneAndUpdate().
			SetUpsert(true).SetReturnDocument(options.After))
	if err := result.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}
	var row mgutil.IDField
	err := result.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}
	return row.ID.Hex(), nil
}
