// @generated
impl serde::Serialize for GeodeticInput {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lat.is_some() {
            len += 1;
        }
        if self.lon.is_some() {
            len += 1;
        }
        if self.alt.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.GeodeticInput", len)?;
        if let Some(v) = self.lat.as_ref() {
            match v {
                geodetic_input::Lat::LatDeg(v) => {
                    struct_ser.serialize_field("latDeg", v)?;
                }
                geodetic_input::Lat::LatRad(v) => {
                    struct_ser.serialize_field("latRad", v)?;
                }
            }
        }
        if let Some(v) = self.lon.as_ref() {
            match v {
                geodetic_input::Lon::LonDeg(v) => {
                    struct_ser.serialize_field("lonDeg", v)?;
                }
                geodetic_input::Lon::LonRad(v) => {
                    struct_ser.serialize_field("lonRad", v)?;
                }
            }
        }
        if let Some(v) = self.alt.as_ref() {
            match v {
                geodetic_input::Alt::AltM(v) => {
                    struct_ser.serialize_field("altM", v)?;
                }
                geodetic_input::Alt::AltKm(v) => {
                    struct_ser.serialize_field("altKm", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GeodeticInput {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lat_deg",
            "latDeg",
            "lat_rad",
            "latRad",
            "lon_deg",
            "lonDeg",
            "lon_rad",
            "lonRad",
            "alt_m",
            "altM",
            "alt_km",
            "altKm",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LatDeg,
            LatRad,
            LonDeg,
            LonRad,
            AltM,
            AltKm,
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
                            "latDeg" | "lat_deg" => Ok(GeneratedField::LatDeg),
                            "latRad" | "lat_rad" => Ok(GeneratedField::LatRad),
                            "lonDeg" | "lon_deg" => Ok(GeneratedField::LonDeg),
                            "lonRad" | "lon_rad" => Ok(GeneratedField::LonRad),
                            "altM" | "alt_m" => Ok(GeneratedField::AltM),
                            "altKm" | "alt_km" => Ok(GeneratedField::AltKm),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GeodeticInput;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.GeodeticInput")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GeodeticInput, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lat__ = None;
                let mut lon__ = None;
                let mut alt__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LatDeg => {
                            if lat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("latDeg"));
                            }
                            lat__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Lat::LatDeg(x.0));
                        }
                        GeneratedField::LatRad => {
                            if lat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("latRad"));
                            }
                            lat__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Lat::LatRad(x.0));
                        }
                        GeneratedField::LonDeg => {
                            if lon__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lonDeg"));
                            }
                            lon__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Lon::LonDeg(x.0));
                        }
                        GeneratedField::LonRad => {
                            if lon__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lonRad"));
                            }
                            lon__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Lon::LonRad(x.0));
                        }
                        GeneratedField::AltM => {
                            if alt__.is_some() {
                                return Err(serde::de::Error::duplicate_field("altM"));
                            }
                            alt__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Alt::AltM(x.0));
                        }
                        GeneratedField::AltKm => {
                            if alt__.is_some() {
                                return Err(serde::de::Error::duplicate_field("altKm"));
                            }
                            alt__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| geodetic_input::Alt::AltKm(x.0));
                        }
                    }
                }
                Ok(GeodeticInput {
                    lat: lat__,
                    lon: lon__,
                    alt: alt__,
                })
            }
        }
        deserializer.deserialize_struct("common.GeodeticInput", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GeodeticOutput {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lat != 0. {
            len += 1;
        }
        if self.lon != 0. {
            len += 1;
        }
        if self.alt != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.GeodeticOutput", len)?;
        if self.lat != 0. {
            struct_ser.serialize_field("lat", &self.lat)?;
        }
        if self.lon != 0. {
            struct_ser.serialize_field("lon", &self.lon)?;
        }
        if self.alt != 0. {
            struct_ser.serialize_field("alt", &self.alt)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GeodeticOutput {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lat",
            "lon",
            "alt",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Lat,
            Lon,
            Alt,
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
                            "lat" => Ok(GeneratedField::Lat),
                            "lon" => Ok(GeneratedField::Lon),
                            "alt" => Ok(GeneratedField::Alt),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GeodeticOutput;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.GeodeticOutput")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GeodeticOutput, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lat__ = None;
                let mut lon__ = None;
                let mut alt__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Lat => {
                            if lat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lat"));
                            }
                            lat__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Lon => {
                            if lon__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lon"));
                            }
                            lon__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Alt => {
                            if alt__.is_some() {
                                return Err(serde::de::Error::duplicate_field("alt"));
                            }
                            alt__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GeodeticOutput {
                    lat: lat__.unwrap_or_default(),
                    lon: lon__.unwrap_or_default(),
                    alt: alt__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("common.GeodeticOutput", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SamplingOptions {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.step_seconds != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.SamplingOptions", len)?;
        if self.step_seconds != 0 {
            struct_ser.serialize_field("stepSeconds", &self.step_seconds)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SamplingOptions {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "step_seconds",
            "stepSeconds",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StepSeconds,
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
                            "stepSeconds" | "step_seconds" => Ok(GeneratedField::StepSeconds),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SamplingOptions;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.SamplingOptions")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SamplingOptions, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut step_seconds__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StepSeconds => {
                            if step_seconds__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stepSeconds"));
                            }
                            step_seconds__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(SamplingOptions {
                    step_seconds: step_seconds__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("common.SamplingOptions", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SatelliteIdentifier {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.kind.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.SatelliteIdentifier", len)?;
        if let Some(v) = self.kind.as_ref() {
            match v {
                satellite_identifier::Kind::NoradId(v) => {
                    struct_ser.serialize_field("noradId", v)?;
                }
                satellite_identifier::Kind::SatelliteName(v) => {
                    struct_ser.serialize_field("satelliteName", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SatelliteIdentifier {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "norad_id",
            "noradId",
            "satellite_name",
            "satelliteName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            NoradId,
            SatelliteName,
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
                            "noradId" | "norad_id" => Ok(GeneratedField::NoradId),
                            "satelliteName" | "satellite_name" => Ok(GeneratedField::SatelliteName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SatelliteIdentifier;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.SatelliteIdentifier")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SatelliteIdentifier, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut kind__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::NoradId => {
                            if kind__.is_some() {
                                return Err(serde::de::Error::duplicate_field("noradId"));
                            }
                            kind__ = map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| satellite_identifier::Kind::NoradId(x.0));
                        }
                        GeneratedField::SatelliteName => {
                            if kind__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satelliteName"));
                            }
                            kind__ = map_.next_value::<::std::option::Option<_>>()?.map(satellite_identifier::Kind::SatelliteName);
                        }
                    }
                }
                Ok(SatelliteIdentifier {
                    kind: kind__,
                })
            }
        }
        deserializer.deserialize_struct("common.SatelliteIdentifier", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TimeRange {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start.is_some() {
            len += 1;
        }
        if self.end.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.TimeRange", len)?;
        if let Some(v) = self.start.as_ref() {
            struct_ser.serialize_field("start", v)?;
        }
        if let Some(v) = self.end.as_ref() {
            struct_ser.serialize_field("end", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TimeRange {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start",
            "end",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Start,
            End,
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
                            "start" => Ok(GeneratedField::Start),
                            "end" => Ok(GeneratedField::End),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TimeRange;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.TimeRange")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TimeRange, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut start__ = None;
                let mut end__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Start => {
                            if start__.is_some() {
                                return Err(serde::de::Error::duplicate_field("start"));
                            }
                            start__ = map_.next_value()?;
                        }
                        GeneratedField::End => {
                            if end__.is_some() {
                                return Err(serde::de::Error::duplicate_field("end"));
                            }
                            end__ = map_.next_value()?;
                        }
                    }
                }
                Ok(TimeRange {
                    start: start__,
                    end: end__,
                })
            }
        }
        deserializer.deserialize_struct("common.TimeRange", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UnitSettings {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.distance_unit != 0 {
            len += 1;
        }
        if self.angle_unit != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.UnitSettings", len)?;
        if self.distance_unit != 0 {
            let v = unit_settings::DistanceUnit::try_from(self.distance_unit)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.distance_unit)))?;
            struct_ser.serialize_field("distanceUnit", &v)?;
        }
        if self.angle_unit != 0 {
            let v = unit_settings::AngleUnit::try_from(self.angle_unit)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.angle_unit)))?;
            struct_ser.serialize_field("angleUnit", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UnitSettings {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "distance_unit",
            "distanceUnit",
            "angle_unit",
            "angleUnit",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            DistanceUnit,
            AngleUnit,
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
                            "distanceUnit" | "distance_unit" => Ok(GeneratedField::DistanceUnit),
                            "angleUnit" | "angle_unit" => Ok(GeneratedField::AngleUnit),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UnitSettings;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.UnitSettings")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UnitSettings, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut distance_unit__ = None;
                let mut angle_unit__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::DistanceUnit => {
                            if distance_unit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("distanceUnit"));
                            }
                            distance_unit__ = Some(map_.next_value::<unit_settings::DistanceUnit>()? as i32);
                        }
                        GeneratedField::AngleUnit => {
                            if angle_unit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("angleUnit"));
                            }
                            angle_unit__ = Some(map_.next_value::<unit_settings::AngleUnit>()? as i32);
                        }
                    }
                }
                Ok(UnitSettings {
                    distance_unit: distance_unit__.unwrap_or_default(),
                    angle_unit: angle_unit__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("common.UnitSettings", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for unit_settings::AngleUnit {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "ANGLE_UNIT_UNSPECIFIED",
            Self::Degrees => "ANGLE_UNIT_DEGREES",
            Self::Radians => "ANGLE_UNIT_RADIANS",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for unit_settings::AngleUnit {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ANGLE_UNIT_UNSPECIFIED",
            "ANGLE_UNIT_DEGREES",
            "ANGLE_UNIT_RADIANS",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = unit_settings::AngleUnit;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "ANGLE_UNIT_UNSPECIFIED" => Ok(unit_settings::AngleUnit::Unspecified),
                    "ANGLE_UNIT_DEGREES" => Ok(unit_settings::AngleUnit::Degrees),
                    "ANGLE_UNIT_RADIANS" => Ok(unit_settings::AngleUnit::Radians),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for unit_settings::DistanceUnit {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "DISTANCE_UNIT_UNSPECIFIED",
            Self::Meters => "DISTANCE_UNIT_METERS",
            Self::Kilometers => "DISTANCE_UNIT_KILOMETERS",
            Self::Miles => "DISTANCE_UNIT_MILES",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for unit_settings::DistanceUnit {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "DISTANCE_UNIT_UNSPECIFIED",
            "DISTANCE_UNIT_METERS",
            "DISTANCE_UNIT_KILOMETERS",
            "DISTANCE_UNIT_MILES",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = unit_settings::DistanceUnit;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "DISTANCE_UNIT_UNSPECIFIED" => Ok(unit_settings::DistanceUnit::Unspecified),
                    "DISTANCE_UNIT_METERS" => Ok(unit_settings::DistanceUnit::Meters),
                    "DISTANCE_UNIT_KILOMETERS" => Ok(unit_settings::DistanceUnit::Kilometers),
                    "DISTANCE_UNIT_MILES" => Ok(unit_settings::DistanceUnit::Miles),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for Vector3 {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.x != 0. {
            len += 1;
        }
        if self.y != 0. {
            len += 1;
        }
        if self.z != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("common.Vector3", len)?;
        if self.x != 0. {
            struct_ser.serialize_field("x", &self.x)?;
        }
        if self.y != 0. {
            struct_ser.serialize_field("y", &self.y)?;
        }
        if self.z != 0. {
            struct_ser.serialize_field("z", &self.z)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Vector3 {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "x",
            "y",
            "z",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            X,
            Y,
            Z,
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
                            "x" => Ok(GeneratedField::X),
                            "y" => Ok(GeneratedField::Y),
                            "z" => Ok(GeneratedField::Z),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Vector3;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct common.Vector3")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Vector3, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut x__ = None;
                let mut y__ = None;
                let mut z__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::X => {
                            if x__.is_some() {
                                return Err(serde::de::Error::duplicate_field("x"));
                            }
                            x__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Y => {
                            if y__.is_some() {
                                return Err(serde::de::Error::duplicate_field("y"));
                            }
                            y__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Z => {
                            if z__.is_some() {
                                return Err(serde::de::Error::duplicate_field("z"));
                            }
                            z__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(Vector3 {
                    x: x__.unwrap_or_default(),
                    y: y__.unwrap_or_default(),
                    z: z__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("common.Vector3", FIELDS, GeneratedVisitor)
    }
}
