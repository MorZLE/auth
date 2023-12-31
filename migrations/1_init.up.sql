create table if not exists users (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    login    text not null unique,
    passHash blob not null,
    app_id   INTEGER not null,
    foreign key(app_id) references apps(id)
);

create table if not exists admins (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER not null,
    lvl     INTEGER not null,
    app_id  INTEGER not null,
    foreign key(user_id) references users(id),
    foreign key(app_id) references apps(id)
);

create index if not exists idx_login on users(login);

create table if not exists apps (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    name   text not null unique,
    secret text not null unique
);
