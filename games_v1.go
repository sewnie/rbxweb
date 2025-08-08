package rbxweb

import (
	"net/url"
)

// GamesServiceV1 partially handles the 'games/v1' Roblox Web API.
type GamesServiceV1 service

// CreatorID represents a Game creator on Roblox.
type CreatorID int64

// PlaceID represents a Place on Roblox.
type PlaceID int64

// UniverseID represents a Universe on Roblox.
type UniverseID int64

// AvatarType represents the type of a Roblox avatar.
type AvatarType string

const (
	AvatarTypeR6     AvatarType = "MorphToR6"
	AvatarTypeChoice AvatarType = "PlayerChoice"
	AvatarTypeR15    AvatarType = "MorphToR15"
)

// Creator implements the GameCreator API model.
type Creator struct {
	ID               CreatorID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	IsRNVAccount     bool      `json:"isRNVAccount"`
	HasVerifiedBadge bool      `json:"hasVerifiedBadge"`
}

// GameDetail implements the GameDetailResponse API model.
type GameDetail struct {
	ID                        PlaceID    `json:"id"`
	RootID                    PlaceID    `json:"rootPlaceId"`
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

// PlaceDetail implements the PlaceDetails API model.
type PlaceDetail struct {
	ID                PlaceID    `json:"placeId"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	SourceName        string     `json:"sourceName"`
	SourceDescription string     `json:"sourceDescription"`
	URL               string     `json:"url"`
	Builder           string     `json:"builder"`
	BuilderID         CreatorID  `json:"builderId"`
	HasVerifiedBadge  bool       `json:"hasVerifiedBadge"`
	IsPlayable        bool       `json:"isPlayable"`
	ReasonProhibited  string     `json:"reasonProhibited"`
	UniverseID        UniverseID `json:"universeId"`
	RootID            PlaceID    `json:"universeRootPlaceId"`
	Price             int64      `json:"price"`
	ImageToken        string     `json:"imageToken"`
}

// GetGamesDetail returns a list of the game details of each given Universe ID.
func (g *GamesServiceV1) ListGamesDetails(uids []UniverseID) ([]GameDetail, error) {
	gdr := struct {
		Data []GameDetail `json:"data"`
	}{}

	query := url.Values{"universeIds": formatSlice(uids)}
	err := g.Client.Execute("GET", "games", path("v1/games", query), nil, &gdr)
	if err != nil {
		return nil, err
	}

	return gdr.Data, nil
}

// GetGameDetail returns the given Universe ID's game details.
//
// If none are found, nil will be returned.
func (g *GamesServiceV1) GetGameDetail(uid UniverseID) (*GameDetail, error) {
	return getList(g.ListGamesDetails([]UniverseID{uid}))
}

// ListPlacesDetail returns a list of the place details of each given Place ID.
func (g *GamesServiceV1) ListPlacesDetails(pids []PlaceID) ([]PlaceDetail, error) {
	var pds []PlaceDetail

	query := url.Values{"placeIds": formatSlice(pids)}
	err := g.Client.Execute("GET", "games", path("v1/games/multiget-place-details", query), nil, &pds)
	if err != nil {
		return nil, err
	}

	return pds, nil
}

// GetPlaceDetails returns the given Place ID's Place details.
//
// If none are found, nil will be returned.
func (g *GamesServiceV1) GetPlaceDetail(placeID PlaceID) (*PlaceDetail, error) {
	return getList(g.ListPlacesDetails([]PlaceID{placeID}))
}
