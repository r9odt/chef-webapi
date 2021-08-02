package web

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"

	"github.com/JIexa24/chef-webapi/database"
	"github.com/JIexa24/chef-webapi/logging"

	"github.com/go-chef/chef"
)

// App is a global application object.
var App *Application

// NewApplication set empty application structure.
func NewApplication(env string, reloadChannel chan struct{}) {
	App = &Application{
		Client:          nil,
		DB:              nil,
		AppKey:          nil,
		Logger:          nil,
		WorkerDirectory: "/app/worker/data",
		SSHKeyPath:      "",
		ChefKeyPath:     "",
		chefClientName:  "",
		chefURL:         "",
		chefMux:         sync.Mutex{},
		Env:             env,
		SessionExpire:   86400,
		LDAP:            nil,
		UsersLastSeen:   make(map[string]string),
		ReloadChannel:   reloadChannel,
	}
}

// ConfigureLDAP configures ldap.
func (App *Application) ConfigureLDAP(baseDN, bindAddress, bindPrefix, bindSuffix string) {
	App.LDAP = &LDAPData{
		BaseDN:      baseDN,
		BindAddress: bindAddress,
		BindPrefix:  bindPrefix,
		BindSuffix:  bindSuffix,
	}
}

// ConfigureLogger configures logger.
func (App *Application) ConfigureLogger(l logging.Logger) {
	App.Logger = l
}

// GetChefClientConfig return existing chef client config.
func (App *Application) GetChefClientConfig() *chef.Client {
	App.chefMux.Lock()
	defer App.chefMux.Unlock()
	return App.Client
}

// CreateChefClientConfig creates config for chef client.
func (App *Application) CreateChefClientConfig() error {
	App.chefMux.Lock()
	defer App.chefMux.Unlock()
	// read a client key
	key, err := os.ReadFile(App.ChefKeyPath)
	if err != nil {
		return err
	}

	App.Client, err = chef.NewClient(&chef.Config{
		Name:    App.chefClientName,
		Key:     string(key),
		BaseURL: App.chefURL,
	})

	return err
}

// ConfigureChefClient configures chef client.
func (App *Application) ConfigureChefClient(name, url, keyPath string) error {
	App.ChefKeyPath = keyPath
	App.chefClientName = name
	App.chefURL = url
	return App.CreateChefClientConfig()
}

// ConfigureDatabase configures database.
func (App *Application) ConfigureDatabase(databaseProvider, sessionProvider,
	databaseName, databaseUser, databasePassword,
	databaseHost, databasePort string) error {
	var err error
	databaseParams := database.NewParams()
	databaseParams.DatabaseProvider = databaseProvider
	databaseParams.SessionProvider = sessionProvider
	databaseParams.User = databaseUser
	databaseParams.Password = databasePassword
	databaseParams.Host = databaseHost
	databaseParams.Port = databasePort
	databaseParams.Name = databaseName
	App.DB, err = database.NewDatabaseConnector(App.Env, databaseParams, App.Logger)
	return err
}

// ConfigureApp configures App.
func (App *Application) ConfigureApp(workerDir, appKeyPath, sshKeyPath string, sessionExpire int64) {
	App.SessionExpire = sessionExpire
	App.WorkerDirectory = workerDir
	App.SSHKeyPath = sshKeyPath
	data, err := os.ReadFile(appKeyPath)
	if err != nil {
		App.Logger.Fatal(err.Error())
	}
	App.AppKey, err = parseRsaPrivateKey(data)
	if err != nil {
		App.Logger.Fatal(err.Error())
	}
}

func parseRsaPrivateKey(str []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(str)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

// ConcatenateStringWithDelimeter constructing string like
// 'names[0],names[1], ... ,names[n]' if delimeter as ","
func ConcatenateStringWithDelimeter(delimeter string, names []string) string {
	var result string
	var length = len(names)
	for i, v := range names {
		result = result + v
		if i < length-1 {
			result = result + delimeter
		}
	}
	return result
}
