// clientsettings selectively implements the 'users' Roblox Web API.
package users

import (
	"github.com/sewnie/rbxweb/internal/api"
)

// UsersServiceV1 partially handles the 'users/v1' Roblox Web API.
type UsersServiceV1 api.Service

// GetAuthenticated returns the minimal authenticated user.
func (u *UsersServiceV1) GetAuthenticated() (*AuthenticatedUser, error) {
	var au AuthenticatedUser

	err := u.Client.Execute("GET", "users", "v1/users/authenticated", nil, &au)
	if err != nil {
		return nil, err
	}

	return &au, nil
}

// GetUsers returns a list of users by their IDs.
func (u *UsersServiceV1) ListUsers(uid UserIDRequest) ([]User, error) {
	ur := struct {
		Data []User `json:"data"`
	}{}

	err := u.Client.Execute("GET", "users", "v1/users", uid, &ur)
	if err != nil {
		return nil, err
	}

	return ur.Data, nil
}

// GetUser returns a user by their ID.
//
// If none are found, nil will be returned.
func (u *UsersServiceV1) GetUser(uid UserIDRequest) (*User, error) {
	return api.GetList(u.ListUsers(uid))
}
