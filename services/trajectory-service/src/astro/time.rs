use chrono::{DateTime, Datelike, Timelike, Utc};
use uom::si::angle::{degree, radian};
use uom::si::f64::{Angle, Time};
use uom::si::time::second;

use super::consts::{
    DAYS_IN_MONTH, DAYS_IN_YEAR, DAYS_PER_CENTURY, GMST_BASE, GMST_COEFF1, GMST_COEFF2,
    GMST_COEFF3, JULIAN_DAY_BASE, JULIAN_DAY_OFFSET, JULIAN_YEAR_OFFSET, SECONDS_TO_DEGREES,
    TWO_PI,
};

/// Convert UTC datetime to Greenwich Sidereal Time (GST) in radians.
///
/// GST is the angle between the Greenwich meridian and the vernal equinox,
/// and is used as an intermediate step for converting between ECI (Earth-Centered Inertial)
/// and ECEF (Earth-Centered Earth-Fixed) coordinate frames.
///
/// Reference formulas:
/// - Astronomical Almanac, 1984, p. B6
/// - Jean Meeus, *Astronomical Algorithms*, Ch. 12
/// - [SOFA/IAU Standards](https://www.iausofa.org/)
pub fn utc_to_gst(datetime: DateTime<Utc>) -> Angle {
    let jd = datetime_to_julian(&datetime);

    // julian centuries since J2000.0
    // J2000.0 corresponds to JD = 2451545.0 (2000-01-01 12:00 TT)
    let t = (jd - JULIAN_DAY_BASE) / DAYS_PER_CENTURY;

    // compute Greenwich Mean Sidereal Time (GMST) in seconds
    let gmst_sec = gmst_seconds(t);

    let angle_deg = gmst_sec.get::<second>() / SECONDS_TO_DEGREES;
    let angle = Angle::new::<degree>(angle_deg);

    // convert seconds to radians, normalize to [0, 2π)
    normalize_angle(angle)
}

/// Convert a UTC datetime into Julian Date (JD).
///
/// Julian Date is a continuous count of days since 4713 BC, used in astronomy.
/// Formula follows Meeus (1998), Ch. 7.
///
/// See: [NASA Julian Date explanation](https://ssd.jpl.nasa.gov/tools/jdc/#/)
fn datetime_to_julian(datetime: &DateTime<Utc>) -> f64 {
    let year = datetime.year();
    let month = datetime.month();
    let day = datetime.day();

    let hour = f64::from(datetime.hour());
    let min = f64::from(datetime.minute());
    let sec = f64::from(datetime.second()) + f64::from(datetime.nanosecond()) / 1_000_000_000.0;

    // adjust months so jan/feb are treated as months 13/14 of previous year
    let (yy, mm) = if month <= 2 {
        (year - 1, month + 12)
    } else {
        (year, month)
    };

    // gregorian calendar correction terms
    let a = (f64::from(yy) / 100.0).floor();
    let b = 2.0 - a + (a / 4.0).floor();

    // integer Julian day number (midnight-based)
    let jd_day = (DAYS_IN_YEAR * (f64::from(yy) + JULIAN_YEAR_OFFSET)).floor()
        + (DAYS_IN_MONTH * (f64::from(mm) + 1.0)).floor()
        + f64::from(day)
        + b
        - JULIAN_DAY_OFFSET;

    // fractional day from time-of-day
    let day_fraction = (hour + min / 60.0 + sec / 3600.0) / 24.0;

    jd_day + day_fraction
}

/// Compute GMST (Greenwich Mean Sidereal Time) in seconds.
///
/// Polynomial expansion in Julian centuries T since J2000.0:
///
/// ```text
/// GMST_sec = 67310.54841
///          + (876600*3600 + 8640184.812866) * T
///          + 0.093104 * T²
///          - 6.2e-6 * T³
/// ```
///
/// Formula from the IAU 1982 convention, also in Meeus (1998), Ch. 12.
fn gmst_seconds(t: f64) -> Time {
    let gmst = GMST_COEFF3
        .mul_add(t, GMST_COEFF2)
        .mul_add(t, GMST_COEFF1)
        .mul_add(t, GMST_BASE);

    Time::new::<second>(gmst)
}

