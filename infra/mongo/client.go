package mongo

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/iamrz1/ab-auth/infra"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo holds necessery fields and
// mongo database session to connect
type Mongo struct {
	*mongo.Client
	database *mongo.Database
	name     string
	lgr      logger.Logger
}

// New returns a new instance of mongodb using session s
func New(ctx context.Context, uri, name string, timeout time.Duration, opts ...Option) (*Mongo, error) {
	minPoolSize := uint64(10)
	maxPoolSize := uint64(100)
	connectionOption := &options.ClientOptions{
		SocketTimeout:          &timeout,
		ConnectTimeout:         &timeout,
		MaxPoolSize:            &maxPoolSize,
		MinPoolSize:            &minPoolSize,
		ServerSelectionTimeout: &timeout,
		RetryWrites:            utils.TrueP(),
		ReadPreference:         readpref.Secondary(),
		//ReplicaSet:             nil,
		//Direct:                 nil,
	}
	log.Println("hitting mongo connect...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri), connectionOption)
	if err != nil {
		return nil, err
	}

	log.Println("pinging db...")
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	log.Println("mongo connected...")

	db := &Mongo{
		Client:   client,
		database: client.Database(name),
		name:     name,
	}
	for _, opt := range opts {
		opt.apply(db)
	}
	return db, nil
}

// Option is mongo db option
type Option interface {
	apply(*Mongo)
}

// OptionFunc implements Option interface
type OptionFunc func(db *Mongo)

func (f OptionFunc) apply(db *Mongo) {
	f(db)
}

// SetLogger sets logger
func SetLogger(lgr logger.Logger) Option {
	return OptionFunc(func(db *Mongo) {
		db.lgr = lgr
	})
}

func (d *Mongo) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx, readpref.Primary())
}

func (d *Mongo) Close(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

// EnsureIndices creates indices for collection col
func (d *Mongo) EnsureIndices(ctx context.Context, collection string, inds []infra.DbIndex) error {
	d.lgr.Infof("EnsureIndices", "", "creating indices for", collection)
	db := d.database
	indexModels := []mongo.IndexModel{}
	for _, ind := range inds {
		keys := bson.D{}
		for _, k := range ind.Keys {
			keys = append(keys, bson.E{k.Key, k.Asc})
		}
		opts := options.Index()
		if ind.Unique != nil {
			opts.SetUnique(*ind.Unique)
		}
		if ind.Sparse != nil {
			opts.SetSparse(*ind.Sparse)
		}
		if ind.Name != "" {
			opts.SetName(ind.Name)
		}
		if ind.ExpireAfter != nil {
			opts.SetExpireAfterSeconds(int32(ind.ExpireAfter.Seconds()))
		}
		im := mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		}
		indexModels = append(indexModels, im)
	}
	if _, err := db.Collection(collection).Indexes().CreateMany(ctx, indexModels); err != nil {
		return err
	}
	return nil
}

// DropIndices drops indices from collection col
func (d *Mongo) DropIndices(ctx context.Context, collection string, inds []infra.DbIndex) error {
	d.lgr.Infof("DropIndices", "", "dropping indices from", collection)
	if _, err := d.database.Collection(collection).Indexes().DropAll(ctx); err != nil {
		return err
	}
	return nil
}

// Insert inserts doc into collection
func (d *Mongo) Insert(ctx context.Context, collection string, doc interface{}) error {
	d.lgr.Infof("Insert", "", "insert into", collection)
	if _, err := d.database.Collection(collection).InsertOne(ctx, doc); err != nil {
		return err
	}
	return nil
}

// Update updates existing doc in the collection
func (d *Mongo) Update(ctx context.Context, collection string, filter, doc interface{}) (int64, error) {
	d.lgr.Infof("Update", "", "update in", collection)
	update := bson.M{"$set": doc}
	res, err := d.database.Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {

		return 0, err
	}
	return res.MatchedCount, nil
}

