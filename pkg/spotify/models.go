package spotify

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
	Device       Device      `json:"device"`
	ShuffleState bool        `json:"shuffle_state"`
	RepeatState  RepeatState `json:"repeat_state"`
	ProgressMs   int         `json:"progress_ms"`
	Item         Track       `json:"item"`
	IsPlaying    bool        `json:"is_playing"`
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
