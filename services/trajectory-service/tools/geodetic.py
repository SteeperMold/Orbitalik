from pyproj import Transformer

# WGS-84 ECEF -> geodetic
#
# EPSG:4978 = geocentric / ECEF
# EPSG:4979 = geographic 3D / latitude, longitude, ellipsoidal height

transformer = Transformer.from_crs(
    "EPSG:4978",
    "EPSG:4979",
    always_xy=True,
)

# ECEF coordinates in kilometres.
x_km = 3182.645531427102
y_km = 1424.012805324485
z_km = 5322.841235217519

lon_deg, lat_deg, alt_m = transformer.transform(
    x_km * 1000.0,
    y_km * 1000.0,
    z_km * 1000.0,
)

print(f"latitude_deg  = {lat_deg:.12f}")
print(f"longitude_deg = {lon_deg:.12f}")
print(f"altitude_km   = {alt_m / 1000.0:.12f}")