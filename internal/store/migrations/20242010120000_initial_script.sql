-- +migrate Up

CREATE table songs (
    id uuid primary key,
    release_date varchar,
    name varchar not null,
    music_group varchar not null,
    text varchar,
    link  varchar,
    deleted bool,

    CONSTRAINT unique_song UNIQUE (release_date, name, music_group)
);

-- +migrate Down

DROP TABLE songs;
