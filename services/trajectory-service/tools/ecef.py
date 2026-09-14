from pyproj import Transformer

import consts

# WGS-84:
# EPSG:4979 = geodetic 3D (lat, lon, height)
# EPSG:4978 = ECEF / geocentric

transformer = Transformer.from_crs(
    "EPSG:4979",
    "EPSG:4978",
    always_xy=True,
)

x, y, z = transformer.transform(consts.TEST_LON, consts.TEST_LAT, consts.TEST_ALT)

print(f"x = {x / 1000:.12f} km")
print(f"y = {y / 1000:.12f} km")
print(f"z = {z / 1000:.12f} km")
