package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents a spotify web api client.
type Client interface {
	// GetDevices fetches the active devices
	// of the user.
	GetDevices() (*DeviceResponse, error)

	// GetPlayer fetches the player state
	// of the user.
	GetPlayer() (*Player, error)

	// Search searches on spotify with a given query.
	Search(q string, searchType string) (*SearchResponse, error)

	// PlayTracks plays a given track on
	// a given device.
	PlayTrack(uri string, deviceId string) error

	// GetLikedTracks fetches the liked tracks
	// of the user.
	GetLikedTracks(limit int, offset int) (*LikeTracksResponse, error)

	// AddQueue adds a given tracks to
	// the queue for a given device.
	AddQueue(uri string, deviceId string) error

	// Pause pauses the playback on a given device.
	Pause(deviceId string) error

	// Play resumes the playback on a given device.
	Play(deviceId string) error

	// Next changes to the next track on a given device.
	Next(deviceId string) error

	// Previous changes to the previous track on a given device.
	Previous(deviceId string) error

	// GetQueue fetches the player queue.
	GetQueue() (*QueueResponse, error)

	// GetRecentlyPlayedTracks fetches the recently
	// played tracks of the user.
	GetRecentlyPlayedTracks() (*RecentlyPlayedResponse, error)

	// GetSavedAlbums fetches the saved albums in
	// a library of the user.
	GetSavedAlbums(limit int, offset int) (*SavedAlbumResponse, error)

	// PlayContext plays a given context (playlist, album, etc.)
	// on a given device and can optionally handle a given uri in the context.
	PlayContext(contextUri string, deviceId string, uri ...string) error

	// SetShuffleState sets the shuffle state of the player for
	// a given device.
	SetShuffleState(deviceId string, state bool) error

	// SetRepeatMode sets the repeat mode of the player for
	// a given device.
	SetRepeatMode(deviceId string, state RepeatState) error

	// GetAlbum fetches a album by id.
	GetAlbum(id string) (*AlbumWithTracks, error)

	// GetPlaylist fetches a playlist by id.
	GetPlaylist(id string) (*Playlist, error)

	// GetPlaylistTracks fetches one page of tracks of a playlist.
	GetPlaylistTracks(id string, limit int, offset int) (*PlaylistItemsResponse, error)

	// AreTracksSaved checks for each given track id if it
	// is saved in the library of the user.
	AreTracksSaved(ids []string) ([]bool, error)

	// GetAllLikedTracks fetches all liked tracks of the user.
	GetAllLikedTracks() ([]Track, error)
}

type client struct {
	refreshToken string
	accessToken  string

	authClient AuthClient
	httpClient *http.Client
}

const (
	spotifyApiBaseUrl = "https://api.spotify.com/v1"
	httpTimeout       = 10 * time.Second

	RepeatTrack   RepeatState = "track"
	RepeatContext RepeatState = "context"
	RepeatOff     RepeatState = "off"
)

// URIToID parses the id from a given uri.
func URIToID(uri string) string {
	split := strings.Split(uri, ":")
	if len(split) != 3 {
		return ""
	}
	return split[2]
}

