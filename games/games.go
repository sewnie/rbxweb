// games selectively implements the Roblox Games Web API.
package games

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/apprehensions/rbxweb"
)

type (
	UniverseID int64
	UserID     int64
)

type AvatarType string

const (
	MorphToR6    AvatarType = "MorphToR6"
	PlayerChoice AvatarType = "PlayerChoice"
	MorphToR15   AvatarType = "MorphToR15"
)

// Creator implements the GameCreator API model.
type Creator struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsRNVAccount     bool   `json:"isRNVAccount"`
	HasVerifiedBadge bool   `json:"hasVerifiedBadge"`
}

// GameDetail implements the GameDetailResponse API model.
type GameDetail struct {
	ID                        UniverseID `json:"id"`
	RootPlaceID               UniverseID `json:"rootPlaceId"`
	Name                      string     `json:"name"`
	Description               string     `json:"description"`
	SourceName                string     `json:"sourceName"`
	SourceDescription         string     `json:"sourceDescription"`
	Creator                   Creator    `json:"creator"`
	Price                     int64      `json:"price"`
	AllowedGearGenres         []string   `json:"allowedGearGenres"`
	AllowedGearCategories     []string   `json:"allowedGearCategories"`
	IsGenreEnforced           bool       `json:"isGenreEnforced"`
	CopyingAllowed            bool       `json:"copyingAllowed"`
	Playing                   int64      `json:"playing"`
	Visits                    int64      `json:"visits"`
	MaxPlayers                int32      `json:"maxPlayers"`
	Created                   string     `json:"created"`
	Updated                   string     `json:"updated"`
	StudioAccessToApisAllowed bool       `json:"studioAccessToApisAllowed"`
	CreateVipServersAllowed   bool       `json:"createVipServersAllowed"`
	UniverseAvatarType        AvatarType `json:"universeAvatarType"`
	Genre                     string     `json:"genre"`
	IsAllGenre                bool       `json:"isAllGenre"`
	IsFavoritedByUser         bool       `json:"isFavoritedByUser"`
	FavoritedCount            int64      `json:"favoritedCount"`
}

type gameDetailsResponse struct {
	Data []GameDetail `json:"data"`
}

// GetGameDetails returns the given Universe ID's game details.
func GetGameDetail(universeID UniverseID) (*GameDetail, error) {
	gds, err := GetGamesDetail([]UniverseID{universeID})
	if err != nil {
		return nil, err
	}

	if len(gds) == 0 {
		return nil, rbxweb.ErrNoData
	}
	return &gds[0], nil
}

// GetGamesDetail returns a list of the game details of each given Universe ID.
func GetGamesDetail(universeIDs []UniverseID) ([]GameDetail, error) {
	var gdr gameDetailsResponse

	if len(universeIDs) == 0 {
		return nil, errors.New("universeIDs missing")
	}

	var uids []string
	for _, uid := range universeIDs {
		uids = append(uids, strconv.FormatInt(int64(uid), 10))
	}

	query := url.Values{"universeIds": uids}

	err := rbxweb.Request("GET", rbxweb.GetURL("games", "v1/games", query), nil, &gdr)
	if err != nil {
		return nil, err
	}

	return gdr.Data, nil
}
