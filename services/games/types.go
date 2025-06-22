package games

import (
	"github.com/apprehensions/rbxweb/services/universes"
)

// CreatorID represents a Game creator on Roblox.
type CreatorID int64

// PlaceID represents a Place on Roblox.
type PlaceID int64

// AvatarType represents the type of a Roblox avatar.
type AvatarType string

const (
	MorphToR6    AvatarType = "MorphToR6"
	PlayerChoice AvatarType = "PlayerChoice"
	MorphToR15   AvatarType = "MorphToR15"
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
	ID                PlaceID              `json:"placeId"`
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	SourceName        string               `json:"sourceName"`
	SourceDescription string               `json:"sourceDescription"`
	URL               string               `json:"url"`
	Builder           string               `json:"builder"`
	BuilderID         CreatorID            `json:"builderId"`
	HasVerifiedBadge  bool                 `json:"hasVerifiedBadge"`
	IsPlayable        bool                 `json:"isPlayable"`
	ReasonProhibited  string               `json:"reasonProhibited"`
	UniverseID        universes.UniverseID `json:"universeId"`
	RootID            PlaceID              `json:"universeRootPlaceId"`
	Price             int64                `json:"price"`
	ImageToken        string               `json:"imageToken"`
}
