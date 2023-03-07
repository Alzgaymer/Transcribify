package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"yt-video-transcriptor/logging"
	"yt-video-transcriptor/models"
)

func TestGetVideoTranscription(t *testing.T) {

	// Create a new router instance
	r := chi.NewRouter()
	logger, err := logging.New(
		logging.WithDevelopment(true),
	)
	if err != nil {
		t.Errorf(err.Error())
	}
	route := NewRoute(logger, &http.Client{
		Timeout: 30 * time.Second,
	})
	// Mount the GetVideoTranscription handler under the /api/v1 route
	r.Get("/api/v1/videos", route.GetVideoTranscription)

	// Define test cases in a table format
	testCases := []struct {
		name           string
		url            string
		responseStatus int
		responseData   []models.YTVideo
		responseErr    error
	}{
		{
			name:           "Success",
			url:            "/api/v1/videos?v=dQw4w9WgXcQ&lang=en",
			responseStatus: http.StatusOK,
			responseData: []models.YTVideo{
				{
					Title:           "Rick Astley - Never Gonna Give You Up (Official Music Video)",
					Description:     "The official video for “Never Gonna Give You Up” by Rick Astley\nTaken from the album ‘Whenever You Need Somebody’ – deluxe 2CD and digital deluxe out 6th May 2022 Pre-order here – https://RickAstley.lnk.to/WYNS2022ID\n\n“Never Gonna Give You Up” was a global smash on its release in July 1987, topping the charts in 25 countries including Rick’s native UK and the US Billboard Hot 100.  It also won the Brit Award for Best single in 1988. Stock Aitken and Waterman wrote and produced the track which was the lead-off single and lead track from Rick’s debut LP “Whenever You Need Somebody”.  The album was itself a UK number one and would go on to sell over 15 million copies worldwide.\n\nThe legendary video was directed by Simon West – who later went on to make Hollywood blockbusters such as Con Air, Lara Croft – Tomb Raider and The Expendables 2.  The video passed the 1bn YouTube views milestone on 28 July 2021.\n\nSubscribe to the official Rick Astley YouTube channel: https://RickAstley.lnk.to/YTSubID\n\nFollow Rick Astley:\nFacebook: https://RickAstley.lnk.to/FBFollowID \nTwitter: https://RickAstley.lnk.to/TwitterID \nInstagram: https://RickAstley.lnk.to/InstagramID \nWebsite: https://RickAstley.lnk.to/storeID \nTikTok: https://RickAstley.lnk.to/TikTokID\n\nListen to Rick Astley:\nSpotify: https://RickAstley.lnk.to/SpotifyID \nApple Music: https://RickAstley.lnk.to/AppleMusicID \nAmazon Music: https://RickAstley.lnk.to/AmazonMusicID \nDeezer: https://RickAstley.lnk.to/DeezerID \n\nLyrics:\nWe’re no strangers to love\nYou know the rules and so do I\nA full commitment’s what I’m thinking of\nYou wouldn’t get this from any other guy\n\nI just wanna tell you how I’m feeling\nGotta make you understand\n\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\n\nWe’ve known each other for so long\nYour heart’s been aching but you’re too shy to say it\nInside we both know what’s been going on\nWe know the game and we’re gonna play it\n\nAnd if you ask me how I’m feeling\nDon’t tell me you’re too blind to see\n\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\n\n#RickAstley #NeverGonnaGiveYouUp #WheneverYouNeedSomebody #OfficialMusicVideo",
					AvailableLangs:  []string{"en"},
					LengthInSeconds: "212",
					Thumbnails: []models.Thumbnails{
						{
							Url:    "https://i.ytimg.com/vi_webp/dQw4w9WgXcQ/default.webp",
							Width:  120,
							Height: 90,
						},
						{
							Url:    "https://i.ytimg.com/vi_webp/dQw4w9WgXcQ/mqdefault.webp",
							Width:  320,
							Height: 180,
						},
						{
							Url:    "https://i.ytimg.com/vi_webp/dQw4w9WgXcQ/hqdefault.webp",
							Width:  480,
							Height: 360,
						},
						{
							Url:    "https://i.ytimg.com/vi_webp/dQw4w9WgXcQ/sddefault.webp",
							Width:  640,
							Height: 480,
						},
						{
							Url:    "https://i.ytimg.com/vi/dQw4w9WgXcQ/hq720.jpg?sqp=-oaymwEcCK4FEIIDSEbyq4qpAw4IARUAAIhCGAFwAcABBg==\u0026rs=AOn4CLAtyUwHA-QTnSIeRMIY_9t9RnBjkA",
							Width:  686,
							Height: 386,
						},
					},
					Transcription: []models.Transcription{
						{
							Subtitle: "[Music]",
							Start:    0,
							Dur:      26.359,
						},
						{
							Subtitle: "you know the rules",
							Start:    22.64,
							Dur:      3.719,
						},
						{
							Subtitle: "[Music]",
							Start:    28.33,
							Dur:      16.31,
						},
						{
							Subtitle: "gotta make you understand",
							Start:    40.399,
							Dur:      4.241,
						},
						{
							Subtitle: "[Music]",
							Start:    44.92,
							Dur:      11.75,
						},
						{
							Subtitle: "goodbye",
							Start:    54.64,
							Dur:      6.079,
						},
						{
							Subtitle: "[Music]",
							Start:    56.67,
							Dur:      6.13,
						},
						{
							Subtitle: "we\u0026#39;ve known each other",
							Start:    60.719,
							Dur:      4.16,
						},
						{
							Subtitle: "for so long",
							Start:    62.8,
							Dur:      2.9,
						},
						{
							Subtitle: "your heart\u0026#39;s been",
							Start:    64.879,
							Dur:      6.401,
						},
						{
							Subtitle: "[Music]",
							Start:    65.7,
							Dur:      3.06,
						},
						{
							Subtitle: "going aching",
							Start:    71.28,
							Dur:      0,
						},
						{
							Subtitle: "[Music]",
							Start:    90.53,
							Dur:      3.15,
						},
						{
							Subtitle: "never gonna say goodbye",
							Start:    95.759,
							Dur:      14.801,
						},
						{
							Subtitle: "[Music]",
							Start:    99,
							Dur:      13.96,
						},
						{
							Subtitle: "never gonna make you",
							Start:    110.56,
							Dur:      4.28,
						},
						{
							Subtitle: "gonna say cry",
							Start:    112.96,
							Dur:      1.88,
						},
						{
							Subtitle: "[Music]",
							Start:    115.86,
							Dur:      32.22,
						},
						{
							Subtitle: "i",
							Start:    153.44,
							Dur:      6.04,
						},
						{
							Subtitle: "just want to tell you how i\u0026#39;m feeling",
							Start:    154.64,
							Dur:      4.84,
						},
						{
							Subtitle: "[Music]",
							Start:    160.01,
							Dur:      19.18,
						},
						{
							Subtitle: "[Music]",
							Start:    183.47,
							Dur:      16.87,
						},
						{
							Subtitle: "never gonna is you down",
							Start:    197.12,
							Dur:      13.069,
						},
						{
							Subtitle: "[Music]",
							Start:    200.34,
							Dur:      9.849,
						},
					},
				},
			},
			responseErr: nil,
		},
		{
			name:           "Invalid Input",
			url:            "/api/v1//",
			responseStatus: http.StatusNotFound,
			responseData:   nil,
			responseErr:    nil,
		},
		// Add more test cases as needed
	}
	server := httptest.NewServer(r)
	defer server.Client()
	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {

			response, err := http.Get(server.URL + testCase.url)
			require.Nil(t, err)

			split := strings.Split(response.Status, " ")
			require.NotNil(t, split)

			status, err := strconv.Atoi(split[0])
			require.Nil(t, err)

			// Status
			assert.Equal(t, testCase.responseStatus, status)

			//Get video from response
			video, _ := route.responseToYTVideo(response)

			// Title
			assert.Equal(t, testCase.responseData, video)

		})
	}

}
