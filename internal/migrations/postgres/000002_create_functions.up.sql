create or replace function get_user(p_login text, p_password text default null) returns setof users as
$$
begin
    if p_password is null then
        return query select *
                     from users u
                     where u.email = p_login;
    else
        return query select *
                     from users u
                     where u.email = p_login and u.password = p_password;
    end if;
end;
$$
    language plpgsql;


create or replace function put_user(p_email TEXT,p_password TEXT) returns setof users as
$$
begin
    INSERT INTO users (email, password)
    VALUES (p_email, p_password)
    ON CONFLICT (email) do nothing;

    RETURN QUERY SELECT * FROM users WHERE email = p_email;
end;
$$
    language plpgsql;