// InsertMany inserts docs into collection
func (d *Mongo) InsertMany(ctx context.Context, collection string, docs []interface{}) error {
	d.lgr.Infof("InsertMany", "", "insert many into", collection)
	if _, err := d.database.Collection(collection).InsertMany(ctx, docs); err != nil {
		return err
	}
	return nil
}

// FindOne finds a doc by query
func (d *Mongo) FindOne(ctx context.Context, collection string, q interface{}, v interface{}, sort ...interface{}) error {
	d.lgr.Infof("FindOne", "", "find", q, "from", collection)
	findOneOpts := options.FindOne()
	if len(sort) > 0 {
		findOneOpts = findOneOpts.SetSort(sort[0])
	}
	err := d.database.Collection(collection).FindOne(ctx, q, findOneOpts).Decode(v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return infra.ErrNotFound
		}
		return err
	}
	return nil
}

// Count counts documents
func (d *Mongo) Count(ctx context.Context, collection string, filter interface{}) (int64, error) {
	d.lgr.Infof("Count", "", "count", filter, "from", collection)

	n, err := d.database.Collection(collection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// FindAndCount counts from found results
func (d *Mongo) FindAndCount(ctx context.Context, collection string, filter interface{}) (int64, error) {
	d.lgr.Infof("FindAndCount", "", "count", filter, "from", collection)

	cursor, err := d.database.Collection(collection).Find(ctx, filter)
	if err != nil {
		return 0, err
	}

	count := int64(cursor.RemainingBatchLength())
	cursor.Close(ctx)
	return count, nil
}

// List finds list of docs that matches query with skip and limit
func (d *Mongo) List(ctx context.Context, collection string, filter interface{}, page, limit int64, v interface{}, sort ...interface{}) error {
	d.lgr.Infof("List", "", "list", filter, "from", collection)
	skip := (page - 1) * limit
	findOpts := options.Find().SetSkip(skip).SetLimit(limit)
	if len(sort) > 0 {
		findOpts = findOpts.SetSort(sort[0])
	}
	cursor, err := d.database.Collection(collection).Find(ctx, filter, findOpts)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}

	cursor.RemainingBatchLength()

	return nil
}

// Aggregate runs aggregation q on docs and store the result on v
func (d *Mongo) Aggregate(ctx context.Context, collection string, q interface{}, v interface{}) error {
	d.lgr.Infof("Aggregate", "", "aggregate", q, "from", collection)
	cursor, err := d.database.Collection(collection).Aggregate(ctx, q)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}
	return nil
}

func (d *Mongo) AggregateWithDiskUse(ctx context.Context, collection string, q []infra.DbQuery, v interface{}) error {
	d.lgr.Infof("AggregateWithDiskUse", "", "aggregate", q, "from", collection)
	opt := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := d.database.Collection(collection).Aggregate(ctx, q, opt)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}
	return nil
}

func (d *Mongo) Distinct(ctx context.Context, collection, field string, q infra.DbQuery, v interface{}) error {
	d.lgr.Infof("Distinct", "", "aggregate", q, "from", collection)
	interfaces, err := d.database.Collection(collection).Distinct(ctx, field, q)
	if err != nil {
		return err
	}
	data, err := json.Marshal(interfaces)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (d *Mongo) PartialUpdateMany(ctx context.Context, collection string, filter infra.DbQuery, data interface{}) error {
	_, err := d.database.Collection(collection).UpdateMany(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return err
	}
	return nil
}

func (d *Mongo) PartialUpdateManyByQuery(ctx context.Context, collection string, filter infra.DbQuery, query infra.UnorderedDbQuery) error {
	_, err := d.database.Collection(collection).UpdateMany(ctx, filter, query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Mongo) BulkUpdate(ctx context.Context, collection string, models []mongo.WriteModel) error {
	_, err := d.database.Collection(collection).BulkWrite(ctx, models)
	return err
}

func (d *Mongo) DeleteMany(ctx context.Context, collection string, filter interface{}) error {
	_, err := d.database.Collection(collection).DeleteMany(ctx, filter)
	return err
}

func (d *Mongo) DeleteOne(ctx context.Context, collection string, filter interface{}) (int64, error) {
	res, err := d.database.Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, err
}