// NewClient creates a new spotify web
// api client.
func NewClient(
	refreshToken string,
	clientId string,
	clientSecret string,
) Client {
	return &client{
		refreshToken: refreshToken,
		authClient: NewAuthClient(
			clientId, clientSecret, "", []string{},
		),
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// doRequestWithToken is am internal implementation to
// execute the request with an access token.
// It requests a new one when no token exists.
func (c *client) doRequestWithToken(req *http.Request) (*http.Response, error) {
	if c.accessToken == "" {
		token, err := c.authClient.RequestRefreshedToken(c.refreshToken)
		if err != nil {
			return nil, err
		}
		c.accessToken = token
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	return c.httpClient.Do(req)
}

// doRequest is an internal implementation to execute
// a request with an access token and retry the request
// when the token is expired.
// It also returns an error for status code >= 400.
func (c *client) doRequest(req *http.Request) (*http.Response, error) {
	res, err := c.doRequestWithToken(req)

	if res.StatusCode == http.StatusUnauthorized {
		c.accessToken = ""
		return c.doRequestWithToken(req)
	} else if res.StatusCode >= 400 {
		fmt.Println(res.StatusCode)
		return nil, errors.New("failed")
	}

	return res, err
}

// getResult is an internal implementation to read
// and unmarshal the json http response to a struct.
func (c *client) getResult(res *http.Response, dest interface{}) error {
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, dest)
}

func (c *client) GetDevices() (*DeviceResponse, error) {
	url := fmt.Sprintf("%s/me/player/devices", spotifyApiBaseUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data DeviceResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) Search(q string, searchType string) (*SearchResponse, error) {
	params := url.Values{}
	params.Add("q", q)
	params.Add("type", searchType)
	// TODO: make limit adjustable
	params.Add("limit", "10")

	url := fmt.Sprintf("%s/search?%s", spotifyApiBaseUrl, params.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data SearchResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) GetPlayer() (*Player, error) {
	url := fmt.Sprintf("%s/me/player", spotifyApiBaseUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var data Player
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) PlayTrack(uri string, deviceId string) error {
	u := fmt.Sprintf("%s/me/player/play", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	reqBody := map[string]interface{}{
		"uris": []string{uri},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) GetLikedTracks(limit int, offset int) (*LikeTracksResponse, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	u := fmt.Sprintf("%s/me/tracks?%s", spotifyApiBaseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data LikeTracksResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) AddQueue(uri string, deviceId string) error {
	params := url.Values{}
	params.Add("uri", uri)

	if deviceId != "" {
		params.Add("device_id", deviceId)
	}

	u := fmt.Sprintf("%s/me/player/queue?%s", spotifyApiBaseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) Pause(deviceId string) error {
	u := fmt.Sprintf("%s/me/player/pause", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := http.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) Play(deviceId string) error {
	u := fmt.Sprintf("%s/me/player/play", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := http.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) Next(deviceId string) error {
	u := fmt.Sprintf("%s/me/player/next", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := http.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) Previous(deviceId string) error {
	u := fmt.Sprintf("%s/me/player/previous", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := http.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) GetQueue() (*QueueResponse, error) {
	u := fmt.Sprintf("%s/me/player/queue", spotifyApiBaseUrl)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data QueueResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) GetRecentlyPlayedTracks() (*RecentlyPlayedResponse, error) {
	u := fmt.Sprintf("%s/me/player/recently-played", spotifyApiBaseUrl)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data RecentlyPlayedResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) GetSavedAlbums(limit int, offset int) (*SavedAlbumResponse, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	u := fmt.Sprintf("%s/me/albums?%s", spotifyApiBaseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data SavedAlbumResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) PlayContext(contextUri string, deviceId string, uri ...string) error {
	u := fmt.Sprintf("%s/me/player/play", spotifyApiBaseUrl)

	if deviceId != "" {
		params := url.Values{}
		params.Add("device_id", deviceId)
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	reqBody := map[string]interface{}{
		"context_uri": contextUri,
	}

	if len(uri) > 0 {
		reqBody["offset"] = map[string]interface{}{
			"uri": uri[0],
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) SetShuffleState(deviceId string, state bool) error {
	params := url.Values{}
	params.Add("state", strconv.FormatBool(state))

	if deviceId != "" {
		params.Add("device_id", deviceId)
	}

	u := fmt.Sprintf("%s/me/player/shuffle?%s", spotifyApiBaseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) SetRepeatMode(deviceId string, state RepeatState) error {
	params := url.Values{}
	params.Add("state", string(state))

	if deviceId != "" {
		params.Add("device_id", deviceId)
	}

	u := fmt.Sprintf("%s/me/player/repeat?%s", spotifyApiBaseUrl, params.Encode())

	req, err := http.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *client) GetAlbum(id string) (*AlbumWithTracks, error) {
	u := fmt.Sprintf("%s/albums/%s", spotifyApiBaseUrl, id)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data AlbumWithTracks
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) GetPlaylist(id string) (*Playlist, error) {
	params := url.Values{}
	params.Add("fields", "id,name,uri")

	u := fmt.Sprintf("%s/playlists/%s?%s", spotifyApiBaseUrl, id, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data Playlist
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) GetPlaylistTracks(id string, limit int, offset int) (*PlaylistItemsResponse, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))

	u := fmt.Sprintf("%s/playlists/%s/items?%s", spotifyApiBaseUrl, id, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var data PlaylistItemsResponse
	if err := c.getResult(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *client) AreTracksSaved(ids []string) ([]bool, error) {
	const chunkSize = 50

	saved := make([]bool, 0, len(ids))
	for start := 0; start < len(ids); start += chunkSize {
		end := start + chunkSize
		if end > len(ids) {
			end = len(ids)
		}

		params := url.Values{}
		params.Add("ids", strings.Join(ids[start:end], ","))

		u := fmt.Sprintf("%s/me/tracks/contains?%s", spotifyApiBaseUrl, params.Encode())

		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}

		res, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		var chunk []bool
		if err := c.getResult(res, &chunk); err != nil {
			return nil, err
		}

		saved = append(saved, chunk...)
	}

	return saved, nil
}

func (c *client) GetAllLikedTracks() ([]Track, error) {
	const pageSize = 50

	tracks := make([]Track, 0, pageSize)
	offset := 0
	for {
		res, err := c.GetLikedTracks(pageSize, offset)
		if err != nil {
			return nil, err
		}

		for _, item := range res.Items {
			tracks = append(tracks, item.Track)
		}

		offset += len(res.Items)
		if len(res.Items) < pageSize || offset >= res.Total {
			break
		}
	}

	return tracks, nil
}
