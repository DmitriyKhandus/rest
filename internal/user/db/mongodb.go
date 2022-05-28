package db

import (
	"context"
	"fmt"

	"github.com/DmitriyKhandus/rest-api/internal/user"
	"github.com/DmitriyKhandus/rest-api/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug(("create user"))
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user %v", user)
	}
	d.logger.Debug("convert insertedId to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert oid to hex %s", oid)
}
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to converthex to ObjId. hex=%s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return u, fmt.Errorf("failed to find user by id: %s", id)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode result userId= %s", id)
	}
	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return fmt.Errorf("failed to convert userId to objId. id=%s", user.Id)
	}
	filter := bson.M{"_id": objectId}
	userBytes, err := bson.Marshal(user)

	if err != nil {
		return fmt.Errorf("failed to marshal user error %v", err)
	}
	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes %v", err)
	}

	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query error %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("not found user")
	} else {
		d.logger.Tracef("found counts %v updated %v", result.MatchedCount, result.ModifiedCount)
	}

	return nil

}

func (d *db) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert userId to objId. id=%s", id)
	}
	filter := bson.M{"_id": objectId}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query error: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("not found user")
	}
	d.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}