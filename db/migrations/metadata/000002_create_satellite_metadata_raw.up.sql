CREATE TABLE satellite_metadata_raw
(
    id         BIGSERIAL PRIMARY KEY,

    norad_id   INT         NOT NULL,
    cospar_id  TEXT,

    source     TEXT        NOT NULL,

    payload    JSONB       NOT NULL,

    fetched_at TIMESTAMPTZ NOT NULL,
    stored_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sat_metadata_raw_norad
    ON satellite_metadata_raw (norad_id);

CREATE INDEX idx_sat_metadata_raw_source
    ON satellite_metadata_raw (source);

CREATE INDEX idx_sat_metadata_raw_fetched_at
    ON satellite_metadata_raw (fetched_at DESC);