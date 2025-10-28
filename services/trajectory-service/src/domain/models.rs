use std::fmt;
use std::fmt::Formatter;

#[derive(Debug, Clone)]
pub enum SatelliteIdentifier {
    NoradId(u32),
    Name(String),
}

impl fmt::Display for SatelliteIdentifier {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match self {
            Self::NoradId(id) => write!(f, "NORAD ID {id}"),
            Self::Name(name) => write!(f, "Satellite '{name}'"),
        }
    }
}
