package views

import (
	"testing"

	"github.com/davidborzek/spofi/internal/format"

	"github.com/davidborzek/spofi/pkg/spotify"
)

func track(name string, artist string, id string) spotify.Track {
	return spotify.Track{
		ID:   id,
		URI:  "spotify:track:" + id,
		Name: name,
		Artists: []spotify.Artist{
			{Name: artist, URI: "spotify:artist:" + artist},
		},
		Album: spotify.Album{
			Name: "Some Album",
			URI:  "spotify:album:" + id,
		},
	}
}

func TestBuildSearchRowsOrder(t *testing.T) {
	results := searchResults{
		tracks: []spotify.Track{
			track("Song", "Artist", "1"),
			track("Other", "Artist", "2"),
		},
		context:      []spotify.Track{track("Song", "Artist", "9")},
		contextURI:   "spotify:playlist:abc",
		contextLabel: "playlist: Chill",
		favorites:    []spotify.Track{track("Song", "Artist", "1")},
	}

	rows := buildSearchRows(results, "")

	want := []string{
		"Song | Artist  (playlist: Chill)",
		"Song | Artist  (liked songs)",
		"Song | Artist  (single track)",
		"Song | Artist  (album: Some Album)",
		"Other | Artist  (single track)",
	}

	if len(rows) != len(want) {
		t.Fatalf("expected %d rows, got %d: %+v", len(want), len(rows), rows)
	}

	for i, title := range want {
		if rows[i].Title != format.FormatIcon("", title) {
			t.Errorf("row %d: expected %q, got %q", i, title, rows[i].Title)
		}
	}
}

func TestBuildSearchRowsDeduplicatesContext(t *testing.T) {
	results := searchResults{
		context: []spotify.Track{
			track("Song", "Artist", "1"),
			track("Song", "Artist", "1"),
		},
		contextURI:   "spotify:album:x",
		contextLabel: "album: X",
	}

	rows := buildSearchRows(results, "")

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d: %+v", len(rows), rows)
	}
}

func TestBuildSearchRowsSkipsEmptyResults(t *testing.T) {
	rows := buildSearchRows(searchResults{}, "")

	if len(rows) != 0 {
		t.Fatalf("expected no rows, got %+v", rows)
	}
}

func TestDecodeActionResult(t *testing.T) {
	tests := []struct {
		value string
		want  searchAction
	}{
		{
			value: encodeActionResult(searchAction{
				Kind:    rowActionContext,
				Context: "spotify:playlist:abc",
				Track:   "spotify:track:1",
			}),
			want: searchAction{
				Kind:    rowActionContext,
				Context: "spotify:playlist:abc",
				Track:   "spotify:track:1",
			},
		},
		{
			value: "spotify:track:2",
			want: searchAction{
				Kind:  rowActionTrack,
				Track: "spotify:track:2",
			},
		},
	}

	for _, test := range tests {
		got := decodeActionResult(test.value)
		if got != test.want {
			t.Errorf("decodeActionResult(%q) = %+v, want %+v", test.value, got, test.want)
		}
	}
}
