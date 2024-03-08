// users selectively implements the Roblox Users Web API.
package users

import (
	"github.com/apprehensions/rbxweb"
)

type UserID int64

// AuthenticatedUser implements the AuthenticatedUserResponse API model.
type AuthenticatedUser struct {
	ID          UserID `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json":displayName"`
}

// GetAuthenticated returns the minimal authenticated user.
func GetAuthenticated() (*AuthenticatedUser, error) {
	var au AuthenticatedUser

	err := rbxweb.Request("GET", rbxweb.GetURL("users", "v1/users/authenticated", nil), nil, &au)
	if err != nil {
		return nil, err
	}

	return &au, nil
}

// UserIDRequest implements the MultiGetByUserIdRequest API model.
type UserIDRequest struct {
	IDs                []UserID `json:"userIds"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

// User implements the VerifiedBadgeUserResponse API model.
type User struct {
	Verified    bool   `json:"hasVerifiedBadge"`
	ID          UserID `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type userResponse struct {
	Data []User `json:"data"`
}

// GetUsers returns a list of users by their IDs.
func GetUsers(req UserIDRequest) ([]User, error) {
	var ur userResponse

	err := rbxweb.Request("GET", rbxweb.GetURL("users", "v1/users", nil), req, &ur)
	if err != nil {
		return nil, err
	}

	return ur.Data, nil
}

// GetUser returns a user by their ID.
func GetUser(req UserIDRequest) (*User, error) {
	us, err := GetUsers(req)
	if err != nil {
		return nil, err
	}

	if len(us) == 0 {
		return nil, rbxweb.ErrNoData
	}
	return &us[0], nil
}
