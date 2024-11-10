package web

import (
	"crypto/rsa"
	"sync"

	"github.com/r9odt/chef-webapi/database"
	"github.com/r9odt/chef-webapi/logging"

	"github.com/go-chef/chef"
)

// Application describes application.
type Application struct {
	Env             string
	WorkerDirectory string
	AppKey          *rsa.PrivateKey
	SessionExpire   int64
	Logger          logging.Logger
	DB              *database.DBConnector
	Client          *chef.Client
	LDAP            *LDAPData
	UsersLastSeen   map[string]string
	SSHKeyPath      string
	ChefKeyPath     string
	chefClientName  string
	chefURL         string
	chefMux         sync.Mutex
	ReloadChannel   chan struct{}
	StopRequest     chan struct{}
}

// LDAPData Describes the details of the connection to the LDAP server.
// Binding by prefix-username-suffix.
type LDAPData struct {
	BaseDN      string
	BindAddress string
	BindPrefix  string
	BindSuffix  string
}