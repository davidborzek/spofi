package spotify

import "strings"

type RepeatState string
type DeviceResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	VolumePercent int    `json:"volume_percent"`
}

type PagingResult struct {
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

type SearchResponse struct {
	Tracks SearchTrackResult `json:"tracks"`
	Albums SearchAlbumResult `json:"albums"`
}

type SearchTrackResult struct {
	Items []Track `json:"items"`
	PagingResult
}

type SearchAlbumResult struct {
	Items []Album `json:"items"`
	PagingResult
}

type Artist struct {
	ID   string `json:"id"`
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type Album struct {
	Artists     []Artist `json:"artists"`
	ID          string   `json:"id"`
	URI         string   `json:"uri"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"release_date"`
	TotalTracks int      `json:"total_tracks"`
}

type Track struct {
	Album       Album    `json:"album"`
	Artists     []Artist `json:"artists"`
	DiscNumber  int      `json:"disc_number"`
	DurationMs  int      `json:"duration_ms"`
	ID          string   `json:"id"`
	URI         string   `json:"uri"`
	Name        string   `json:"name"`
	TrackNumber int      `json:"track_number"`
}

type Player struct {
	Device       Device         `json:"device"`
	ShuffleState bool           `json:"shuffle_state"`
	RepeatState  RepeatState    `json:"repeat_state"`
	ProgressMs   int            `json:"progress_ms"`
	Item         Track          `json:"item"`
	IsPlaying    bool           `json:"is_playing"`
	Context      *PlayerContext `json:"context"`
}

// PlayerContext represents the context (playlist, album, artist, ...)
// the current playback was started from.
type PlayerContext struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// Playlist represents a spotify playlist.
type Playlist struct {
	ID   string `json:"id"`
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// PlaylistItemsResponse represents a page of playlist items.
type PlaylistItemsResponse struct {
	Items []struct {
		// Item is the current field, Track is the deprecated alias
		// kept for compatibility with the /tracks endpoint.
		Item  *Track `json:"item"`
		Track *Track `json:"track"`
	} `json:"items"`
	PagingResult
}

// Tracks returns the playable tracks of the page,
// skipping episodes, local files and empty entries.
func (r *PlaylistItemsResponse) Tracks() []Track {
	tracks := make([]Track, 0, len(r.Items))
	for _, item := range r.Items {
		t := item.Item
		if t == nil {
			t = item.Track
		}

		if t == nil || !strings.HasPrefix(t.URI, "spotify:track:") {
			continue
		}

		tracks = append(tracks, *t)
	}

	return tracks
}

type LikeTracksResponse struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
	PagingResult
}

type QueueResponse struct {
	Queue []Track `json:"queue"`
}

type RecentlyPlayedResponse struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

type AlbumWithTracks struct {
	Album
	Tracks struct {
		Items []Track `json:"items"`
		PagingResult
	} `json:"tracks"`
}

type SavedAlbumResponse struct {
	Items []struct {
		Album AlbumWithTracks `json:"album"`
	} `json:"items"`
	PagingResult
}
