package rbxweb

import (
	"strconv"
)

type (
	UniverseID int64
	CreatorID  int64
	PlaceID    int64
)

// GetPlaceUniverse returns the given Place's Universe ID.
//
// Uses the undocumented API service universes/v1
func GetPlaceUniverse(placeID PlaceID) (UniverseID, error) {
	r := struct {
		uid UniverseID `json:"universeId"`
	}{}

	pid := strconv.FormatInt(int64(placeID), 10)
	err := Request("GET",
		GetURL("apis", "v1/universes/places/"+pid+"/universe", nil), nil, &r)
	if err != nil {
		return 0, err
	}

	return r.uid, nil
}
