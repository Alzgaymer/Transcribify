create or replace function is_text_empty(str text) returns bool as
    $$
    begin
        return (str <> '') is not true or str is null;
    end;
    $$
    language plpgsql;

create or replace function get_user(
    p_email text,
    p_password text default null)
    returns setof users as
$$
begin

    if p_email is null or is_text_empty(p_email) then
        raise data_exception using hint ='cannot fund user with empty login', message ='enter login';
    else
        if p_password is null or is_text_empty(p_password) then
            return query    select *
                            from users u
                            where u.email = p_email;
        else
            return query    select *
                            from users u
                            where u.email = p_email
                            and   u.password = p_password;
        end if;
    end if;

end;
$$
    language plpgsql;

create or replace function put_user(
    p_email text, p_password text)
    returns setof users as
$$
begin

    if is_text_empty(p_email) or is_text_empty(p_password) then
        raise DATA_EXCEPTION using hint = 'enter user email and password', message = 'empty email or password';
    end if;

    insert into users (email, password)
    values (p_email, p_password)
    on conflict (email) do nothing;

    return query select * from users where email = p_email;
end;
$$
    language plpgsql;

create or replace function put_video(
    p_title text,
    p_description text,
    p_available_langs text[],
    p_length_in_seconds text,
    p_thumbnails jsonb,
    p_transcription jsonb,
    p_video_id char(11),
    p_language char(2)
    )
    returns setof video as
$$
begin
    insert into video (
    title, description, available_langs, length_in_seconds, thumbnails, transcription, video_id, language
    )
    values (
    p_title, p_description, p_available_langs, p_length_in_seconds, p_thumbnails, p_transcription, p_video_id, p_language
    )
    on conflict (video_id, language) do nothing;

    return query select * from video where video_id = p_video_id;

end;
$$
    language plpgsql;

create or replace function get_user_videos(
    p_user_id integer
    )
    returns setof user_videos as
$$
begin
    return query
    select * from user_videos
    where  user_id = p_user_id;
end;
$$
    language plpgsql;

create or replace procedure put_user_video(
    p_user_id integer,
    p_video_id char(11)
)  as
$$
begin

    if p_user_id <= 0 or is_text_empty(p_video_id::text) or length(p_video_id) <> 11 then
        raise data_exception using hint = 'enter user-id or video-id', message = 'empty user-id or video-id';
    end if;

    insert into user_videos(user_id, video_id)
    values (p_user_id, p_video_id);

end;
$$
    language plpgsql;