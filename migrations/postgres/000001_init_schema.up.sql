CREATE TABLE IF NOT EXISTS video_data (
        id SERIAL PRIMARY KEY,

        video_id VARCHAR(11) NOT NULL,
        language VARCHAR(2) NOT NULL,
        json_data JSONB NOT NULL
);
