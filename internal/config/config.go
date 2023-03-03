package config

type AppConfiguration struct {
	Port string `env:"APP_PORT"`
}

type APIConfiguration struct {
	Key string `env:"API_KEY"`
	API string `env:"API_URL"`
}

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

type Transcription struct {
	Subtitle any     `json:"subtitle"`
	Start    float64 `json:"start"`
	Dur      float64 `json:"dur"`
}
