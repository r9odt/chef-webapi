package mongo

import (
	"context"

	"github.com/JIexa24/chef-webapi/database/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionAppModuleName string = "modules"

// CreateAppModule creates module with given parameters.
func (db *DBConnector) CreateAppModule(name, comment string,
	initialStatus bool) (*interfaces.AppModuleEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionAppModuleName)
	m, err := db.GetAppModuleByName(name)
	if err != nil {
		return nil, err
	}
	if m != nil {
		return nil, nil
	}
	module := interfaces.NewEmptyAppModule()

	module.Name = name
	module.Comment = comment
	module.IsON = initialStatus

	r, err := collection.InsertOne(context.TODO(), module)
	if err != nil {
		return nil, err
	}

	id := r.InsertedID.(primitive.ObjectID)
	return db.GetAppModuleByID(id.Hex())
}

// GetAllAppModules return all modules.
func (db *DBConnector) GetAllAppModules() ([]interfaces.AppModuleEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionAppModuleName)
	result := make([]interfaces.AppModuleEntry, 0)
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

// GetAppModuleByID return module by given id.
func (db *DBConnector) GetAppModuleByID(id string) (*interfaces.AppModuleEntry, error) {
	var result interfaces.AppModuleEntry
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionAppModuleName)
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

// GetAppModuleByName return module by given name.
func (db *DBConnector) GetAppModuleByName(name string) (*interfaces.AppModuleEntry, error) {
	var result interfaces.AppModuleEntry
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionAppModuleName)
	filter := bson.D{primitive.E{Key: "name", Value: name}}
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

// UpdateAppModuleByID update module by given id.
func (db *DBConnector) UpdateAppModuleByID(id string,
	entry *interfaces.AppModuleEntry) error {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionAppModuleName)
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
