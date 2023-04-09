CREATE TABLE video (
         id SERIAL PRIMARY KEY,
         title TEXT,
         description TEXT,
         available_langs TEXT[],
         length_in_seconds TEXT,
         thumbnails JSONB,
         transcription JSONB,
         video_id CHAR(11),
         language CHAR(2),
         UNIQUE(video_id, language)
);

CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name TEXT NOT NULL,
       email TEXT UNIQUE NOT NULL,
       password TEXT NOT NULL,
       refresh_token TEXT
);

CREATE TABLE user_videos (
         id SERIAL PRIMARY KEY,
         user_id INT NOT NULL,
         video_id CHAR(11) NOT NULL,
         FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
         FOREIGN KEY (video_id) REFERENCES video (video_id) ON DELETE CASCADE
);

