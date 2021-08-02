package database

import "github.com/JIexa24/chef-webapi/database/interfaces"

// CreateUser creates user with given parameters.
func (db *DBConnector) CreateUser(username, password,
	fullName string, isAdmin, isBlocked bool) (*interfaces.UserEntry, error) {
	return db.Database.CreateUser(username, password,
		fullName, isAdmin, isBlocked)
}

// GetUserByUsername return user by given name.
func (db *DBConnector) GetUserByUsername(name string) (*interfaces.UserEntry, error) {
	return db.Database.GetUserByUsername(name)
}

// GetUserByID return user by given id.
func (db *DBConnector) GetUserByID(id string) (*interfaces.UserEntry, error) {
	return db.Database.GetUserByID(id)
}

// UpdateUserByID update user by given id.
func (db *DBConnector) UpdateUserByID(id string, user *interfaces.UserEntry) error {
	return db.Database.UpdateUserByID(id, user)
}

// GetAllUsers return all users.
func (db *DBConnector) GetAllUsers() ([]interfaces.UserEntry, error) {
	return db.Database.GetAllUsers()
}

// DeleteUserByID delete user by given id.
func (db *DBConnector) DeleteUserByID(id string) error {
	return db.Database.DeleteUserByID(id)
}
