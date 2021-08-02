package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	web "github.com/JIexa24/chef-webapi"

	"github.com/JIexa24/chef-webapi/database/interfaces"
	"github.com/JIexa24/chef-webapi/encryption"
	"github.com/JIexa24/chef-webapi/httpserver/middleware"
)

// AuthenticationSessionInfo contains session information.
type AuthenticationSessionInfo struct {
	Session string `json:"session"`
}

// AuthenticationUserInfo contains user information.
type AuthenticationUserInfo struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	FullName           string `json:"fullName"`
	Avatar             string `json:"avatar"`
	LastLogin          string `json:"lastLogin"`
	NeedPasswordChange bool   `json:"needPasswordChange"`
}

// ExtractAuthenticationUserInfo extracts all needed fields
// from the user's record..
func ExtractAuthenticationUserInfo(
	user *interfaces.UserEntry) *AuthenticationUserInfo {
	if user == nil {
		return nil
	}
	return &AuthenticationUserInfo{
		ID:                 user.ID,
		Username:           user.Username,
		FullName:           user.FullName,
		Avatar:             user.Avatar,
		LastLogin:          user.LastLogin,
		NeedPasswordChange: user.NeedPasswordChange,
	}
}

// AuthenticationUserPermissions contains admin flag.
type AuthenticationUserPermissions struct {
	IsAdmin bool `json:"admin"`
}

// ExtractAuthenticationUserPermissions extract permissions fields of user.
func ExtractAuthenticationUserPermissions(
	user *interfaces.UserEntry) *AuthenticationUserPermissions {
	if user == nil {
		return nil
	}
	return &AuthenticationUserPermissions{
		IsAdmin: user.IsAdmin,
	}
}

// AuthenticationPing contains is answer of question of user authenticate.
type AuthenticationPing struct {
	Authenticate       bool `json:"authenticate"`
	NeedPasswordChange bool `json:"needPasswordChange"`
}

type authenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticate is a function to authenticate a user.
func Authenticate(r *http.Request) (*AuthenticationSessionInfo, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req = &authenticationRequest{}

	err = json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}

	user, err := web.App.DB.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	passwd, err := encryption.RSADecrypt(req.Password, *web.App.AppKey)
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsBlocked ||
		encryption.CheckPasswordHASH([]byte(passwd), []byte(user.Password)) != nil {
		return nil, fmt.Errorf("bad creditinals for user %v", req.Username)
	}

	expire := time.Now().Unix() + web.App.SessionExpire
	s, err := web.App.DB.CreateSession(user.Username, expire)
	if err != nil {
		return nil, err
	}
	session := &AuthenticationSessionInfo{
		Session: s.UUID,
	}

	user.LastLogin = time.Unix(time.Now().Unix(), 0).Format(interfaces.TimeFormat)
	err = web.App.DB.UpdateUserByID(user.ID, user)
	return session, err
}

// CheckSession will check the existence of the session.
func CheckSession(r *http.Request) (*AuthenticationPing, error) {
	session := r.Header.Get(middleware.SessionHeader)
	if session == "" || session == "null" {
		return &AuthenticationPing{Authenticate: false}, nil
	}
	s, err := web.App.DB.GetSessionByUUID(session)
	nowTime := time.Now().Unix()
	if (s == nil || err != nil) || (s != nil && s.Expire < nowTime) {
		return &AuthenticationPing{Authenticate: false}, err
	}

	user, err := GetUserBySession(r)
	if err == nil && user != nil {
		web.App.UsersLastSeen[user.Username] =
			time.Unix(nowTime, 0).Format(interfaces.TimeFormat)
	}

	if s.Expire-int64(float64(web.App.SessionExpire)*0.2) < nowTime {
		s.Expire = nowTime + web.App.SessionExpire
		err = web.App.DB.UpdateSessionByUUID(session, s)
	}
	return &AuthenticationPing{Authenticate: true,
		NeedPasswordChange: user.NeedPasswordChange}, err
}
