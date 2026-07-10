package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"turtle/core/lgr"
	"turtle/credentials"
)

func ExportMongo(connstr string) {
	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI(connstr) // Replace with your MongoDB URI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	another_client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		lgr.Error(err.Error())
		return
	}

	target := another_client.Database(credentials.GetDBName())
	source := DB.Db(credentials.GetDBName())

	MigrateMongo(target, source)
}

func ImportMongo(connstr string) {
	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI(connstr) // Replace with your MongoDB URI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	another_client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		lgr.Error(err.Error())
		return
	}

	source := another_client.Database(credentials.GetDBName())
	target := DB.Db(credentials.GetDBName())

	MigrateMongo(target, source)
}

func MigrateMongo(target *mongo.Database, source *mongo.Database) {

	ctx := context.TODO()

	for _, container := range DB.ListContainers() {
		target_ct := target.Collection(container)

		cursor, err := source.Collection(container).Find(ctx, bson.M{})

		if err != nil {
			lgr.Error(err.Error())
			continue
		}

		for cursor.Next(context.TODO()) {
			var elem map[string]interface{}

			err := cursor.Decode(&elem)
			if err != nil {
				lgr.Error(err.Error())
				continue
			}

			uid_ptr, _ := elem["uid"]

			if uid_ptr != nil {
				uid := uid_ptr.(string)
				if val, err := target_ct.CountDocuments(ctx, bson.M{"uid": uid}); err != nil && val != 0 {
					target_ct.UpdateOne(ctx, bson.M{"uid": uid}, elem)
				} else {
					target_ct.InsertOne(ctx, elem)
				}
			} else {
				target_ct.InsertOne(ctx, elem)
			}

		}
	}
}
