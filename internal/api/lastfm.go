package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	lastFMBaseURL   = "https://ws.audioscrobbler.com/2.0/"
	lastFMAPIKey    = "6f1733ce690e7e654e47f3061df509d3"
	lastFMAPISecret = "a055effdfebb4fd83b66cdd830497b67"
)

type Client struct {
	apiKey    string
	apiSecret string
	client    *http.Client
}

type AuthResponse struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	Subscriber int    `json:"subscriber"`
}

type TrackSearchResult struct {
	Results struct {
		TrackMatches struct {
			Track []TrackSearchInfo `json:"track"`
		} `json:"trackmatches"`
	} `json:"results"`
}

type TrackSearchInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
	Image  []struct {
		Text string `json:"#text"`
		Size string `json:"size"`
	} `json:"image"`
	Mbid string `json:"mbid"`
}

type AlbumSearchResult struct {
	Results struct {
		AlbumMatches struct {
			Album []AlbumSearchInfo `json:"album"`
		} `json:"albummatches"`
	} `json:"results"`
}

type TrackInfo struct {
	Name     string     `json:"name"`
	Artist   ArtistInfo `json:"artist"`
	Album    AlbumInfo  `json:"album"`
	URL      string     `json:"url"`
	Date     DateInfo   `json:"date"`
	Duration string     `json:"duration"`
	Image    []struct {
		Text string `json:"#text"`
		Size string `json:"size"`
	} `json:"image"`
	Mbid string `json:"mbid"`
}

type ArtistInfo struct {
	Text string `json:"#text"`
	Mbid string `json:"mbid"`
}

type AlbumInfo struct {
	Text string `json:"#text"`
	Mbid string `json:"mbid"`
}

type AlbumSearchInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
	Image  []struct {
		Text string `json:"#text"`
		Size string `json:"size"`
	} `json:"image"`
	Mbid string `json:"mbid"`
}

type DateInfo struct {
	Text string `json:"#text"`
	UTS  string `json:"uts"`
}

type AlbumDetailResponse struct {
	Album struct {
		Name   string `json:"name"`
		Artist string `json:"artist"`
		Tracks struct {
			Track []struct {
				Name     string `json:"name"`
				Duration int    `json:"duration"`
				Artist   struct {
					Name string `json:"name"`
				} `json:"artist"`
			} `json:"track"`
		} `json:"tracks"`
	} `json:"album"`
}

type RecentTracksResponse struct {
	RecentTracks struct {
		Track []TrackInfo `json:"track"`
	} `json:"recenttracks"`
}

type ErrorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

type ScrobbleResponse struct {
	Scrobbles struct {
		Scrobble struct {
			Track struct {
				Name   string `json:"name"`
				Artist struct {
					Text string `json:"#text"`
				} `json:"artist"`
			} `json:"track"`
			Timestamp string `json:"timestamp"`
		} `json:"scrobble"`
	} `json:"scrobbles"`
}

