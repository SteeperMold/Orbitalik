use chrono::{DateTime, TimeZone, Utc};
use prost_types::Timestamp;
use crate::domain::errors::TimestampConversionError;

pub trait ToChrono {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError>;
}

impl ToChrono for Timestamp {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError> {
        let nanos = u32::try_from(self.nanos)?;

        Utc.timestamp_opt(self.seconds, nanos)
            .single()
            .ok_or(TimestampConversionError::InvalidTimestamp)
    }
}

pub trait ToProtoTimestamp {
    fn to_proto_timestamp(&self) -> Result<Timestamp, TimestampConversionError>;
}

impl ToProtoTimestamp for DateTime<Utc> {
    fn to_proto_timestamp(&self) -> Result<Timestamp, TimestampConversionError> {
        let nanos = i32::try_from(self.timestamp_subsec_nanos())?;

        Ok(Timestamp {
            seconds: self.timestamp(),
            nanos,
        })
    }
}