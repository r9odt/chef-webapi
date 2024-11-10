package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	web "github.com/r9odt/chef-webapi"

	"github.com/r9odt/chef-webapi/database/interfaces"
	"github.com/r9odt/chef-webapi/encryption"
	"github.com/r9odt/chef-webapi/httpserver/middleware"

	mergeSort "github.com/r9odt/go-mergeSort"
)

type userRequest struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	FullName           string `json:"fullName"`
	Avatar             string `json:"avatar"`
	IsAdmin            bool   `json:"admin"`
	IsBlocked          bool   `json:"blocked"`
	NeedPasswordChange bool   `json:"needPasswordChange"`
}

// GetUserBySession return struct with user information.
// Information get by user from session.
func GetUserBySession(r *http.Request) (*interfaces.UserEntry, error) {
	session := r.Header.Get(middleware.SessionHeader)
	if session == "" || session == "null" {
		return nil, nil
	}
	s, err := web.App.DB.GetSessionByUUID(session)
	if s == nil || err != nil {
		return nil, err
	}
	user, err := web.App.DB.GetUserByUsername(s.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByID return struct with user information.
// Information get by id.
func GetUserByID(r *http.Request) (*interfaces.UserEntry, error) {
	userID := middleware.GetID(r)
	user, err := web.App.DB.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	result := interfaces.CopyUserData(*user)
	if val, ok := web.App.UsersLastSeen[result.Username]; ok {
		result.LastSeen = val
	}

	return result, nil
}

// UpdateUserByID updates user information.
// Information update by id.
func UpdateUserByID(r *http.Request) error {
	userID := middleware.GetID(r)
	currentUser, err := GetUserBySession(r)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return fmt.Errorf("user from session not found")
	}
	user, err := web.App.DB.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	updatedUser := &userRequest{}
	err = json.Unmarshal(body, updatedUser)
	if err != nil {
		return err
	}
	ImplementToUserUpdatedInfo(currentUser.ID, user, updatedUser)
	return web.App.DB.UpdateUserByID(user.ID, user)
}

// ImplementToUserUpdatedInfo implement some updated fields into user's entry.
func ImplementToUserUpdatedInfo(currentUserID string, user *interfaces.UserEntry,
	updatedUser *userRequest) {
	if updatedUser == nil || user == nil {
		return
	}
	user.FullName = updatedUser.FullName
	user.Avatar = updatedUser.Avatar
	user.NeedPasswordChange = updatedUser.NeedPasswordChange

	// Prevent self ban for users.
	if currentUserID != updatedUser.ID {
		user.IsBlocked = updatedUser.IsBlocked
		user.IsAdmin = updatedUser.IsAdmin
	}

	// Update only if password has been received.
	if updatedUser.Password != "" {
		passwd, err := encryption.RSADecrypt(updatedUser.Password, *web.App.AppKey)
		if err == nil {
			var pwd []byte
			var err error
			if pwd, err = encryption.GetPasswordHASH([]byte(passwd)); err == nil {
				user.Password = string(pwd)
			} else {
				web.App.Logger.Errorf(
					"ImplementToUserUpdatedInfo [encryption.GetPasswordHASH]: %s",
					err.Error())
			}
		} else {
			web.App.Logger.Errorf(
				"ImplementToUserUpdatedInfo [encryption.RSADecrypt]: %s", err.Error())
		}
	}
}

// DeleteUserByID delete user from database.
func DeleteUserByID(r *http.Request) error {
	userID := middleware.GetID(r)
	user, err := web.App.DB.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	return web.App.DB.DeleteUserByID(user.ID)
}

// GetAllUsers return users information.
func GetAllUsers(parameters *APIParameters) ([]interfaces.UserEntry, error) {
	sortField := parameters.Sort
	sortOrder := parameters.Order
	query := parameters.Q
	idFilter := parameters.ID
	list, err := web.App.DB.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var sliceToSort = make([]mergeSort.Interface, 0)
	for i := range list {
		if strings.Contains(list[i].Username, query) {
			n := interfaces.CopyUserData(list[i])
			contain := false
			for i := range idFilter {
				if idFilter[i] == n.ID {
					contain = true
					break
				}
			}
			if len(idFilter) <= 0 || contain {
				if val, ok := web.App.UsersLastSeen[n.Username]; ok {
					n.LastSeen = val
				}
				sliceToSort = append(sliceToSort, *n)
			}
		}
	}

	compareUsernameFunction := func(a, b mergeSort.Interface) bool {
		return a.(interfaces.UserEntry).Username <
			b.(interfaces.UserEntry).Username
	}
	compareFullNameFunction := func(a, b mergeSort.Interface) bool {
		return a.(interfaces.UserEntry).FullName <
			b.(interfaces.UserEntry).FullName
	}
	compareLastLoginFunction := func(a, b mergeSort.Interface) bool {
		ta, _ := time.Parse(interfaces.TimeFormat, a.(interfaces.UserEntry).LastLogin)
		tb, _ := time.Parse(interfaces.TimeFormat, b.(interfaces.UserEntry).LastLogin)
		return ta.Unix() < tb.Unix()
	}
	compareLastSeenFunction := func(a, b mergeSort.Interface) bool {
		ta, _ := time.Parse(interfaces.TimeFormat, a.(interfaces.UserEntry).LastSeen)
		tb, _ := time.Parse(interfaces.TimeFormat, b.(interfaces.UserEntry).LastSeen)
		return ta.Unix() < tb.Unix()
	}

	switch sortField {
	default:
		sliceToSort = mergeSort.Sort(sliceToSort, compareUsernameFunction, false)
	case "username":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareUsernameFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareUsernameFunction, false)
		}
	case "fullName":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareFullNameFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareFullNameFunction, false)
		}
	case "lastLogin":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareLastLoginFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareLastLoginFunction, false)
		}
	case "lastSeen":
		switch sortOrder {
		case "DESC":
			sliceToSort = mergeSort.Sort(sliceToSort, compareLastSeenFunction, true)
		default:
			sliceToSort = mergeSort.Sort(sliceToSort, compareLastSeenFunction, false)
		}
	}

	var result []interfaces.UserEntry = make([]interfaces.UserEntry,
		len(sliceToSort))
	for i, v := range sliceToSort {
		result[i] = v.(interfaces.UserEntry)
	}

	return result, nil
}

// CreateUser creates user in database.
func CreateUser(r *http.Request) (*interfaces.UserEntry, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req = &userRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}

	checkUser, err := web.App.DB.GetUserByUsername(req.Username)
	if checkUser != nil {
		return nil, fmt.Errorf("user %s already exist", req.Username)
	}

	if err != nil {
		return nil, err
	}

	passwd, err := encryption.RSADecrypt(req.Password,
		*web.App.AppKey)
	if err == nil {
		if pwd, err := encryption.GetPasswordHASH([]byte(passwd)); err == nil {
			req.Password = string(pwd)
		}
	}

	return web.App.DB.CreateUser(req.Username, req.Password, req.FullName,
		req.IsAdmin, req.IsBlocked)
}

// Logout delete session.
func Logout(r *http.Request) error {
	session := r.Header.Get(middleware.SessionHeader)
	if session == "" || session == "null" {
		return nil
	}
	return web.App.DB.DeleteSessionByUUID(session)
}

// GetUserIsBlocked delete session.
func GetUserIsBlocked(r *http.Request) bool {
	user, err := GetUserBySession(r)
	if err != nil {
		return true
	}
	if user == nil {
		return false
	}
	return user.IsBlocked
}
