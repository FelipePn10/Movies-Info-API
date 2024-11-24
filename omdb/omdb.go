package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Results struct {
	Search       []Search `json:"Search"`
	TotalResults string   `json:"totalResults"`
	Response     string   `json:"Response"`
}

type Search struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

func SearchMovies(apiKey, title string) (Results, error) {
	url := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&s=%s", apiKey, title)

	resp, err := http.Get(url)
	if err != nil {
		return Results{}, fmt.Errorf("failed to make request to OMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Results{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var results Results
	if err = json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Results{}, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if results.Response != "True" {
		return Results{}, fmt.Errorf("OMDB API error: no results or invalid response")
	}

	return results, nil
}
