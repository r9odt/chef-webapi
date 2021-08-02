package interfaces

// Database describes the database object interface.
type Database interface {
	// REST
	CreateUser(string, string, string, bool, bool) (*UserEntry, error)
	GetAllUsers() ([]UserEntry, error)
	GetUserByID(string) (*UserEntry, error)
	DeleteUserByID(string) error
	UpdateUserByID(string, *UserEntry) error
	// REST extension
	GetUserByUsername(string) (*UserEntry, error)

	CreateSession(string, int64) (*SessionEntry, error)
	GetSessionByUUID(string) (*SessionEntry, error)
	UpdateSessionByUUID(string, *SessionEntry) error
	DeleteSessionByUUID(string) error
	DeleteExpireSession() error

	// REST
	CreateTask(string, string, string, string, bool, bool) (*TaskEntry, error)
	GetAllTasks() ([]TaskEntry, error)
	GetTaskByID(string) (*TaskEntry, error)
	UpdateTaskStatusByID(string, string) error
	// REST extension
	СheckIfTaskAlreadyCreate(string, string) bool
	GetWaitingTask() (*TaskEntry, error)
	GetLastCompleteTaskByResourceAndName(string, string) (*TaskEntry, error)
	UpdateTaskStatusAtStartup() error

	// REST
	CreateAppModule(string, string, bool) (*AppModuleEntry, error)
	GetAllAppModules() ([]AppModuleEntry, error)
	GetAppModuleByID(string) (*AppModuleEntry, error)
	UpdateAppModuleByID(string, *AppModuleEntry) error
	// REST extension
	GetAppModuleByName(string) (*AppModuleEntry, error)
	
	UploadFile(string, string) error
	DownloadFile(string) []byte
	CheckFile(string) bool

	Close()
}
