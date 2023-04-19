drop function get_user(p_login text, p_password text);
drop function put_user(p_email text,p_password text);
drop function put_video(
    in p_title text,
    in p_description text,
    in p_available_langs text[],
    in p_length_in_seconds text,
    in p_thumbnails jsonb,
    in p_transcription jsonb,
    in p_video_id char(11),
    in p_language char(2)
);
drop function get_user_videos(
    p_user_id integer
);
drop function put_user_video(
    p_user_id integer,
    p_video_id char(11)
) ;

drop function is_text_empty(str text);