package rbxweb

// UsersServiceV1 partially handles the 'users/v1' Roblox Web API.
type UsersServiceV1 service

// UserID represents a user on Roblox.
type UserID int64

// AuthenticatedUser implements the AuthenticatedUserResponse API model.
type AuthenticatedUser struct {
	ID          UserID `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// User implements the VerifiedBadgeUserResponse API model.
type User struct {
	Verified    bool   `json:"hasVerifiedBadge"`
	ID          UserID `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// UserIDRequest implements the MultiGetByUserIdRequest API model.
type UserIDRequest struct {
	IDs                []UserID `json:"userIds"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

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
	return getList(u.ListUsers(uid))
}
