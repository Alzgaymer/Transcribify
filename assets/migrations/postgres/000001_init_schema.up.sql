create table IF NOT EXISTS video (
        id serial primary key,
        title text,
        description text,
        available_langs text[],
        length_in_seconds text, --- not integer because response from API is "160"
        thumbnails jsonb,
        transcription jsonb,
        video_id text CONSTRAINT valid_video_id CHECK ( length(video_id) = 11) ,
        language text CONSTRAINT valid_language CHECK ( length(language) = 2 )
);

create table IF NOT EXISTS users (
        id serial primary key,
        email text unique ,
        password text not null
);

create table IF NOT EXISTS user_videos (
        id serial primary key,
        user_id int not null,
        video_id int not null,
        foreign key (user_id) references users (id),
        foreign key (video_id) references video (id)
);
