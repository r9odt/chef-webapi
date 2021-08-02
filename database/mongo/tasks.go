package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/JIexa24/chef-webapi/database/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionTaskName string = "tasks"

// CreateTask creates task with given parameters.
func (db *DBConnector) CreateTask(resource, name, resources,
	initiatorID string, onlyResource, selectedResource bool) (*interfaces.TaskEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
	t := time.Now()

	entry := interfaces.NewEmptyTask()
	entry.Resource = resource
	entry.Name = name
	entry.InitiatorID = initiatorID
	entry.Timestamp = t.Unix()
	entry.OnlyResource = onlyResource
	entry.Resources = resources
	entry.SelectedResource = selectedResource

	result, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		return nil, err
	}
	return db.GetTaskByID(result.InsertedID.(primitive.ObjectID).Hex())
}

// СheckIfTaskAlreadyCreate checks for the existence of a task with an error or completed status.
func (db *DBConnector) СheckIfTaskAlreadyCreate(resource, name string) bool {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
	result := make([]interfaces.TaskEntry, 0)
	filter := bson.M{
		"$and": []bson.M{
			{"resource": bson.M{"$eq": resource}},
			{"name": bson.M{"$eq": name}},
			{"status": bson.M{"$ne": "Complete"}},
			{"status": bson.M{"$ne": "Error"}},
		},
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return false
	}

	if err := cursor.All(context.TODO(), &result); err != nil {
		return false
	}
	return len(result) > 0
}

// GetAllTasks return all tasks.
func (db *DBConnector) GetAllTasks() ([]interfaces.TaskEntry, error) {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
	result := make([]interfaces.TaskEntry, 0)
	filter := bson.D{{}}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
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
		result[i].Date = time.Unix(result[i].Timestamp,
			0).Format(interfaces.TimeFormat)
	}
	return result, nil
}

// GetWaitingTask return one task with 'waiting' status.
func (db *DBConnector) GetWaitingTask() (*interfaces.TaskEntry, error) {
	var result interfaces.TaskEntry
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)

	filter := bson.M{"status": "Waiting"}
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

// GetTaskByID return task by given id.
func (db *DBConnector) GetTaskByID(id string) (*interfaces.TaskEntry, error) {
	var result interfaces.TaskEntry
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
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
	result.Date = time.Unix(result.Timestamp, 0).Format(interfaces.TimeFormat)
	return &result, nil
}

// UpdateTaskStatusByID update task status by given id.
func (db *DBConnector) UpdateTaskStatusByID(id, newStatus string) error {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
	entry, err := db.GetTaskByID(id)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("task not found")
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	entry.ObjectID = entry.ID
	entry.Status = newStatus
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

// UpdateTaskStatusAtStartup update tasks with status 'InProgress' to 'Error'.
func (db *DBConnector) UpdateTaskStatusAtStartup() error {
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)

	filter := bson.M{"status": "InProgress"}
	update := bson.M{"$set": bson.M{"status": "Error"}}

	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return err
	}

	return nil
}

// GetLastCompleteTaskByResourceAndName return last completed task for given resource.
func (db *DBConnector) GetLastCompleteTaskByResourceAndName(resource,
	name string) (*interfaces.TaskEntry, error) {
	var result interfaces.TaskEntry
	collection := db.DB.Database(db.DatabaseName).Collection(
		collectionTaskName)
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	filter := bson.M{"resource": resource, "name": name, "status": "Complete"}
	err := collection.FindOne(context.TODO(), filter, findOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return nil, err
	}

	result.ID = result.ObjectID
	result.Date = time.Unix(result.Timestamp, 0).Format(interfaces.TimeFormat)
	return &result, nil
}
