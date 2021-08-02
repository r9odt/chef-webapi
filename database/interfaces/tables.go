package interfaces

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// TimeFormat is string for convert timestamp to string.
const TimeFormat string = "Mon Jan 2 15:04:05"

// UserEntry describes a user record.
type UserEntry struct {
	ObjectID  string `bson:"_id,omitempty" json:"-"`
	ID        string `json:"id" bson:"-"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"-" bson:"password"`
	FullName  string `json:"fullName" bson:"fullName"`
	Avatar    string `json:"avatar" bson:"avatar"`
	IsAdmin   bool   `json:"admin" bson:"admin"`
	IsBlocked bool   `json:"blocked" bson:"blocked"`
	LastLogin string `json:"lastLogin" bson:"lastLogin"`
	NeedPasswordChange bool `json:"needPasswordChange" bson:"needPasswordChange"`
	LastSeen  string `json:"lastSeen" bson:"-"`
}

// GetBSOND Returns Object as BSON.d
func (v *UserEntry) GetBSOND() bson.D {
	return bson.D{
		{Key: "username", Value: v.Username},
		{Key: "avatar", Value: v.Avatar},
		{Key: "password", Value: v.Password},
		{Key: "blocked", Value: v.IsBlocked},
		{Key: "admin", Value: v.IsAdmin},
		{Key: "fullName", Value: v.FullName},
		{Key: "lastLogin", Value: v.LastLogin},
		{Key: "needPasswordChange", Value: v.NeedPasswordChange},
	}
}

// NewEmptyUser returns a new empty user record.
func NewEmptyUser() UserEntry {
	return UserEntry{
		ID:        "",
		Username:  "",
		Password:  "",
		FullName:  "",
		Avatar:    "",
		IsAdmin:   false,
		IsBlocked: false,
		LastLogin: "Never logged in",
		LastSeen: "Never seen yet",
		NeedPasswordChange: false,
	}
}

// CopyUserData returns a copy of the user's record.
func CopyUserData(v UserEntry) *UserEntry {
	return &UserEntry{
		ID:        v.ID,
		Username:  v.Username,
		Password:  v.Password,
		FullName:  v.FullName,
		Avatar:    v.Avatar,
		IsAdmin:   v.IsAdmin,
		IsBlocked: v.IsBlocked,
		LastLogin: v.LastLogin,
		LastSeen: "Never seen yet",
		NeedPasswordChange: v.NeedPasswordChange,
	}
}

// SessionEntry describes the user session record.
type SessionEntry struct {
	ObjectID string `bson:"_id,omitempty" json:"-"`
	ID       string `json:"id" bson:"-"`
	UUID     string `json:"uuid" bson:"-"`
	Username string `json:"username" bson:"username"`
	Expire   int64  `json:"expire" bson:"expire"`
}

// GetBSOND Returns Object as BSON.d
func (v *SessionEntry) GetBSOND() bson.D {
	return bson.D{
		{Key: "username", Value: v.Username},
		{Key: "expire", Value: v.Expire},
	}
}

// NewEmptySession returns a new empty session record.
func NewEmptySession() SessionEntry {
	return SessionEntry{
		ID:       "",
		UUID:     "",
		Username: "",
		Expire:   -1,
	}
}

// CopySessionData returns a copy of the session record.
func CopySessionData(v SessionEntry) *SessionEntry {
	return &SessionEntry{
		ID:       v.ID,
		UUID:     v.UUID,
		Username: v.Username,
		Expire:   v.Expire,
	}
}

// TaskEntry describes a record of the task.
type TaskEntry struct {
	ObjectID         string `bson:"_id,omitempty" json:"-"`
	ID               string `json:"id" bson:"-"`
	Resource         string `json:"resource" bson:"resource"`
	Name             string `json:"name" bson:"name"`
	Status           string `json:"status" bson:"status"`
	InitiatorID      string `json:"initiatorID" bson:"initiatorID"`
	Date             string `json:"date" bson:"-"`
	Timestamp        int64  `json:"timestamp" bson:"timestamp"`
	Log              string `json:"log,omitempty" bson:"log"`
	OnlyResource     bool   `json:"onlyResource" bson:"onlyResource"`
	Resources        string `json:"resources" bson:"resources"`
	SelectedResource bool   `json:"selectedResource" bson:"selectedResource"`
}

// GetBSOND Returns Object as BSON.d
func (v *TaskEntry) GetBSOND() bson.D {
	return bson.D{
		{Key: "resource", Value: v.Resource},
		{Key: "name", Value: v.Name},
		{Key: "status", Value: v.Status},
		{Key: "initiatorID", Value: v.InitiatorID},
		{Key: "date", Value: v.Date},
		{Key: "timestamp", Value: v.Timestamp},
		{Key: "log", Value: v.Log},
		{Key: "onlyResource", Value: v.OnlyResource},
		{Key: "resources", Value: v.Resources},
		{Key: "selectedResource", Value: v.SelectedResource},
	}
}

// NewEmptyTask returns a new empty task record.
func NewEmptyTask() *TaskEntry {
	return &TaskEntry{
		ID:               "",
		Resource:         "",
		Name:             "",
		Status:           "Waiting",
		InitiatorID:      "",
		Timestamp:        -1,
		Log:              "",
		OnlyResource:     false,
		Resources:        "",
		SelectedResource: false,
	}
}

// CopyTaskData returns a copy of the task record.
func CopyTaskData(v *TaskEntry) *TaskEntry {
	return &TaskEntry{
		ID:               v.ID,
		Resource:         v.Resource,
		Name:             v.Name,
		InitiatorID:      v.InitiatorID,
		Status:           v.Status,
		Log:              v.Log,
		Date:             time.Unix(v.Timestamp, 0).Format(TimeFormat),
		Timestamp:        v.Timestamp,
		OnlyResource:     v.OnlyResource,
		Resources:        v.Resources,
		SelectedResource: v.SelectedResource,
	}
}

// AppModuleEntry describes the user session record.
type AppModuleEntry struct {
	ObjectID string `bson:"_id,omitempty" json:"-"`
	ID       string `json:"id" bson:"-"`
	Name     string `json:"name" bson:"name"`
	Comment  string `json:"comment" bson:"comment"`
	IsON     bool   `json:"isON" bson:"isON"`
}

// GetBSOND Returns Object as BSON.d
func (v *AppModuleEntry) GetBSOND() bson.D {
	return bson.D{
		{Key: "name", Value: v.Name},
		{Key: "comment", Value: v.Comment},
		{Key: "isON", Value: v.IsON},
	}
}

// NewEmptyAppModule returns a new empty session record.
func NewEmptyAppModule() AppModuleEntry {
	return AppModuleEntry{
		ID:      "",
		Name:    "",
		IsON:    false,
		Comment: "",
	}
}

// CopyAppModuleData returns a copy of the session record.
func CopyAppModuleData(v AppModuleEntry) *AppModuleEntry {
	return &AppModuleEntry{
		ID:      v.ID,
		Name:    v.Name,
		IsON:    v.IsON,
		Comment: v.Comment,
	}
}