/// Normalize an angle in radians to the range [0, 2π).
fn normalize_angle(angle: Angle) -> Angle {
    let wrapped = angle.get::<radian>() % TWO_PI;

    let wrapped = if wrapped < 0.0 {
        wrapped + TWO_PI
    } else {
        wrapped
    };

    Angle::new::<radian>(wrapped)
}

#[cfg(test)]
mod tests {
    use super::super::test_utils::assert_close;
    use super::*;

    use chrono::TimeZone;

    #[allow(clippy::expect_used)]
    fn utc(year: i32, month: u32, day: u32, hour: u32, minute: u32, seconds: u32) -> DateTime<Utc> {
        Utc.with_ymd_and_hms(year, month, day, hour, minute, seconds)
            .single()
            .expect("tests must have avalid UTC datetime")
    }

    #[test]
    fn julian_date_j2000() {
        // J2000.0 = 2000-01-01 12:00:00 UTC
        let datetime = utc(2000, 1, 1, 12, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_451_545.0, 1e-9);
    }

    #[test]
    fn julian_date_j2000_midnight() {
        let datetime = utc(2000, 1, 1, 0, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_451_544.5, 1e-9);
    }

    #[test]
    fn julian_date_1999() {
        let datetime = utc(1999, 1, 1, 0, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_451_179.5, 1e-9);
    }

    #[test]
    fn julian_date_1987() {
        let datetime = utc(1987, 1, 27, 0, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_446_822.5, 1e-9);
    }

    #[test]
    fn julian_date_with_fractional_day() {
        // 1987-06-19 12:00 UTC = 2446966.0
        let datetime = utc(1987, 6, 19, 12, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_446_966.0, 1e-9);
    }

    #[test]
    fn julian_date_handles_february() {
        // tests the month <= 2 branch
        let datetime = utc(1988, 1, 27, 0, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_447_187.5, 1e-9);
    }

    #[test]
    fn julian_date_handles_leap_year() {
        let datetime = utc(1988, 6, 19, 12, 0, 0);

        let jd = datetime_to_julian(&datetime);

        assert_close(jd, 2_447_332.0, 1e-9);
    }

    #[test]
    fn gmst_at_j2000() {
        // at J2000, T=0, so the polynomial reduces to GMST_BASE
        let result = gmst_seconds(0.0);

        assert_close(result.get::<second>(), GMST_BASE, 1e-9);
    }

    #[test]
    fn normalize_angle_leaves_normal_angle_unchanged() {
        let angle = Angle::new::<radian>(1.0);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), 1.0, 1e-12);
    }

    #[test]
    fn normalize_angle_wraps_angle_above_two_pi() {
        let angle = Angle::new::<radian>(TWO_PI + 1.0);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), 1.0, 1e-12);
    }

    #[test]
    fn normalize_angle_wraps_multiple_revolutions() {
        let angle = Angle::new::<radian>(3.0 * TWO_PI + 1.5);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), 1.5, 1e-12);
    }

    #[test]
    fn normalize_angle_wraps_negative_angle() {
        let angle = Angle::new::<radian>(-1.0);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), TWO_PI - 1.0, 1e-12);
    }

    #[test]
    fn normalize_angle_wraps_negative_multiple_revolutions() {
        let angle = Angle::new::<radian>(-3.0 * TWO_PI - 1.5);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), TWO_PI - 1.5, 1e-12);
    }

    #[test]
    fn normalize_angle_zero() {
        let angle = Angle::new::<radian>(0.0);

        let result = normalize_angle(angle);

        assert_close(result.get::<radian>(), 0.0, 1e-12);
    }

    #[test]
    fn gst_at_j2000() {
        // IAU 1982 GMST polynomial at J2000.0:
        //
        // GMST = 67310.54841 seconds
        //      = 280.460618375 degrees
        //      = approximately 4.8949612127 radians
        let datetime = utc(2000, 1, 1, 12, 0, 0);

        let result = utc_to_gst(datetime);

        assert_close(result.get::<radian>(), 4.894_961_212_823_059, 1e-9);
    }

    #[test]
    fn gst_is_normalized() {
        let datetime = utc(2000, 1, 1, 12, 0, 0);

        let result = utc_to_gst(datetime);
        let radians = result.get::<radian>();

        assert!(radians >= 0.0);
        assert!(radians < TWO_PI);
    }
}
