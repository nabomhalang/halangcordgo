package youtube

import "google.golang.org/api/youtube/v3"

type Youtube struct {
	client *youtube.Service
}

type Video struct {
	ID        string
	Title     string
	Thumbnail string
	Duration  float64
}
