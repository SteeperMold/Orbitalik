from math import atan2, cos, sin, sqrt, radians, degrees

from sgp4.api import Satrec, jday, WGS72

import consts

A = 6378.137
E2 = 0.0066943799901413165
TWO_PI = 2.0 * 3.141592653589793


def assert_close(name, actual, expected, tolerance):
    error = abs(actual - expected)

    if error > tolerance:
        raise AssertionError(
            f"{name}: expected {expected}, "
            f"got {actual}, error {error} > tolerance {tolerance}"
        )


def print_vector(name, vector):
    print(f"{name}:")
    print(f"  x = {vector[0]:.15f} km")
    print(f"  y = {vector[1]:.15f} km")
    print(f"  z = {vector[2]:.15f} km")


def datetime_to_julian(dt):
    year = dt.year
    month = dt.month
    day = dt.day

    hour = dt.hour
    minute = dt.minute
    second = dt.second

    if month <= 2:
        year -= 1
        month += 12

    a = year // 100
    b = 2 - a + a // 4

    jd_day = (
            int(365.25 * (year + 4716))
            + int(30.6001 * (month + 1))
            + day
            + b
            - 1524.5
    )

    day_fraction = (
                           hour
                           + minute / 60.0
                           + second / 3600.0
                   ) / 24.0

    return jd_day + day_fraction


def utc_to_gmst(dt):
    jd = datetime_to_julian(dt)

    t = (jd - 2451545.0) / 36525.0

    gmst_seconds = (
            67310.54841
            + (876600.0 * 3600.0 + 8640184.812866) * t
            + 0.093104 * t * t
            - 6.2e-6 * t * t * t
    )

    gmst_degrees = gmst_seconds / 240.0

    gmst_degrees %= 360.0

    return radians(gmst_degrees)


def teme_to_ecef(teme, gst):
    x, y, z = teme

    sin_gst = sin(gst)
    cos_gst = cos(gst)

    x_ecef = cos_gst * x + sin_gst * y
    y_ecef = -sin_gst * x + cos_gst * y
    z_ecef = z

    return (
        x_ecef,
        y_ecef,
        z_ecef,
    )


def ecef_to_geodetic(ecef):
    x, y, z = ecef

    lon = atan2(y, x)

    r = sqrt(x * x + y * y)

    lat = atan2(
        z,
        r * (1.0 - E2),
    )

    while True:
        sin_lat = sin(lat)
        cos_lat = cos(lat)

        n = A / sqrt(
            1.0 - E2 * sin_lat * sin_lat
        )

        h = r / cos_lat - n

        new_lat = atan2(
            E2 * n * sin_lat + z,
            r,
        )

        if abs(new_lat - lat) < 1e-12:
            lat = new_lat
            break

        lat = new_lat

    return (
        degrees(lat),
        degrees(lon),
        h,
    )


def main():
    satellite = Satrec.twoline2rv(
        consts.ISS_LINE_1,
        consts.ISS_LINE_2,
        WGS72,
    )

    jd, fr = jday(
        consts.TEST_DATETIME.year,
        consts.TEST_DATETIME.month,
        consts.TEST_DATETIME.day,
        consts.TEST_DATETIME.hour,
        consts.TEST_DATETIME.minute,
        consts.TEST_DATETIME.second,
    )

    error, position, velocity = satellite.sgp4(jd, fr)

    if error != 0:
        raise RuntimeError(
            f"SGP4 propagation failed with error code {error}"
        )

    teme = position

    print()
    print("==================================================")
    print("TEME")
    print("==================================================")

    print_vector("TEME", teme)

    gst = utc_to_gmst(consts.TEST_DATETIME)

    ecef = teme_to_ecef(
        teme,
        gst,
    )

    print()
    print("==================================================")
    print("ECEF")
    print("==================================================")

    print_vector("ECEF", ecef)

    print()
    print(f"GMST = {degrees(gst):.15f} degrees")
    print(f"GMST = {gst:.15f} radians")

    latitude, longitude, altitude = ecef_to_geodetic(ecef)

    print()
    print("==================================================")
    print("GEODETIC")
    print("==================================================")

    print(
        f"latitude  = {latitude:.15f} degrees"
    )

    print(
        f"longitude = {longitude:.15f} degrees"
    )

    print(
        f"altitude  = {altitude:.15f} km"
    )

    print()
    print("==================================================")
    print("TIME")
    print("==================================================")

    print(
        f"UTC = {consts.TEST_DATETIME.isoformat()}"
    )


if __name__ == "__main__":
    main()