func NewClient() *Client {
	return &Client{
		apiKey:    lastFMAPIKey,
		apiSecret: lastFMAPISecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetAPIKey() string {
	return c.apiKey
}

func (c *Client) GetAPISecret() string {
	return c.apiSecret
}

func (c *Client) GetSessionKey(username, password string) (string, error) {
	signature := c.createSignature(map[string]string{
		"api_key":  c.GetAPIKey(),
		"method":   "auth.getMobileSession",
		"password": password,
		"username": username,
	})

	params := url.Values{
		"api_key":  {c.GetAPIKey()},
		"method":   {"auth.getMobileSession"},
		"password": {password},
		"username": {username},
		"api_sig":  {signature},
		"format":   {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return "", err
	}

	var authResp struct {
		Session struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"session"`
	}
	if err := json.Unmarshal([]byte(data), &authResp); err != nil {
		var errResp ErrorResponse
		if err := json.Unmarshal([]byte(data), &errResp); err == nil {
			return "", fmt.Errorf("last.fm error %d: %s", errResp.Error, errResp.Message)
		}
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}

	if authResp.Session.Key == "" {
		return "", fmt.Errorf("no session key received from last.fm")
	}

	return authResp.Session.Key, nil
}

func (c *Client) SearchTrack(artist, track string) ([]TrackSearchInfo, error) {
	params := url.Values{
		"method":  {"track.search"},
		"api_key": {c.GetAPIKey()},
		"artist":  {artist},
		"track":   {track},
		"format":  {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result TrackSearchResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return result.Results.TrackMatches.Track, nil
}

func (c *Client) SearchAlbum(artist, album string) ([]AlbumSearchInfo, error) {
	params := url.Values{
		"method":  {"album.search"},
		"api_key": {c.GetAPIKey()},
		"artist":  {artist},
		"album":   {album},
		"format":  {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result AlbumSearchResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return result.Results.AlbumMatches.Album, nil
}

func (c *Client) GetAlbumInfo(artist, album string) (*AlbumDetailResponse, error) {
	params := url.Values{
		"method":      {"album.getInfo"},
		"api_key":     {c.GetAPIKey()},
		"artist":      {artist},
		"album":       {album},
		"autocorrect": {"1"},
		"format":      {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result AlbumDetailResponse
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ValidateAlbumForTrack(artist, track, album string) (string, error) {
	albumInfo, err := c.GetAlbumInfo(artist, album)
	if err != nil {
		return "", err
	}

	for _, t := range albumInfo.Album.Tracks.Track {
		if t.Name == track {
			return albumInfo.Album.Name, nil
		}
	}

	return "", nil
}

func (c *Client) ScrobbleTrack(artist, track, timestamp, sessionKey string, album string) error {
	signatureParams := map[string]string{
		"api_key":   c.GetAPIKey(),
		"artist":    artist,
		"method":    "track.scrobble",
		"sk":        sessionKey,
		"timestamp": timestamp,
		"track":     track,
	}
	if album != "" {
		signatureParams["album"] = album
	}

	signature := c.createSignature(signatureParams)

	params := url.Values{
		"api_key":   {c.GetAPIKey()},
		"artist":    {artist},
		"method":    {"track.scrobble"},
		"sk":        {sessionKey},
		"timestamp": {timestamp},
		"track":     {track},
		"api_sig":   {signature},
		"format":    {"json"},
	}
	if album != "" {
		params.Add("album", album)
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return err
	}

	var errResp ErrorResponse
	if err := json.Unmarshal([]byte(data), &errResp); err == nil && errResp.Error != 0 {
		return fmt.Errorf("last.fm error %d: %s", errResp.Error, errResp.Message)
	}

	return nil
}

func (c *Client) GetRecentTracks(username string, limit int) ([]TrackInfo, error) {
	params := url.Values{
		"method":  {"user.getrecenttracks"},
		"api_key": {c.GetAPIKey()},
		"user":    {username},
		"limit":   {fmt.Sprintf("%d", limit)},
		"format":  {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result RecentTracksResponse
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return result.RecentTracks.Track, nil
}

type TrackMetadata struct {
	Duration int
	Artist   string
	Track    string
	Album    string
}

func (c *Client) GetTrackInfo(artist, track string) (*TrackMetadata, error) {
	params := url.Values{
		"method":      {"track.getInfo"},
		"api_key":     {c.GetAPIKey()},
		"artist":      {artist},
		"track":       {track},
		"autocorrect": {"1"},
		"format":      {"json"},
	}

	data, _, err := c.doRequest(params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Track struct {
			Duration string `json:"duration"`
			Name     string `json:"name"`
			Artist   struct {
				Name string `json:"name"`
			} `json:"artist"`
			Album struct {
				Title string `json:"title"`
			} `json:"album"`
		} `json:"track"`
	}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	duration := 0
	if result.Track.Duration != "" {
		fmt.Sscanf(result.Track.Duration, "%d", &duration)
	}

	return &TrackMetadata{
		Duration: duration,
		Artist:   result.Track.Artist.Name,
		Track:    result.Track.Name,
		Album:    result.Track.Album.Title,
	}, nil
}

func (c *Client) createSignature(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sig string
	for _, k := range keys {
		sig += k + params[k]
	}
	sig += c.GetAPISecret()

	hash := md5.Sum([]byte(sig)) //nolint:gosec // Last.fm API requires MD5 signature
	return hex.EncodeToString(hash[:])
}

func (c *Client) doRequest(params url.Values) (string, int, error) {
	req, err := http.NewRequest("POST", lastFMBaseURL, strings.NewReader(params.Encode()))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", statusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return string(bodyBytes), statusCode, nil
}
