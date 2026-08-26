from datetime import timedelta

from skyfield.api import EarthSatellite, load, wgs84

import consts

START = consts.TEST_DATETIME
END = START + timedelta(hours=6)

MIN_ELEVATION_DEG = 0.0

ts = load.timescale()

satellite = EarthSatellite(consts.ISS_LINE_1, consts.ISS_LINE_2, ts=ts)

observer = wgs84.latlon(
    consts.TEST_LAT,
    consts.TEST_LON,
    elevation_m=consts.TEST_ALT * 1000.0,
)

t0 = ts.from_datetime(START)
t1 = ts.from_datetime(END)

times, events = satellite.find_events(
    observer,
    t0,
    t1,
    altitude_degrees=MIN_ELEVATION_DEG,
)

print("PASS EVENTS")

for t, event in zip(times, events):
    dt = t.utc_datetime()

    names = {
        0: "AOS",
        1: "MAX",
        2: "LOS",
    }

    print(
        names[event],
        dt.isoformat(),
    )

    if event == 1:
        difference = satellite - observer
        topocentric = difference.at(t)
        alt, az, distance = topocentric.altaz()

        print(f"  elevation_deg = {alt.degrees:.15f}")
        print(f"  azimuth_deg   = {az.degrees:.15f}")
        print(f"  range_km      = {distance.km:.15f}")
