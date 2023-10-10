package youtube

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nabomhalang/halangcordgo/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	log *config.Logger = config.NewLogger("youtube")
	YT  *Youtube
)

func init() {
	var err error

	if config.Get().YoutubeKey != "" {
		log.Infof("Initializing youtube...")
		YT, err = NewYoutube(config.Get().YoutubeKey)
		if err != nil {
			log.Fatalf("Failed to initialize youtube: %s", err.Error())
		}
	}
}

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

func (y *Youtube) Search(query string, maxReults int64) []Video {
	list, _ := y.client.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(maxReults).Type("video").Do()
	response, err := y.client.Videos.List([]string{"snippet", "contentDetails"}).Id(list.Items[0].Id.VideoId).Do()

	if err != nil {
		log.Errorf("Error while searching video: %s", err.Error())
		return nil
	}

	result := make([]Video, len(response.Items))
	for i, item := range response.Items {
		duration, _ := time.ParseDuration(strings.TrimPrefix(strings.ToLower(response.Items[0].ContentDetails.Duration), "pt"))

		result[i] = Video{
			Title:       item.Snippet.Title,
			Thumbnail:   getBestThumbnail(item.Snippet.Thumbnails),
			ID:          item.Id,
			Link:        fmt.Sprintf(YouTubeURL, item.Id),
			Description: item.Snippet.Description,
			Channel:     item.Snippet.ChannelTitle,
			Duration:    duration.Seconds(),
		}
	}

	return result
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
