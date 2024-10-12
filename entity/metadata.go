package entity

import "errors"

type Metadata struct {
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HdURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

var ErrNotFound = errors.New("not found")
