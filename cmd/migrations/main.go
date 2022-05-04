package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/db"
	"github.com/fredriksiemund/tournament-planner/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type schema struct {
	name string
	json primitive.M
}

var schemas = []*schema{
	{name: "users", json: mongodb.UserSchema},
}

func main() {
	down := flag.Bool("down", false, "Clear database")
	connStr := flag.String("connStr", "mongodb://root:root@localhost:27017", "MongoDb connection string")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := db.MongoConnect(*connStr, ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	database := client.Database("tournajay")

	if *down {
		collections, err := database.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			log.Println(1)
			log.Fatal(err)
		}
		for _, coll := range collections {
			err = database.Collection(coll).Drop(ctx)
			if err != nil {
				log.Println(2)
				log.Fatal(err)
			}
		}
	} else {
		for _, v := range schemas {
			validator := bson.M{"$jsonSchema": v.json}
			opts := options.CreateCollection().SetValidator(validator)
			err = database.CreateCollection(context.Background(), v.name, opts)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
