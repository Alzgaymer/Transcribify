package models

import (
	"encoding/json"
	"strconv"
)

type YTVideo struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	AvailableLangs  []string        `json:"availableLangs"`  //nolint:tagliatelle
	LengthInSeconds string          `json:"lengthInSeconds"` //nolint:tagliatelle
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
	VideoID  string `json:"v" validate:"len=11,ascii"`
	Language string `json:"lang" validate:"bcp47_language_tag"`
}
