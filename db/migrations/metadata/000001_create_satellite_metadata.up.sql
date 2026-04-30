CREATE TABLE satellite_metadata
(
    norad_id           INTEGER PRIMARY KEY,
    cospar_id          TEXT,

    name               TEXT        NOT NULL,
    aliases            TEXT[]      NOT NULL DEFAULT '{}',

    object_type        TEXT        NOT NULL,
    mission_type       TEXT        NOT NULL,
    orbit_regime       TEXT        NOT NULL,

    operator           TEXT,
    owner              TEXT,
    constellation      TEXT,

    launch_date        TIMESTAMPTZ,
    launch_site        TEXT,
    launch_vehicle     TEXT,

    operational_status TEXT        NOT NULL,

    frequencies        JSONB       NOT NULL DEFAULT '[]',
    sources            JSONB       NOT NULL DEFAULT '[]',

    updated_at         TIMESTAMPTZ NOT NULL
);
