package dao

import (
	"context"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var mongoURI string

func TestMongo_ResolveAccountID(t *testing.T) {
	ctx := context.Background()
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	// 固定objectid
	mgutil.NewObjID = func() primitive.ObjectID {
		objID, _ := primitive.ObjectIDFromHex("jbojbojbo")
		return objID
	}
	id, err := m.ResolveAccountID(ctx, "123")
	if err != nil {
		t.Fatalf("failed resolve account id for 123: %v", err)
	} else {
		want := "xxxxxxxx"
		if id != want {
			t.Fatalf("resolve account id: want: %q, got: %q", want, id)
		}
	}

}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
