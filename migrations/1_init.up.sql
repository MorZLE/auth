create table if not exists users (
    id       INTEGER PRIMARY KEY  AUTOINCREMENT,
    login    text not null unique,
    passHash blob not null,
    isadmin  boolean not null default false
);

create index if not exists  idx_login on users(login);

create table if not exists apps (
    id     INTEGER PRIMARY KEY   AUTOINCREMENT,
    name   text not null unique,
    secret text not null unique
);
