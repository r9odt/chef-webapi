package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/r9odt/chef-webapi/database/interfaces"
	"github.com/r9odt/chef-webapi/database/mongo"
	"github.com/r9odt/chef-webapi/encryption"
	"github.com/r9odt/chef-webapi/logging"
)

// Params describes database parameters.
type Params struct {
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	DatabaseProvider string
	SessionProvider  string
}

// ConstructDSN from parameters into string.
func (d *Params) ConstructDSN() string {
	dsn := ""
	switch d.DatabaseProvider {
	case "mongo":
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
			d.User, d.Password, d.Host, d.Port, d.Name)
	default:
	}
	return dsn
}

// NewParams return empty database params.
func NewParams() *Params {
	return &Params{
		Host:             "unknown",
		Port:             "unknown",
		User:             "unknown",
		Password:         "unknown",
		Name:             "unknown",
		DatabaseProvider: "unknown",
		SessionProvider:  "unknown",
	}
}

// DBConnector describes the structure of the database and
// implements the Database interface.
type DBConnector struct {
	Database        interfaces.Database
	SessionDatabase interfaces.Database
	Params  *Params
	StopRequest     chan struct{}
	waitGroup       sync.WaitGroup
	pause           time.Duration
}

// NewDatabaseConnector return database connector with session connector.
// Initialize memory db if databaseProvider in-memory.
func NewDatabaseConnector(env string, dbParams *Params, l logging.Logger) (*DBConnector, error) {
	db := &DBConnector{
		Database:        nil,
		SessionDatabase: nil,
		Params:  dbParams,
		StopRequest:     make(chan struct{}),
		waitGroup:       sync.WaitGroup{},
		pause:           time.Hour,
	}
	dsn := db.Params.ConstructDSN()
	switch db.Params.DatabaseProvider {
	case "mongo":
		database, err := mongo.New(dsn, db.Params.Name, l)
		if err != nil {
			return nil, err
		}
		db.Database = database
		db.initialize(l)
	default:
		return nil, fmt.Errorf("unknown provider %s",
			db.Params.DatabaseProvider)
	}
	switch db.Params.SessionProvider {
	case "database", db.Params.DatabaseProvider:
		db.SessionDatabase = db.Database
	default:
		return nil, fmt.Errorf("unknown provider %s",
			db.Params.SessionProvider)
	}
	db.waitGroup.Add(1)
	go db.sessionWatcher(l)
	return db, nil
}

// Close connector.
func (db *DBConnector) Close() {
	close(db.StopRequest)
	db.waitGroup.Wait()
	switch db.Params.DatabaseProvider {
	case "mongo":
		db.Database.Close()
	default:
	}
}

func (db *DBConnector) initialize(l logging.Logger) {
	err := db.CreateAppModules()
	if err != nil {
		l.Fatalf("%s", err.Error())
	}
	users, err := db.GetAllUsers()
	if err != nil {
		l.Fatalf("%s", err.Error())
	}
	if len(users) < 1 {
		admpwd, err := encryption.GetPasswordHASH([]byte("admin"))
		if err != nil {
			l.Fatalf("Database initialize: %s", err.Error())
		}
		testpwd, err := encryption.GetPasswordHASH([]byte("tester"))
		if err != nil {
			l.Fatalf("Database initialize: %s", err.Error())
		}
		user, err := db.CreateUser("admin", string(admpwd),
			"Admin^)", true, false)
		if err == nil && user != nil {
			user.Avatar = "https://svoi.sibnet.ru/photos/80x80/1399830_0.jpg"
			if err := db.UpdateUserByID(user.ID, user); err != nil {
				l.Errorf("Database initialize [db.UpdateUserByID]: %s", err.Error())
			}
		}
		if err != nil {
			l.Errorf("Database initialize: %s", err.Error())
		}
		_, err = db.CreateUser("tester", string(testpwd),
			"Just mortal",
			false, false)
		if err != nil {
			l.Errorf("Database initialize: %s", err.Error())
		}
	}
}

func (db *DBConnector) sessionWatcher(l logging.Logger) {
	defer db.waitGroup.Done()
	// Flow check request for completion of work.
	exit := false
	go func() {
		<-db.StopRequest
		exit = true
	}()
	timeoutChannel := make(chan bool)
	for !exit {
		if err := db.Database.DeleteExpireSession(); err != nil {
			l.Errorf("Database DeleteExpireSession: %s", err.Error())
		}
		go func() {
			// Waiting one of two signals.
			select {
			// Pause not to load the CPU.
			case <-time.After(db.pause):
			case <-db.StopRequest:
			}
			timeoutChannel <- true
		}()
		<-timeoutChannel
	}
}

// UploadFile uploads file to database.
func (db *DBConnector) UploadFile(path, filename string) error {
	return db.Database.UploadFile(path, filename)
}

// DownloadFile downloads file from database.
func (db *DBConnector) DownloadFile(filename string) []byte {
	return db.Database.DownloadFile(filename)
}

// CheckFile check file in database.
func (db *DBConnector) CheckFile(filename string) bool {
	return db.Database.CheckFile(filename)
}
