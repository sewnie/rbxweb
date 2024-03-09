// games selectively implements the Roblox Games Web API.
package games

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/apprehensions/rbxweb"
)

type AvatarType string

const (
	MorphToR6    AvatarType = "MorphToR6"
	PlayerChoice AvatarType = "PlayerChoice"
	MorphToR15   AvatarType = "MorphToR15"
)

// Creator implements the GameCreator API model.
type Creator struct {
	ID               rbxweb.CreatorID `json:"id"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	IsRNVAccount     bool             `json:"isRNVAccount"`
	HasVerifiedBadge bool             `json:"hasVerifiedBadge"`
}

// GameDetail implements the GameDetailResponse API model.
type GameDetail struct {
	ID                        rbxweb.PlaceID `json:"id"`
	RootID                    rbxweb.PlaceID `json:"rootPlaceId"`
	Name                      string         `json:"name"`
	Description               string         `json:"description"`
	SourceName                string         `json:"sourceName"`
	SourceDescription         string         `json:"sourceDescription"`
	Creator                   Creator        `json:"creator"`
	Price                     int64          `json:"price"`
	AllowedGearGenres         []string       `json:"allowedGearGenres"`
	AllowedGearCategories     []string       `json:"allowedGearCategories"`
	IsGenreEnforced           bool           `json:"isGenreEnforced"`
	CopyingAllowed            bool           `json:"copyingAllowed"`
	Playing                   int64          `json:"playing"`
	Visits                    int64          `json:"visits"`
	MaxPlayers                int32          `json:"maxPlayers"`
	Created                   string         `json:"created"`
	Updated                   string         `json:"updated"`
	StudioAccessToApisAllowed bool           `json:"studioAccessToApisAllowed"`
	CreateVipServersAllowed   bool           `json:"createVipServersAllowed"`
	UniverseAvatarType        AvatarType     `json:"universeAvatarType"`
	Genre                     string         `json:"genre"`
	IsAllGenre                bool           `json:"isAllGenre"`
	IsFavoritedByUser         bool           `json:"isFavoritedByUser"`
	FavoritedCount            int64          `json:"favoritedCount"`
}

// PlaceDetail implements the PlaceDetails API model.
type PlaceDetail struct {
	ID                rbxweb.PlaceID    `json:"placeId"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	SourceName        string            `json:"sourceName"`
	SourceDescription string            `json:"sourceDescription"`
	URL               string            `json:"url"`
	Builder           string            `json:"builder"`
	BuilderID         rbxweb.CreatorID  `json:"builderId"`
	HasVerifiedBadge  bool              `json:"hasVerifiedBadge"`
	IsPlayable        bool              `json:"isPlayable"`
	ReasonProhibited  string            `json:"reasonProhibited"`
	UniverseID        rbxweb.UniverseID `json:"universeId"`
	RootID            rbxweb.PlaceID    `json:"universeRootPlaceId"`
	Price             int64             `json:"price"`
	ImageToken        string            `json:"imageToken"`
}

type gameDetailsResponse struct {
	Data []GameDetail `json:"data"`
}

// GetGameDetails returns the given Universe ID's game details.
func GetGameDetail(universeID rbxweb.UniverseID) (*GameDetail, error) {
	gds, err := GetGamesDetail([]rbxweb.UniverseID{universeID})
	if err != nil {
		return nil, err
	}

	if len(gds) == 0 {
		return nil, rbxweb.ErrNoData
	}
	return &gds[0], nil
}

// GetGamesDetail returns a list of the game details of each given Universe ID.
func GetGamesDetail(universeIDs []rbxweb.UniverseID) ([]GameDetail, error) {
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

// GetPlacesDetail returns a list of the place details of each given Place ID.
func GetPlacesDetail(placeIDs []rbxweb.PlaceID) ([]PlaceDetail, error) {
	var pds []PlaceDetail

	if len(placeIDs) == 0 {
		return nil, errors.New("placeIDs missing")
	}

	var uids []string
	for _, uid := range placeIDs {
		uids = append(uids, strconv.FormatInt(int64(uid), 10))
	}

	query := url.Values{"placeIds": uids}

	err := rbxweb.Request("GET", rbxweb.GetURL("games", "v1/games/multiget-place-details", query), nil, &pds)
	if err != nil {
		return nil, err
	}

	return pds, nil
}

// GetPlaceDetails returns the given Place ID's Place details.
func GetPlaceDetail(placeID rbxweb.PlaceID) (*PlaceDetail, error) {
	pds, err := GetPlacesDetail([]rbxweb.PlaceID{placeID})
	if err != nil {
		return nil, err
	}

	if len(pds) == 0 {
		return nil, rbxweb.ErrNoData
	}
	return &pds[0], nil
}
