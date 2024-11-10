package mongo

import (
	"context"

	"github.com/r9odt/chef-webapi/database/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionUserName string = "users"

// CreateUser creates user with given parameters.
func (db *DBConnector) CreateUser(username, password,
	fullName string, isAdmin, isBlocked bool) (*interfaces.UserEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	user := interfaces.NewEmptyUser()

	user.Username = username
	user.Password = password
	user.FullName = fullName
	user.IsBlocked = isBlocked
	user.IsAdmin = isAdmin
	user.NeedPasswordChange = true

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return db.GetUserByUsername(username)
}

// GetUserByUsername return user by given name.
func (db *DBConnector) GetUserByUsername(name string) (
	*interfaces.UserEntry, error) {
	var result interfaces.UserEntry
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	filter := bson.D{primitive.E{Key: "username", Value: name}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return nil, err
	}

	result.ID = result.ObjectID
	return &result, nil
}

// GetUserByID return user by given id.
func (db *DBConnector) GetUserByID(id string) (*interfaces.UserEntry, error) {
	var result interfaces.UserEntry
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	objectID, err := primitive.ObjectIDFromHex(id)
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

	result.ID = result.ObjectID
	return &result, nil
}

// GetAllUsers return all users.
func (db *DBConnector) GetAllUsers() ([]interfaces.UserEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	result := make([]interfaces.UserEntry, 0)
	filter := bson.D{{}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return result, err
	}

	if err := cursor.All(context.TODO(), &result); err != nil {
		return result, err
	}
	for i := range result {
		result[i].ID = result[i].ObjectID
	}
	return result, nil
}

// DeleteUserByID delete user by given id.
func (db *DBConnector) DeleteUserByID(id string) error {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	objectID, err := primitive.ObjectIDFromHex(id)
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

// UpdateUserByID update user by given id.
func (db *DBConnector) UpdateUserByID(id string,
	entry *interfaces.UserEntry) error {
	collection := db.DB.Database(db.DatabaseName).Collection(collectionUserName)
	objectID, err := primitive.ObjectIDFromHex(id)
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
