package infra

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// DB interface wraps the databse
type DB interface {
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
	EnsureIndices(ctx context.Context, collection string, inds []DbIndex) error
	DropIndices(ctx context.Context, collection string, inds []DbIndex) error
	Insert(ctx context.Context, collection string, doc interface{}) error
	Update(ctx context.Context, collection string, filter, doc interface{}) (int64, error)
	InsertMany(ctx context.Context, collection string, v []interface{}) error
	List(ctx context.Context, collection string, filter interface{}, page, limit int64, v interface{}, sort ...interface{}) error
	FindOne(ctx context.Context, collection string, filter interface{}, v interface{}, sort ...interface{}) error
	PartialUpdateMany(ctx context.Context, collection string, filter DbQuery, data interface{}) error
	PartialUpdateManyByQuery(ctx context.Context, collection string, filter DbQuery, query UnorderedDbQuery) error
	BulkUpdate(ctx context.Context, collection string, models []mongo.WriteModel) error
	Aggregate(ctx context.Context, collection string, q interface{}, v interface{}) error
	AggregateWithDiskUse(ctx context.Context, collection string, q []DbQuery, v interface{}) error
	Distinct(ctx context.Context, collection, field string, q DbQuery, v interface{}) error
	DeleteMany(ctx context.Context, collection string, filter interface{}) error
	DeleteOne(ctx context.Context, collection string, filter interface{}) (int64, error)
	Count(ctx context.Context, collection string, filter interface{}) (int64, error)
	FindAndCount(ctx context.Context, collection string, filter interface{}) (int64, error)
}

// DbIndex holds database index
type DbIndex struct {
	Name   string
	Keys   []DbIndexKey
	Unique *bool
	Sparse *bool

	// If ExpireAfter is defined the server will periodically delete
	// documents with indexed time.Time older than the provided delta.
	ExpireAfter *time.Duration
}

type DbIndexKey struct {
	Key string
	Asc interface{}
}

// DbQuery holds a database query
type DbQuery bson.D
type UnorderedDbQuery bson.M

type BulkWriteModel mongo.WriteModel
