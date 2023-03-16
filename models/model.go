package models

import (
	"encoding/json"
	"strconv"
)

type YTVideo struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	AvailableLangs  []string        `json:"availableLangs"`
	LengthInSeconds string          `json:"lengthInSeconds"`
	Thumbnails      []Thumbnails    `json:"thumbnails"`
	Transcription   []Transcription `json:"transcription"`
}

type Thumbnails struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Transcription implements json.Unmarshaler
type Transcription struct {
	Subtitle string  `json:"subtitle"`
	Start    float64 `json:"start"`
	Dur      float64 `json:"dur"`
}

func (t *Transcription) UnmarshalJSON(data []byte) error {
	type rawString struct {
		Subtitle string  `json:"subtitle"`
		Start    float64 `json:"start"`
		Dur      float64 `json:"dur"`
	}

	rawStr := new(rawString)

	err := json.Unmarshal(data, rawStr)
	if err == nil {
		t.Dur = rawStr.Dur
		t.Start = rawStr.Start
		t.Subtitle = rawStr.Subtitle
		return nil
	}

	type rawInt struct {
		Subtitle int     `json:"subtitle"`
		Start    float64 `json:"start"`
		Dur      float64 `json:"dur"`
	}
	rawI := rawInt{}

	err = json.Unmarshal(data, &rawI)
	if err != nil {
		return err
	}

	t.Start = rawI.Start
	t.Dur = rawI.Dur
	t.Subtitle = strconv.Itoa(rawI.Subtitle)
	return nil
}

type VideoRequest struct {
	VideoID  string `json:"v"`
	Language string `json:"lang"`
}

func YTVideoToJsonb(videos []YTVideo) (json.RawMessage, error) {
	var jsonData []byte
	var err error

	if len(videos) > 0 {
		jsonData, err = json.Marshal(videos)
		if err != nil {
			return nil, err
		}
	} else {
		jsonData = []byte("[]")
	}

	return jsonData, nil
}

func RawMessageToYTVideo(raw json.RawMessage) ([]YTVideo, error) {
	ytVideo := make([]YTVideo, 0)
	err := json.Unmarshal(raw, &ytVideo)
	if err != nil {
		return nil, err
	}
	return ytVideo, nil
}
