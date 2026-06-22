// @generated
impl serde::Serialize for NextPassesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.satellites.is_empty() {
            len += 1;
        }
        if self.observer.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        if self.min_peak_elevation != 0. {
            len += 1;
        }
        if self.min_elevation != 0. {
            len += 1;
        }
        if self.count != 0 {
            len += 1;
        }
        if self.lookahead_range.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("passes.NextPassesRequest", len)?;
        if !self.satellites.is_empty() {
            struct_ser.serialize_field("satellites", &self.satellites)?;
        }
        if let Some(v) = self.observer.as_ref() {
            struct_ser.serialize_field("observer", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        if self.min_peak_elevation != 0. {
            struct_ser.serialize_field("minPeakElevation", &self.min_peak_elevation)?;
        }
        if self.min_elevation != 0. {
            struct_ser.serialize_field("minElevation", &self.min_elevation)?;
        }
        if self.count != 0 {
            struct_ser.serialize_field("count", &self.count)?;
        }
        if let Some(v) = self.lookahead_range.as_ref() {
            struct_ser.serialize_field("lookaheadRange", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for NextPassesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "satellites",
            "observer",
            "units",
            "min_peak_elevation",
            "minPeakElevation",
            "min_elevation",
            "minElevation",
            "count",
            "lookahead_range",
            "lookaheadRange",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Satellites,
            Observer,
            Units,
            MinPeakElevation,
            MinElevation,
            Count,
            LookaheadRange,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "satellites" => Ok(GeneratedField::Satellites),
                            "observer" => Ok(GeneratedField::Observer),
                            "units" => Ok(GeneratedField::Units),
                            "minPeakElevation" | "min_peak_elevation" => Ok(GeneratedField::MinPeakElevation),
                            "minElevation" | "min_elevation" => Ok(GeneratedField::MinElevation),
                            "count" => Ok(GeneratedField::Count),
                            "lookaheadRange" | "lookahead_range" => Ok(GeneratedField::LookaheadRange),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = NextPassesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct passes.NextPassesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<NextPassesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut satellites__ = None;
                let mut observer__ = None;
                let mut units__ = None;
                let mut min_peak_elevation__ = None;
                let mut min_elevation__ = None;
                let mut count__ = None;
                let mut lookahead_range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Satellites => {
                            if satellites__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satellites"));
                            }
                            satellites__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Observer => {
                            if observer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("observer"));
                            }
                            observer__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                        GeneratedField::MinPeakElevation => {
                            if min_peak_elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minPeakElevation"));
                            }
                            min_peak_elevation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MinElevation => {
                            if min_elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minElevation"));
                            }
                            min_elevation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Count => {
                            if count__.is_some() {
                                return Err(serde::de::Error::duplicate_field("count"));
                            }
                            count__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LookaheadRange => {
                            if lookahead_range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lookaheadRange"));
                            }
                            lookahead_range__ = map_.next_value()?;
                        }
                    }
                }
                Ok(NextPassesRequest {
                    satellites: satellites__.unwrap_or_default(),
                    observer: observer__,
                    units: units__,
                    min_peak_elevation: min_peak_elevation__.unwrap_or_default(),
                    min_elevation: min_elevation__.unwrap_or_default(),
                    count: count__.unwrap_or_default(),
                    lookahead_range: lookahead_range__,
                })
            }
        }
        deserializer.deserialize_struct("passes.NextPassesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Pass {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.satellite.is_some() {
            len += 1;
        }
        if self.aos.is_some() {
            len += 1;
        }
        if self.aos_azimuth != 0. {
            len += 1;
        }
        if self.max_elevation_time.is_some() {
            len += 1;
        }
        if self.max_elevation != 0. {
            len += 1;
        }
        if self.max_elevation_azimuth != 0. {
            len += 1;
        }
        if self.los.is_some() {
            len += 1;
        }
        if self.los_azimuth != 0. {
            len += 1;
        }
        if self.duration_seconds != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("passes.Pass", len)?;
        if let Some(v) = self.satellite.as_ref() {
            struct_ser.serialize_field("satellite", v)?;
        }
        if let Some(v) = self.aos.as_ref() {
            struct_ser.serialize_field("aos", v)?;
        }
        if self.aos_azimuth != 0. {
            struct_ser.serialize_field("aosAzimuth", &self.aos_azimuth)?;
        }
        if let Some(v) = self.max_elevation_time.as_ref() {
            struct_ser.serialize_field("maxElevationTime", v)?;
        }
        if self.max_elevation != 0. {
            struct_ser.serialize_field("maxElevation", &self.max_elevation)?;
        }
        if self.max_elevation_azimuth != 0. {
            struct_ser.serialize_field("maxElevationAzimuth", &self.max_elevation_azimuth)?;
        }
        if let Some(v) = self.los.as_ref() {
            struct_ser.serialize_field("los", v)?;
        }
        if self.los_azimuth != 0. {
            struct_ser.serialize_field("losAzimuth", &self.los_azimuth)?;
        }
        if self.duration_seconds != 0 {
            struct_ser.serialize_field("durationSeconds", &self.duration_seconds)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Pass {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "satellite",
            "aos",
            "aos_azimuth",
            "aosAzimuth",
            "max_elevation_time",
            "maxElevationTime",
            "max_elevation",
            "maxElevation",
            "max_elevation_azimuth",
            "maxElevationAzimuth",
            "los",
            "los_azimuth",
            "losAzimuth",
            "duration_seconds",
            "durationSeconds",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Satellite,
            Aos,
            AosAzimuth,
            MaxElevationTime,
            MaxElevation,
            MaxElevationAzimuth,
            Los,
            LosAzimuth,
            DurationSeconds,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "satellite" => Ok(GeneratedField::Satellite),
                            "aos" => Ok(GeneratedField::Aos),
                            "aosAzimuth" | "aos_azimuth" => Ok(GeneratedField::AosAzimuth),
                            "maxElevationTime" | "max_elevation_time" => Ok(GeneratedField::MaxElevationTime),
                            "maxElevation" | "max_elevation" => Ok(GeneratedField::MaxElevation),
                            "maxElevationAzimuth" | "max_elevation_azimuth" => Ok(GeneratedField::MaxElevationAzimuth),
                            "los" => Ok(GeneratedField::Los),
                            "losAzimuth" | "los_azimuth" => Ok(GeneratedField::LosAzimuth),
                            "durationSeconds" | "duration_seconds" => Ok(GeneratedField::DurationSeconds),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Pass;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct passes.Pass")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Pass, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut satellite__ = None;
                let mut aos__ = None;
                let mut aos_azimuth__ = None;
                let mut max_elevation_time__ = None;
                let mut max_elevation__ = None;
                let mut max_elevation_azimuth__ = None;
                let mut los__ = None;
                let mut los_azimuth__ = None;
                let mut duration_seconds__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Satellite => {
                            if satellite__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satellite"));
                            }
                            satellite__ = map_.next_value()?;
                        }
                        GeneratedField::Aos => {
                            if aos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("aos"));
                            }
                            aos__ = map_.next_value()?;
                        }
                        GeneratedField::AosAzimuth => {
                            if aos_azimuth__.is_some() {
                                return Err(serde::de::Error::duplicate_field("aosAzimuth"));
                            }
                            aos_azimuth__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MaxElevationTime => {
                            if max_elevation_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxElevationTime"));
                            }
                            max_elevation_time__ = map_.next_value()?;
                        }
                        GeneratedField::MaxElevation => {
                            if max_elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxElevation"));
                            }
                            max_elevation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MaxElevationAzimuth => {
                            if max_elevation_azimuth__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxElevationAzimuth"));
                            }
                            max_elevation_azimuth__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Los => {
                            if los__.is_some() {
                                return Err(serde::de::Error::duplicate_field("los"));
                            }
                            los__ = map_.next_value()?;
                        }
                        GeneratedField::LosAzimuth => {
                            if los_azimuth__.is_some() {
                                return Err(serde::de::Error::duplicate_field("losAzimuth"));
                            }
                            los_azimuth__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::DurationSeconds => {
                            if duration_seconds__.is_some() {
                                return Err(serde::de::Error::duplicate_field("durationSeconds"));
                            }
                            duration_seconds__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(Pass {
                    satellite: satellite__,
                    aos: aos__,
                    aos_azimuth: aos_azimuth__.unwrap_or_default(),
                    max_elevation_time: max_elevation_time__,
                    max_elevation: max_elevation__.unwrap_or_default(),
                    max_elevation_azimuth: max_elevation_azimuth__.unwrap_or_default(),
                    los: los__,
                    los_azimuth: los_azimuth__.unwrap_or_default(),
                    duration_seconds: duration_seconds__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("passes.Pass", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PassPredictionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.satellites.is_empty() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        if self.observer.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        if self.min_peak_elevation != 0. {
            len += 1;
        }
        if self.min_elevation != 0. {
            len += 1;
        }
        if self.max_results != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("passes.PassPredictionRequest", len)?;
        if !self.satellites.is_empty() {
            struct_ser.serialize_field("satellites", &self.satellites)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        if let Some(v) = self.observer.as_ref() {
            struct_ser.serialize_field("observer", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        if self.min_peak_elevation != 0. {
            struct_ser.serialize_field("minPeakElevation", &self.min_peak_elevation)?;
        }
        if self.min_elevation != 0. {
            struct_ser.serialize_field("minElevation", &self.min_elevation)?;
        }
        if self.max_results != 0 {
            struct_ser.serialize_field("maxResults", &self.max_results)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PassPredictionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "satellites",
            "range",
            "observer",
            "units",
            "min_peak_elevation",
            "minPeakElevation",
            "min_elevation",
            "minElevation",
            "max_results",
            "maxResults",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Satellites,
            Range,
            Observer,
            Units,
            MinPeakElevation,
            MinElevation,
            MaxResults,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "satellites" => Ok(GeneratedField::Satellites),
                            "range" => Ok(GeneratedField::Range),
                            "observer" => Ok(GeneratedField::Observer),
                            "units" => Ok(GeneratedField::Units),
                            "minPeakElevation" | "min_peak_elevation" => Ok(GeneratedField::MinPeakElevation),
                            "minElevation" | "min_elevation" => Ok(GeneratedField::MinElevation),
                            "maxResults" | "max_results" => Ok(GeneratedField::MaxResults),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PassPredictionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct passes.PassPredictionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PassPredictionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut satellites__ = None;
                let mut range__ = None;
                let mut observer__ = None;
                let mut units__ = None;
                let mut min_peak_elevation__ = None;
                let mut min_elevation__ = None;
                let mut max_results__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Satellites => {
                            if satellites__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satellites"));
                            }
                            satellites__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                        GeneratedField::Observer => {
                            if observer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("observer"));
                            }
                            observer__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                        GeneratedField::MinPeakElevation => {
                            if min_peak_elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minPeakElevation"));
                            }
                            min_peak_elevation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MinElevation => {
                            if min_elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minElevation"));
                            }
                            min_elevation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MaxResults => {
                            if max_results__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxResults"));
                            }
                            max_results__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(PassPredictionRequest {
                    satellites: satellites__.unwrap_or_default(),
                    range: range__,
                    observer: observer__,
                    units: units__,
                    min_peak_elevation: min_peak_elevation__.unwrap_or_default(),
                    min_elevation: min_elevation__.unwrap_or_default(),
                    max_results: max_results__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("passes.PassPredictionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PassPredictionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.metadata.is_some() {
            len += 1;
        }
        if !self.passes.is_empty() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("passes.PassPredictionResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        if !self.passes.is_empty() {
            struct_ser.serialize_field("passes", &self.passes)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PassPredictionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
            "passes",
            "range",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
            Passes,
            Range,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "metadata" => Ok(GeneratedField::Metadata),
                            "passes" => Ok(GeneratedField::Passes),
                            "range" => Ok(GeneratedField::Range),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PassPredictionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct passes.PassPredictionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PassPredictionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                let mut passes__ = None;
                let mut range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                        GeneratedField::Passes => {
                            if passes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("passes"));
                            }
                            passes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                    }
                }
                Ok(PassPredictionResponse {
                    metadata: metadata__,
                    passes: passes__.unwrap_or_default(),
                    range: range__,
                })
            }
        }
        deserializer.deserialize_struct("passes.PassPredictionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PassesComputationMetadata {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.propagation_model.is_empty() {
            len += 1;
        }
        if self.computation_time.is_some() {
            len += 1;
        }
        if !self.norad_ids.is_empty() {
            len += 1;
        }
        if !self.satellite_names.is_empty() {
            len += 1;
        }
        if self.tle_epoch.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        if self.satellites_evaluated != 0 {
            len += 1;
        }
        if self.passes_found != 0 {
            len += 1;
        }
        if self.computation_ms != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("passes.PassesComputationMetadata", len)?;
        if !self.propagation_model.is_empty() {
            struct_ser.serialize_field("propagationModel", &self.propagation_model)?;
        }
        if let Some(v) = self.computation_time.as_ref() {
            struct_ser.serialize_field("computationTime", v)?;
        }
        if !self.norad_ids.is_empty() {
            struct_ser.serialize_field("noradIds", &self.norad_ids)?;
        }
        if !self.satellite_names.is_empty() {
            struct_ser.serialize_field("satelliteNames", &self.satellite_names)?;
        }
        if let Some(v) = self.tle_epoch.as_ref() {
            struct_ser.serialize_field("tleEpoch", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        if self.satellites_evaluated != 0 {
            struct_ser.serialize_field("satellitesEvaluated", &self.satellites_evaluated)?;
        }
        if self.passes_found != 0 {
            struct_ser.serialize_field("passesFound", &self.passes_found)?;
        }
        if self.computation_ms != 0 {
            struct_ser.serialize_field("computationMs", &self.computation_ms)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PassesComputationMetadata {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "propagation_model",
            "propagationModel",
            "computation_time",
            "computationTime",
            "norad_ids",
            "noradIds",
            "satellite_names",
            "satelliteNames",
            "tle_epoch",
            "tleEpoch",
            "units",
            "satellites_evaluated",
            "satellitesEvaluated",
            "passes_found",
            "passesFound",
            "computation_ms",
            "computationMs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PropagationModel,
            ComputationTime,
            NoradIds,
            SatelliteNames,
            TleEpoch,
            Units,
            SatellitesEvaluated,
            PassesFound,
            ComputationMs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "propagationModel" | "propagation_model" => Ok(GeneratedField::PropagationModel),
                            "computationTime" | "computation_time" => Ok(GeneratedField::ComputationTime),
                            "noradIds" | "norad_ids" => Ok(GeneratedField::NoradIds),
                            "satelliteNames" | "satellite_names" => Ok(GeneratedField::SatelliteNames),
                            "tleEpoch" | "tle_epoch" => Ok(GeneratedField::TleEpoch),
                            "units" => Ok(GeneratedField::Units),
                            "satellitesEvaluated" | "satellites_evaluated" => Ok(GeneratedField::SatellitesEvaluated),
                            "passesFound" | "passes_found" => Ok(GeneratedField::PassesFound),
                            "computationMs" | "computation_ms" => Ok(GeneratedField::ComputationMs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PassesComputationMetadata;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct passes.PassesComputationMetadata")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PassesComputationMetadata, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut propagation_model__ = None;
                let mut computation_time__ = None;
                let mut norad_ids__ = None;
                let mut satellite_names__ = None;
                let mut tle_epoch__ = None;
                let mut units__ = None;
                let mut satellites_evaluated__ = None;
                let mut passes_found__ = None;
                let mut computation_ms__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PropagationModel => {
                            if propagation_model__.is_some() {
                                return Err(serde::de::Error::duplicate_field("propagationModel"));
                            }
                            propagation_model__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ComputationTime => {
                            if computation_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("computationTime"));
                            }
                            computation_time__ = map_.next_value()?;
                        }
                        GeneratedField::NoradIds => {
                            if norad_ids__.is_some() {
                                return Err(serde::de::Error::duplicate_field("noradIds"));
                            }
                            norad_ids__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::SatelliteNames => {
                            if satellite_names__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satelliteNames"));
                            }
                            satellite_names__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TleEpoch => {
                            if tle_epoch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tleEpoch"));
                            }
                            tle_epoch__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                        GeneratedField::SatellitesEvaluated => {
                            if satellites_evaluated__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satellitesEvaluated"));
                            }
                            satellites_evaluated__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PassesFound => {
                            if passes_found__.is_some() {
                                return Err(serde::de::Error::duplicate_field("passesFound"));
                            }
                            passes_found__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ComputationMs => {
                            if computation_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("computationMs"));
                            }
                            computation_ms__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(PassesComputationMetadata {
                    propagation_model: propagation_model__.unwrap_or_default(),
                    computation_time: computation_time__,
                    norad_ids: norad_ids__.unwrap_or_default(),
                    satellite_names: satellite_names__.unwrap_or_default(),
                    tle_epoch: tle_epoch__,
                    units: units__,
                    satellites_evaluated: satellites_evaluated__.unwrap_or_default(),
                    passes_found: passes_found__.unwrap_or_default(),
                    computation_ms: computation_ms__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("passes.PassesComputationMetadata", FIELDS, GeneratedVisitor)
    }
}
