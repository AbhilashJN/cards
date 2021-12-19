package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection interface {
	InsertOne(context.Context, interface{},
		...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}
