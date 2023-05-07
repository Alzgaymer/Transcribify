drop function get_user(p_login text);

drop procedure put_user(p_email text,p_password text);

drop procedure put_video(
    in p_title text,
    in p_description text,
    in p_available_langs text[],
    in p_length_in_seconds text,
    in p_thumbnails jsonb,
    in p_transcription jsonb,
    in p_video_id char(11),
    in p_language char(2)
);


drop function is_text_empty_or_null(str text);