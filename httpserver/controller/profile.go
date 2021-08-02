package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/database/interfaces"
	"github.com/JIexa24/chef-webapi/encryption"
)

// ProfileUserInfo containы profile info.
type ProfileUserInfo struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
	Username string `json:"username,omitempty"`
	UserID   string `json:"userID,omitempty"`
}

// ProfileUserRequest contains request of profile.
type ProfileUserRequest struct {
	ProfileUserInfo
	Password string `json:"password"`
}

// ExtractProfileUserInfo extract all needed fields from user's entry.
func ExtractProfileUserInfo(
	user *interfaces.UserEntry) *ProfileUserInfo {
	if user == nil {
		return nil
	}
	return &ProfileUserInfo{
		ID:       user.ID,
		UserID:   "",
		Username: "",
		FullName: user.FullName,
		Avatar:   user.Avatar,
	}
}

// ImplementToProfileUserInfo implement all needed fields into user's entry.
func ImplementToProfileUserInfo(profile *ProfileUserRequest,
	user *interfaces.UserEntry) {
	if user == nil || profile == nil {
		return
	}
	user.FullName = profile.FullName
	user.Avatar = profile.Avatar
	if profile.Password != "" {
		passwd, err := encryption.RSADecrypt(profile.Password, *web.App.AppKey)
		if err == nil {
			var pwd []byte
			var err error
			if pwd, err = encryption.GetPasswordHASH([]byte(passwd)); err == nil {
				user.Password = string(pwd)
			} else {
				web.App.Logger.Errorf(
					"ImplementToProfileUserInfo [encryption.GetPasswordHASH]: %s",
					err.Error())
			}
		} else {
			web.App.Logger.Errorf(
				"ImplementToProfileUserInfo [encryption.RSADecrypt]: %s", err.Error())
		}
	}
}

// GetUserProfiles return user profile information.
// Information get only by id in url query.
func GetUserProfiles(params *APIParameters) ([]ProfileUserInfo, error) {
	profiles := make([]ProfileUserInfo, 0)
	for _, id := range params.ID {
		user, err := web.App.DB.GetUserByID(id)
		if err != nil {
			continue
		}
		if user != nil {
			profile := ExtractProfileUserInfo(user)
			profiles = append(profiles, *profile)
		}
	}
	return profiles, nil
}

// GetUserProfileByID returns information about the user's profile.
// Information is requested by id.
func GetUserProfileByID(r *http.Request) (*ProfileUserInfo, error) {
	user, err := GetUserByID(r)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	profile := ExtractProfileUserInfo(user)
	return profile, nil
}

// UpdateCurrentUserProfile updates user information.
// Information update by id.
func UpdateCurrentUserProfile(r *http.Request) error {
	user, err := GetUserBySession(r)
	if err != nil {
		web.App.Logger.Errorf(
			"profileGetCurrentUserAPIHandler [controller.GetUserBySession]: %s",
			err.Error())
		return err
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	profile := &ProfileUserRequest{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		return err
	}
	ImplementToProfileUserInfo(profile, user)
	return web.App.DB.UpdateUserByID(user.ID, user)
}
