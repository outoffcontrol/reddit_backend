package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoDatabase interface {
	Collection(name string) IMongoCollection
}

type IMongoCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (IMongoDeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (IMongoUpdateResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (MongoInsertOneResult, error)
}

type IMongoSingleResult interface {
	Decode(v interface{}) error
}
type IMongoDeleteResult interface {
}

type IMongoUpdateResult interface {
	UnmarshalBSON(b []byte) error
}

type IMongoCursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(ctx context.Context, results interface{}) error
}

type MongoCollection struct {
	Coll *mongo.Collection
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

type MongoDeleteResult struct {
	dr *mongo.DeleteResult
}
type MongoUpdateResult struct {
	ur *mongo.UpdateResult
}
type MongoInsertOneResult struct {
	ir *mongo.InsertOneResult
}
type MongoCursor struct {
	cur *mongo.Cursor
}

func (mur *MongoUpdateResult) UnmarshalBSON(b []byte) error {
	return mur.ur.UnmarshalBSON(b)
}

func (msr *MongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
}

func (mc *MongoCursor) Close(ctx context.Context) error {
	return mc.cur.Close(ctx)
}

func (mc *MongoCursor) Next(ctx context.Context) bool {
	return mc.cur.Next(ctx)
}

func (mc *MongoCursor) Decode(val interface{}) error {
	return mc.cur.Decode(val)
}
func (mc *MongoCursor) All(ctx context.Context, results interface{}) error {
	return mc.cur.All(ctx, results)
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error) {
	cursorResult, err := mc.Coll.Find(ctx, filter, opts...)
	return &MongoCursor{cur: cursorResult}, err
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult {
	singleResult := mc.Coll.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}

func (mc *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (IMongoDeleteResult, error) {
	deleteResult, err := mc.Coll.DeleteOne(ctx, filter, opts...)
	return &MongoDeleteResult{dr: deleteResult}, err
}

func (mc *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (IMongoUpdateResult, error) {
	updateResilt, err := mc.Coll.UpdateOne(ctx, filter, update, opts...)
	return &MongoUpdateResult{ur: updateResilt}, err
}
func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (MongoInsertOneResult, error) {
	insertResilt, err := mc.Coll.InsertOne(ctx, document, opts...)
	return MongoInsertOneResult{ir: insertResilt}, err
}
