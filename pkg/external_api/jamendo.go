package external_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type JamendoClient struct {
	ClientID string
	BaseURL  string
}

func NewJamendoClient(clientID string) *JamendoClient {
	return &JamendoClient{
		ClientID: clientID,
		BaseURL:  "https://api.jamendo.com/v3.0",
	}
}

type JamendoTrack struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ArtistName   string `json:"artist_name"`
	AlbumName    string `json:"album_name"`
	Duration     int    `json:"duration"`
	Image        string `json:"image"`
	Audio        string `json:"audio"`
	AudioDownload string `json:"audiodownload"`
}

type JamendoResponse struct {
	Headers struct {
		Status  string `json:"status"`
		Code    int    `json:"code"`
		ResultsCount int `json:"results_count"`
	} `json:"headers"`
	Results []JamendoTrack `json:"results"`
}

type JamendoFeedResponse struct {
	Headers struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	} `json:"headers"`
	Results []struct {
		Type   string       `json:"type"`
		JoinID string       `json:"joinid"`
		Track  JamendoTrack `json:"track,omitempty"`
	} `json:"results"`
}

func (c *JamendoClient) SearchTracks(query string, limit int) ([]JamendoTrack, error) {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("format", "json")
	params.Add("search", query)
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("include", "musicinfo")

	resp, err := http.Get(fmt.Sprintf("%s/tracks/?%s", c.BaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (c *JamendoClient) GetTrackByID(id string) (*JamendoTrack, error) {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("format", "json")
	params.Add("id", id)

	resp, err := http.Get(fmt.Sprintf("%s/tracks/?%s", c.BaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("track not found")
	}

	return &result.Results[0], nil
}

func (c *JamendoClient) GetFeed(limit int) ([]JamendoTrack, error) {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("format", "json")
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("type", "track") // Get track-based feed

	resp, err := http.Get(fmt.Sprintf("%s/feeds/?%s", c.BaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JamendoFeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var tracks []JamendoTrack
	for _, res := range result.Results {
		if res.Type == "track" && res.Track.ID != "" {
			tracks = append(tracks, res.Track)
		}
	}

	return tracks, nil
}

func (c *JamendoClient) GetPopularTracks(limit int) ([]JamendoTrack, error) {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("format", "json")
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("order", "popularity_total")

	resp, err := http.Get(fmt.Sprintf("%s/tracks/?%s", c.BaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (c *JamendoClient) GetMostPlayedTracks(limit int) ([]JamendoTrack, error) {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("format", "json")
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("order", "listens_total")

	resp, err := http.Get(fmt.Sprintf("%s/tracks/?%s", c.BaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JamendoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}
