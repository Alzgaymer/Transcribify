create table video (
        id serial primary key,
        title text,
        description text,
        available_langs text[],
        length_in_seconds text,
        thumbnails jsonb,
        transcription jsonb,
        video_id char(11) unique ,
        language char(2) unique ,
        unique (video_id, language)
);

create table users (
        id serial primary key,
        email text unique ,
        password text not null,
        unique (id, email)
);

create table user_videos (
        id serial primary key,
        user_id int not null,
        video_id char(11) not null,
        foreign key (user_id) references users (id) on delete cascade,
        foreign key (video_id) references video (video_id) on delete cascade

);
