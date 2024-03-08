package rbxweb

import (
	"net/url"
	"testing"
)

func TestGetURL(t *testing.T) {
	query := url.Values{
		"universeIds": {"189707", "292439477"},
	}

	exp := "https://thumbnails.roblox.com/v1/games/icons"
	if url := GetURL("thumbnails", "v1/games/icons", nil); url != exp {
		t.Fatalf("want %s, got %s", exp, url)
	}

	exp += "?universeIds=189707&universeIds=292439477"
	if url := GetURL("thumbnails", "v1/games/icons", query); url != exp {
		t.Fatalf("want %s, got %s", exp, url)
	}
}
