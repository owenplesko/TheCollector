DROP TABLE summoner_match;
DROP TABLE match;
DROP TABLE rank;
DROP INDEX matches_last_updated;
DROP TABLE summoner;

CREATE TABLE summoner (
    puuid TEXT NOT NULL PRIMARY KEY,
    region TEXT NOT NULL,
    summoner_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    profile_icon_id INT NOT NULL,
    revision_date BIGINT NOT NULL,
    display_name TEXT NOT NULL,
    raw_name TEXT NOT NULL,
    summoner_level INT NOT NULL,
    last_updated BIGINT NOT NULL DEFAULT 0,
    rank_last_updated BIGINT NOT NULL DEFAULT 0,
    matches_last_updated BIGINT NOT NULL DEFAULT 0,

    UNIQUE(raw_name, region)
);

CREATE INDEX matches_last_updated ON summoner(matches_last_updated ASC);

CREATE TABLE rank (
    summoner_puuid TEXT PRIMARY KEY NOT NULL,
    CONSTRAINT fK_summoner_puuid FOREIGN KEY (summoner_puuid) REFERENCES summoner(puuid),

    summoner_region TEXT NOT NULL,
    summoner_id     TEXT NOT NULL,
    queue_type      TEXT NOT NULL,
    tier            TEXT NOT NULL,
    rank            TEXT NOT NULL,
    league_points   INT NOT NULL,
    wins            INT NOT NULL,
    losses          INT NOT NULL,
    hot_streak      BOOLEAN NOT NULL,
    veteran         BOOLEAN NOT NULL,
    fresh_blood     BOOLEAN NOT NULL,
    inactive        BOOLEAN NOT NULL
);

CREATE TABLE match (
    id          TEXT PRIMARY KEY NOT NULL,
    date BIGINT NOT NULL,
    data_version TEXT,
    data JSONB
);

CREATE TABLE summoner_match (
    summoner_puuid TEXT NOT NULL,
    CONSTRAINT fK_summoner_puuid FOREIGN KEY (summoner_puuid) REFERENCES summoner(puuid),
    match_id TEXT NOT NULL,
    CONSTRAINT fk_match_id FOREIGN KEY (match_id) REFERENCES match(id),
    PRIMARY KEY(summoner_puuid, match_id)
)