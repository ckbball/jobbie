package mongo

import (
  "context"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "github.com/ckbball/quik"
)

type DB struct {
  ds *mongo.Collection //
}

func NewMongoCollection(dbName string, collName string, client *mongo.Client) *DB {
  collection := client.Database(dbName).Collection(collName)
  return &DB{
    ds: collection,
  }
}
