package database

import "github.com/JIexa24/chef-webapi/database/interfaces"

// CreateTask creates task with given parameters.
func (db *DBConnector) CreateTask(Resource, Name,
	resources, initiatorID string, onlyResource, selectedResource bool) (*interfaces.TaskEntry, error) {
	return db.Database.CreateTask(Resource, Name,
		resources, initiatorID, onlyResource, selectedResource)
}

// GetWaitingTask return one task with 'waiting' status.
func (db *DBConnector) GetWaitingTask() (*interfaces.TaskEntry, error) {
	return db.Database.GetWaitingTask()
}

// GetAllTasks return all tasks.
func (db *DBConnector) GetAllTasks() ([]interfaces.TaskEntry, error) {
	return db.Database.GetAllTasks()
}

// UpdateTaskStatusByID update task status by given id.
func (db *DBConnector) UpdateTaskStatusByID(id string, newStatus string) error {
	return db.Database.UpdateTaskStatusByID(id, newStatus)
}

// СheckIfTaskAlreadyCreate checks for the existence of a task with an error or completed status.
func (db *DBConnector) СheckIfTaskAlreadyCreate(resource, name string) bool {
	return db.Database.СheckIfTaskAlreadyCreate(resource, name)
}

// UpdateTaskStatusAtStartup update tasks with status 'InProgress' to 'Error'.
func (db *DBConnector) UpdateTaskStatusAtStartup() error {
	return db.Database.UpdateTaskStatusAtStartup()
}

// GetLastCompleteTaskByResourceAndName return last completed task for given resource.
func (db *DBConnector) GetLastCompleteTaskByResourceAndName(resource, name string) (*interfaces.TaskEntry, error) {
	return db.Database.GetLastCompleteTaskByResourceAndName(resource, name)
}

// GetTaskByID return task by given id.
func (db *DBConnector) GetTaskByID(id string) (*interfaces.TaskEntry, error) {
	return db.Database.GetTaskByID(id)
}
