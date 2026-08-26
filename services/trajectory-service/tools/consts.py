import datetime as dt

from skyfield.api import wgs84

ISS_LINE_1 = "1 25544U 98067A   26232.17880947  .00009753  00000+0  18154-3 0  9994"
ISS_LINE_2 = "2 25544  51.6332 343.3775 0007674  65.0551 295.1235 15.49524101581676"

TEST_DATETIME = dt.datetime(
    2026, 1, 1,
    0, 0, 0,
    tzinfo=dt.timezone.utc,
)

TEST_LAT = 56.9496
TEST_LON = 24.1052
TEST_ALT = 0

TEST_OBSERVER = wgs84.latlon(
    latitude_degrees=TEST_LAT,
    longitude_degrees=TEST_LON,
    elevation_m=TEST_ALT,
)
