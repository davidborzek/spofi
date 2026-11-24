package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// fakeTransport serves canned responses per path so the api
// client can be tested without network access.
type fakeTransport struct {
	handler func(req *http.Request) (string, error)
}

func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := f.handler(req)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

func newTestClient(handler func(req *http.Request) (string, error)) *client {
	return &client{
		accessToken: "test",
		httpClient: &http.Client{
			Transport: &fakeTransport{handler: handler},
		},
	}
}

func TestGetAllLikedTracksPagination(t *testing.T) {
	const total = 120

	c := newTestClient(func(req *http.Request) (string, error) {
		q := req.URL.Query()
		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		if limit != 50 {
			t.Errorf("expected limit 50, got %d", limit)
		}

		items := make([]map[string]interface{}, 0, limit)
		for i := offset; i < offset+limit && i < total; i++ {
			items = append(items, map[string]interface{}{
				"track": map[string]interface{}{
					"id":   fmt.Sprintf("id%d", i),
					"uri":  fmt.Sprintf("spotify:track:id%d", i),
					"name": fmt.Sprintf("track %d", i),
				},
			})
		}

		body := map[string]interface{}{
			"items":  items,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		}

		data, err := json.Marshal(body)
		return string(data), err
	})

	tracks, err := c.GetAllLikedTracks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tracks) != total {
		t.Fatalf("expected %d tracks, got %d", total, len(tracks))
	}

	if tracks[len(tracks)-1].ID != fmt.Sprintf("id%d", total-1) {
		t.Errorf("unexpected last track: %+v", tracks[len(tracks)-1])
	}
}

func TestGetAllLikedTracksExactMultiple(t *testing.T) {
	const total = 100

	c := newTestClient(func(req *http.Request) (string, error) {
		q := req.URL.Query()
		offset, _ := strconv.Atoi(q.Get("offset"))
		limit, _ := strconv.Atoi(q.Get("limit"))

		items := make([]map[string]interface{}, 0, limit)
		for i := offset; i < offset+limit && i < total; i++ {
			items = append(items, map[string]interface{}{
				"track": map[string]interface{}{
					"id":   fmt.Sprintf("id%d", i),
					"uri":  fmt.Sprintf("spotify:track:id%d", i),
					"name": fmt.Sprintf("track %d", i),
				},
			})
		}

		data, err := json.Marshal(map[string]interface{}{
			"items": items,
			"total": total,
		})
		return string(data), err
	})

	tracks, err := c.GetAllLikedTracks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tracks) != total {
		t.Fatalf("expected %d tracks, got %d", total, len(tracks))
	}
}

func TestAreTracksSavedChunks(t *testing.T) {
	ids := make([]string, 0, 120)
	for i := 0; i < 120; i++ {
		ids = append(ids, fmt.Sprintf("id%d", i))
	}

	calls := 0
	c := newTestClient(func(req *http.Request) (string, error) {
		calls++

		idsParam := req.URL.Query().Get("ids")
		if idsParam == "" {
			t.Error("expected ids parameter")
		}

		chunk := strings.Split(idsParam, ",")
		if len(chunk) > 50 {
			t.Errorf("chunk too large: %d", len(chunk))
		}

		result := make([]bool, len(chunk))
		for i := range chunk {
			result[i] = i%2 == 0
		}

		data, err := json.Marshal(result)
		return string(data), err
	})

	saved, err := c.AreTracksSaved(ids)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calls != 3 {
		t.Errorf("expected 3 requests, got %d", calls)
	}

	if len(saved) != len(ids) {
		t.Fatalf("expected %d results, got %d", len(ids), len(saved))
	}

	if !saved[0] || saved[1] {
		t.Errorf("unexpected results: %v", saved[:2])
	}
}

func TestPlaylistItemsTracks(t *testing.T) {
	const body = `{
		"items": [
			{"item": {"uri": "spotify:track:1", "name": "one"}},
			{"item": {"uri": "spotify:episode:2", "name": "podcast"}},
			{"item": null},
			{"track": {"uri": "spotify:track:3", "name": "three"}}
		],
		"total": 4
	}`

	var result PlaylistItemsResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tracks := result.Tracks()
	if len(tracks) != 2 {
		t.Fatalf("expected 2 tracks, got %d: %+v", len(tracks), tracks)
	}

	if tracks[0].URI != "spotify:track:1" || tracks[1].URI != "spotify:track:3" {
		t.Errorf("unexpected tracks: %+v", tracks)
	}
}

func TestGetPlaylistTracksQuery(t *testing.T) {
	var gotQuery url.Values
	c := newTestClient(func(req *http.Request) (string, error) {
		gotQuery = req.URL.Query()
		if !strings.HasSuffix(req.URL.Path, "/playlists/abc/items") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		return `{"items": [], "total": 0}`, nil
	})

	if _, err := c.GetPlaylistTracks("abc", 50, 100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotQuery.Get("limit") != "50" || gotQuery.Get("offset") != "100" {
		t.Errorf("unexpected query: %v", gotQuery)
	}
}
