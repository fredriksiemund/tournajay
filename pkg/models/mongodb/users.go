package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserModel struct {
	Coll *mongo.Collection
}

func (m *UserModel) Upsert(id, name, email, picture string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}
	update := bson.M{
		"$set": bson.M{
			"name":    name,
			"email":   email,
			"picture": picture,
		},
	}
	options := options.Update().SetUpsert(true)

	_, err := m.Coll.UpdateOne(ctx, filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) One(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	u := &models.User{}
	err := m.Coll.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
