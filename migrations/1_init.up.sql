create table if not exists users (
    id       integer primary key,
    login    text not null unique,
    passHash blob not null,
    isadmin  boolean not null default false
);

create index if not exists  idx_login on users(login);

create table if not exists apps (
    id     int64  integer primary key,
    name   text not null unique,
    secret text not null unique
);
