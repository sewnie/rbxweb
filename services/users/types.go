package users

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
