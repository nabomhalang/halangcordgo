package sponsorblock

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/nabomhalang/halangcordgo/config"
)

const (
	frameSeconds = 50.00067787
)

var (
	log *config.Logger = config.NewLogger("sponsorblock")
)

// GetSegments returns a map for skipping certain frames of a song
func GetSegments(videoID string) map[int]bool {
	// Gets segments
	req, _ := http.NewRequest("GET", "https://sponsor.ajay.app/api/skipSegments/"+hash(videoID)+"?categories=[\"sponsor\",\"music_offtopic\"]", nil)
	client := http.Client{Timeout: time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Can't get SponsorBlock segments: %s", err)
		return nil
	}

	if resp.StatusCode == http.StatusOK {
		var (
			videos     SponsorBlock
			segmentMap = make(map[int]bool)
		)

		err = json.NewDecoder(resp.Body).Decode(&videos)
		_ = resp.Body.Close()
		if err != nil {
			log.Errorf("Can't unmarshal JSON, %s", err)
			return nil
		}

		for _, v := range videos {
			if v.VideoID == videoID {
				for _, segment := range v.Segments {
					if len(segment.Segment) == 2 {
						segmentMap[int(segment.Segment[0]*frameSeconds)] = true
						segmentMap[int(segment.Segment[1]*frameSeconds)] = true
					}
				}
				return segmentMap
			}
		}
	}

	return nil
}

// returns the first 4 characters of a sha256 hash
func hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))[:4]
}
