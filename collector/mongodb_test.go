package collector

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mockMongoDBDriver struct {
	ChangeStreamData *mockCursor
	AggregateCursor  *mockCursor
}

type mockCursor struct {
	Data    []interface{}
	Current interface{}
}

func (cursor *mockCursor) Decode(val interface{}) error {
	switch val.(type) {
	case *AggregationResult:
		*val.(*AggregationResult) = cursor.Current.(AggregationResult)
	case *ChangeStreamEvent:
		*val.(*ChangeStreamEvent) = cursor.Current.(ChangeStreamEvent)
	}
	return nil
}

func (cursor *mockCursor) Next(ctx context.Context) bool {
	if len(cursor.Data) == 0 {
		return false
	}

	cursor.Current, cursor.Data = cursor.Data[0], cursor.Data[1:]
	return true
}

func (cursor *mockCursor) Close(ctx context.Context) error {
	return nil
}

func (mdb *mockMongoDBDriver) Connect(ctx context.Context, opts ...*options.ClientOptions) error {
	return nil
}

func (mdb *mockMongoDBDriver) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return nil
}

func (mdb *mockMongoDBDriver) Aggregate(ctx context.Context, db string, col string, pipeline bson.A) (Cursor, error) {
	return mdb.AggregateCursor, nil
}

func (mdb *mockMongoDBDriver) Watch(ctx context.Context, db string, col string, pipeline mongo.Pipeline) (Cursor, error) {
	return mdb.AggregateCursor, nil
}
