package database

import "github.com/r9odt/chef-webapi/database/interfaces"

// CreateAppModule creates module with given parameters.
func (db *DBConnector) CreateAppModule(name, comment string,
	initialStatus bool) (*interfaces.AppModuleEntry, error) {
	return db.Database.CreateAppModule(name, comment, initialStatus)
}

// GetAllAppModules return all modules.
func (db *DBConnector) GetAllAppModules() ([]interfaces.AppModuleEntry, error) {
	return db.Database.GetAllAppModules()
}

// GetAppModuleByID return module by given id.
func (db *DBConnector) GetAppModuleByID(id string) (*interfaces.AppModuleEntry, error) {
	return db.Database.GetAppModuleByID(id)
}

// GetAppModuleByName return module by given name.
func (db *DBConnector) GetAppModuleByName(name string) (*interfaces.AppModuleEntry, error) {
	return db.Database.GetAppModuleByName(name)
}

// UpdateAppModuleByID update module by given id.
func (db *DBConnector) UpdateAppModuleByID(id string, entry *interfaces.AppModuleEntry) error {
	return db.Database.UpdateAppModuleByID(id, entry)
}

// CreateAppModules creates app modules.
func (db *DBConnector) CreateAppModules() error {
	if _, err := db.Database.CreateAppModule("ChefWorker",
		"Is chef worker processing tasks", false); err != nil {
		return err
	}
	return nil
}
