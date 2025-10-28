CREATE TABLE IF NOT EXISTS tle
(
    id             SERIAL PRIMARY KEY,
    norad_id       INT  NOT NULL,
    satellite_name TEXT NOT NULL,
    line1          TEXT NOT NULL,
    line2          TEXT NOT NULL,
    epoch          TIMESTAMP WITH TIME ZONE,
    fetched_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),

    UNIQUE (norad_id, epoch)
);

CREATE INDEX IF NOT EXISTS idx_tle_norad ON tle (norad_id);
CREATE INDEX IF NOT EXISTS idx_tle_name ON tle (satellite_name);
