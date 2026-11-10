from skyfield.api import EarthSatellite

import consts

satellite = EarthSatellite(consts.ISS_LINE_1, consts.ISS_LINE_2, "TEST")

ts = satellite.ts

difference = satellite - consts.TEST_OBSERVER
topocentric = difference.at(ts.from_datetime(consts.TEST_DATETIME))

alt, az, distance = topocentric.altaz()

print(f"azimuth_deg   = {az.degrees:.12f}")
print(f"elevation_deg = {alt.degrees:.12f}")
print(f"range_km      = {distance.km:.12f}")
