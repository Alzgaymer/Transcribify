create or replace function is_text_empty_or_null(str text) returns bool as
$$
begin
    return (str <> '') is not true or str is null;
end;
$$
language plpgsql;

create or replace function get_user(
    p_email text)
    returns setof users as
$$
begin

    if p_email is null or is_text_empty_or_null(p_email) then
        raise data_exception using hint ='cannot fund user with empty login', message ='enter login';
    else
        return query    select *
                        from users u
                        where u.email = p_email;
    end if;
end;
$$
    language plpgsql;

create or replace procedure put_user(
    p_email text, p_password text)
     as
$$
begin

    if is_text_empty_or_null(p_email) or is_text_empty_or_null(p_password) then
        raise DATA_EXCEPTION using hint = 'enter user email and password', message = 'empty email or password';
    end if;

    insert into users (email, password)
    values (p_email, p_password)
    on conflict (email) do nothing ;

end;
$$
    language plpgsql;

create or replace procedure put_video(
    p_title text,
    p_description text,
    p_available_langs text[],
    p_length_in_seconds text,
    p_thumbnails jsonb,
    p_transcription jsonb,
    p_video_id char(11),
    p_language char(2)
    )
     as
$$
begin
    insert into video (
    title, description, available_langs, length_in_seconds, thumbnails, transcription, video_id, language
    )
    values (
    p_title, p_description, p_available_langs, p_length_in_seconds, p_thumbnails, p_transcription, p_video_id, p_language
    );

end;
$$
    language plpgsql;

