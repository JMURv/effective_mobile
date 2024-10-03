CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL,
    song_name TEXT NOT NULL,
    release_date DATE NOT NULL,
    lyrics TEXT[] NOT NULL ,
    link TEXT NOT NULL
);
