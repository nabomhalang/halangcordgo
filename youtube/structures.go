package youtube

import "google.golang.org/api/youtube/v3"

type Youtube struct {
	client *youtube.Service
}

type Video struct {
	ID          string
	Title       string
	Thumbnail   string
	Duration    float64
	Link        string
	Description string
	Channel     string
}

type YtDLP struct {
	Duration         float64          `json:"duration"`
	Thumbnail        string           `json:"thumbnail"`
	Extractor        string           `json:"extractor"`
	ID               string           `json:"id"`
	WebpageURL       string           `json:"webpage_url"`
	Title            string           `json:"title"`
	RequestedFormats RequestedFormats `json:"requested_formats"`
}

type RequestedFormats []struct {
	Resolution string `json:"resolution"`
}
