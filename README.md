
# Transcribify

This project is a tool that allows users to easily obtain a transcription of a YouTube video. Using the [Youtube transcriptor](https://rapidapi.com/benrhzala90/api/youtube-transcriptor). The user simply needs to provide the YouTube video URL and the tool will handle the rest. This project is useful for anyone who needs a text version of a YouTube video, such as content creators who want to provide closed captions for their videos or people who want to save a transcript for personal use.


## Run Locally

Clone the project

```bash
  git clone https://github.com/Alzgaymer/transcribify
```

Go to the project directory

```bash
  cd transcribify
```

Install dependencies

```bash
  go mod download
```

Start the server. Make sure you add [.env](https://github.com/Alzgaymer/transcribify#environment-variables)


```bash
  make build & make run 
```
## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`APP_PORT`

`VIDEO_API_KEY` [Youtube transcriptor](https://rapidapi.com/benrhzala90/api/youtube-transcriptor) API key

*optional* `OPENAI_API_KEY`

`DB_USERNAME
DB_PASSWORD
DB_HOST
DB_PORT
DB_DATABASE`

`JWT_SALT`
## API Reference

#### Get video transcription (user autentification required)

```http
  GET /api/v1/video/{id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `lang` | `string` | **Required**. Transcription language |


#### Register

```http
  POST /api/v1/auth/sign-up
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email` | `string` | **Required**.  |
| `password` | `string` | **Required**.  |


#### Login

```http
  POST /api/v1/auth/login
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email` | `string` | **Required**.  |
| `password` | `string` | **Required**.  |


#### Get user searched videos

```http
  GET /api/v1/user/history/{page}?limit=
```
