package youtube

import (
	"context"
	"strings"
	"time"

	"github.com/nabomhalang/halangcordgo/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	log *config.Logger = config.NewLogger("youtube")
)

func NewYoutube(key string) (*Youtube, error) {
	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	return &Youtube{
		client: youtubeService,
	}, nil
}

func (y *Youtube) GetVideo(id string) *Video {
	response, err := y.client.Videos.List([]string{"snippet", "contentDetails"}).Id(id).Do()
	if err != nil {
		log.Errorf("Error while getting video: %s", err.Error())
		return nil
	}

	if len(response.Items) == 0 {
		return nil
	}

	duration, _ := time.ParseDuration(strings.TrimPrefix(strings.ToLower(response.Items[0].ContentDetails.Duration), "pt"))

	return &Video{
		Title:     response.Items[0].Snippet.Title,
		Thumbnail: getBestThumbnail(response.Items[0].Snippet.Thumbnails),
		ID:        id,
		Duration:  duration.Seconds(),
	}
}

func getBestThumbnail(thumbnails *youtube.ThumbnailDetails) string {
	if thumbnails.Maxres != nil {
		return thumbnails.Maxres.Url
	}

	if thumbnails.Standard != nil {
		return thumbnails.Standard.Url
	}

	if thumbnails.High != nil {
		return thumbnails.High.Url
	}

	if thumbnails.Medium != nil {
		return thumbnails.Medium.Url
	}

	if thumbnails.Default != nil {
		return thumbnails.Default.Url
	}

	return ""
}
