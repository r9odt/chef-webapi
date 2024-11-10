package database

import (
	"github.com/r9odt/chef-webapi/database/interfaces"
)

// CreateSession creates session with given parameters.
func (db *DBConnector) CreateSession(username string, expire int64) (*interfaces.SessionEntry, error) {
	return db.SessionDatabase.CreateSession(username, expire)
}

// GetSessionByUUID return session by given uuid.
func (db *DBConnector) GetSessionByUUID(uuid string) (*interfaces.SessionEntry, error) {
	return db.SessionDatabase.GetSessionByUUID(uuid)
}

// UpdateSessionByUUID update session by given uuid.
func (db *DBConnector) UpdateSessionByUUID(uuid string, user *interfaces.SessionEntry) error {
	return db.Database.UpdateSessionByUUID(uuid, user)
}

// DeleteSessionByUUID delete session by given uuid.
func (db *DBConnector) DeleteSessionByUUID(uuid string) error {
	return db.SessionDatabase.DeleteSessionByUUID(uuid)
}

// DeleteExpireSession delete all expire session.
func (db *DBConnector) DeleteExpireSession() error {
	return db.SessionDatabase.DeleteExpireSession()
}
