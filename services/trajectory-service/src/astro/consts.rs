pub const A: f64 = 6378.137; // Equatorial radius in km
pub const F: f64 = 1.0 / 298.257_223_563; // Flattening
pub const E2: f64 = F * (2.0 - F); // Square of eccentricity

// Julian day constants
pub const JULIAN_DAY_OFFSET: f64 = 1524.5;
pub const JULIAN_DAY_BASE: f64 = 2_451_545.0; // J2000 epoch
pub const DAYS_PER_CENTURY: f64 = 36525.0;

// Calendar conversion constants
pub const JULIAN_YEAR_OFFSET: f64 = 4716.0;
pub const DAYS_IN_YEAR: f64 = 365.25;
pub const DAYS_IN_MONTH: f64 = 30.6001;

// GMST calculation constants
pub const GMST_BASE: f64 = 67310.54841;
pub const GMST_COEFF1: f64 = 876_600.0 * 3600.0 + 8_640_184.812_866; // ≈ 3.156008e9
pub const GMST_COEFF2: f64 = 0.093_104;
pub const GMST_COEFF3: f64 = -6.2e-6;

// Conversion
pub const SECONDS_TO_DEGREES: f64 = 240.0;
pub const TWO_PI: f64 = std::f64::consts::TAU;
