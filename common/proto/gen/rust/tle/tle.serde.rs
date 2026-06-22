// @generated
impl serde::Serialize for GetTleRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.identifier.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tle.GetTleRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTleRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "identifier",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Identifier,
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
                            "identifier" => Ok(GeneratedField::Identifier),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTleRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tle.GetTleRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTleRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut identifier__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetTleRequest {
                    identifier: identifier__,
                })
            }
        }
        deserializer.deserialize_struct("tle.GetTleRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTleResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.tle.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tle.GetTleResponse", len)?;
        if let Some(v) = self.tle.as_ref() {
            struct_ser.serialize_field("tle", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTleResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "tle",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Tle,
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
                            "tle" => Ok(GeneratedField::Tle),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTleResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tle.GetTleResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTleResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut tle__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Tle => {
                            if tle__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tle"));
                            }
                            tle__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetTleResponse {
                    tle: tle__,
                })
            }
        }
        deserializer.deserialize_struct("tle.GetTleResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListTlesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("tle.ListTlesRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListTlesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListTlesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tle.ListTlesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListTlesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(ListTlesRequest {
                })
            }
        }
        deserializer.deserialize_struct("tle.ListTlesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListTlesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.tles.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tle.ListTlesResponse", len)?;
        if !self.tles.is_empty() {
            struct_ser.serialize_field("tles", &self.tles)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListTlesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "tles",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Tles,
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
                            "tles" => Ok(GeneratedField::Tles),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListTlesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tle.ListTlesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListTlesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut tles__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Tles => {
                            if tles__.is_some() {
                                return Err(serde::de::Error::duplicate_field("tles"));
                            }
                            tles__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ListTlesResponse {
                    tles: tles__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tle.ListTlesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Tle {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.norad_id != 0 {
            len += 1;
        }
        if !self.satellite_name.is_empty() {
            len += 1;
        }
        if !self.line1.is_empty() {
            len += 1;
        }
        if !self.line2.is_empty() {
            len += 1;
        }
        if self.epoch.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tle.Tle", len)?;
        if self.norad_id != 0 {
            struct_ser.serialize_field("noradId", &self.norad_id)?;
        }
        if !self.satellite_name.is_empty() {
            struct_ser.serialize_field("satelliteName", &self.satellite_name)?;
        }
        if !self.line1.is_empty() {
            struct_ser.serialize_field("line1", &self.line1)?;
        }
        if !self.line2.is_empty() {
            struct_ser.serialize_field("line2", &self.line2)?;
        }
        if let Some(v) = self.epoch.as_ref() {
            struct_ser.serialize_field("epoch", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Tle {
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
            "line1",
            "line2",
            "epoch",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            NoradId,
            SatelliteName,
            Line1,
            Line2,
            Epoch,
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
                            "line1" => Ok(GeneratedField::Line1),
                            "line2" => Ok(GeneratedField::Line2),
                            "epoch" => Ok(GeneratedField::Epoch),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Tle;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tle.Tle")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Tle, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut norad_id__ = None;
                let mut satellite_name__ = None;
                let mut line1__ = None;
                let mut line2__ = None;
                let mut epoch__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::NoradId => {
                            if norad_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("noradId"));
                            }
                            norad_id__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::SatelliteName => {
                            if satellite_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("satelliteName"));
                            }
                            satellite_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Line1 => {
                            if line1__.is_some() {
                                return Err(serde::de::Error::duplicate_field("line1"));
                            }
                            line1__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Line2 => {
                            if line2__.is_some() {
                                return Err(serde::de::Error::duplicate_field("line2"));
                            }
                            line2__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Epoch => {
                            if epoch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("epoch"));
                            }
                            epoch__ = map_.next_value()?;
                        }
                    }
                }
                Ok(Tle {
                    norad_id: norad_id__.unwrap_or_default(),
                    satellite_name: satellite_name__.unwrap_or_default(),
                    line1: line1__.unwrap_or_default(),
                    line2: line2__.unwrap_or_default(),
                    epoch: epoch__,
                })
            }
        }
        deserializer.deserialize_struct("tle.Tle", FIELDS, GeneratedVisitor)
    }
}
