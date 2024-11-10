package mongo

import (
	"context"
	"time"

	"github.com/r9odt/chef-webapi/database/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionSessionName string = "sessions"

// CreateSession creates session with given parameters.
func (db *DBConnector) CreateSession(username string, expire int64) (*interfaces.SessionEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionSessionName)
	session := interfaces.NewEmptySession()

	session.Username = username
	session.Expire = expire

	r, err := collection.InsertOne(context.TODO(), session)
	if err != nil {
		return nil, err
	}

	uuid := r.InsertedID.(primitive.ObjectID)
	return db.GetSessionByUUID(uuid.Hex())
}

// GetSessionByUUID return session by given uuid.
func (db *DBConnector) GetSessionByUUID(uuid string) (*interfaces.SessionEntry, error) {
	var result interfaces.SessionEntry
	collection := db.DB.Database(db.DatabaseName).Collection(collectionSessionName)
	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return nil, err
	}

	result.UUID = result.ObjectID
	return &result, nil
}

// UpdateSessionByUUID update session by given uuid.
func (db *DBConnector) UpdateSessionByUUID(uuid string, entry *interfaces.SessionEntry) error {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionSessionName)
	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return err
	}

	entry.ObjectID = entry.ID
	filter := bson.M{"_id": objectID}
	_, err = collection.UpdateOne(context.TODO(), filter, bson.D{
		{Key: "$set", Value: entry.GetBSOND()},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return err
	}

	return nil
}

// DeleteSessionByUUID delete session by given uuid.
func (db *DBConnector) DeleteSessionByUUID(uuid string) error {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionSessionName)
	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return err
	}

	return nil
}

// DeleteExpireSession delete all expire session.
func (db *DBConnector) DeleteExpireSession() error {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionSessionName)
	filter := bson.M{"expire": bson.M{"$lt": time.Now().Unix()}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return err
	}

	return nil
}
