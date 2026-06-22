// @generated
impl serde::Serialize for Frequency {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.direction != 0 {
            len += 1;
        }
        if self.frequency_mhz != 0. {
            len += 1;
        }
        if self.bandwidth_khz != 0. {
            len += 1;
        }
        if !self.modulation.is_empty() {
            len += 1;
        }
        if !self.mode.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metadata.Frequency", len)?;
        if self.direction != 0 {
            let v = FrequencyDirection::try_from(self.direction)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.direction)))?;
            struct_ser.serialize_field("direction", &v)?;
        }
        if self.frequency_mhz != 0. {
            struct_ser.serialize_field("frequencyMhz", &self.frequency_mhz)?;
        }
        if self.bandwidth_khz != 0. {
            struct_ser.serialize_field("bandwidthKhz", &self.bandwidth_khz)?;
        }
        if !self.modulation.is_empty() {
            struct_ser.serialize_field("modulation", &self.modulation)?;
        }
        if !self.mode.is_empty() {
            struct_ser.serialize_field("mode", &self.mode)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Frequency {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "direction",
            "frequency_mhz",
            "frequencyMhz",
            "bandwidth_khz",
            "bandwidthKhz",
            "modulation",
            "mode",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Direction,
            FrequencyMhz,
            BandwidthKhz,
            Modulation,
            Mode,
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
                            "direction" => Ok(GeneratedField::Direction),
                            "frequencyMhz" | "frequency_mhz" => Ok(GeneratedField::FrequencyMhz),
                            "bandwidthKhz" | "bandwidth_khz" => Ok(GeneratedField::BandwidthKhz),
                            "modulation" => Ok(GeneratedField::Modulation),
                            "mode" => Ok(GeneratedField::Mode),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Frequency;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.Frequency")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Frequency, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut direction__ = None;
                let mut frequency_mhz__ = None;
                let mut bandwidth_khz__ = None;
                let mut modulation__ = None;
                let mut mode__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Direction => {
                            if direction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("direction"));
                            }
                            direction__ = Some(map_.next_value::<FrequencyDirection>()? as i32);
                        }
                        GeneratedField::FrequencyMhz => {
                            if frequency_mhz__.is_some() {
                                return Err(serde::de::Error::duplicate_field("frequencyMhz"));
                            }
                            frequency_mhz__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::BandwidthKhz => {
                            if bandwidth_khz__.is_some() {
                                return Err(serde::de::Error::duplicate_field("bandwidthKhz"));
                            }
                            bandwidth_khz__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Modulation => {
                            if modulation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("modulation"));
                            }
                            modulation__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Mode => {
                            if mode__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mode"));
                            }
                            mode__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(Frequency {
                    direction: direction__.unwrap_or_default(),
                    frequency_mhz: frequency_mhz__.unwrap_or_default(),
                    bandwidth_khz: bandwidth_khz__.unwrap_or_default(),
                    modulation: modulation__.unwrap_or_default(),
                    mode: mode__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metadata.Frequency", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for FrequencyDirection {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "FREQUENCY_DIRECTION_UNSPECIFIED",
            Self::Uplink => "FREQUENCY_DIRECTION_UPLINK",
            Self::Downlink => "FREQUENCY_DIRECTION_DOWNLINK",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for FrequencyDirection {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "FREQUENCY_DIRECTION_UNSPECIFIED",
            "FREQUENCY_DIRECTION_UPLINK",
            "FREQUENCY_DIRECTION_DOWNLINK",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = FrequencyDirection;

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
                    "FREQUENCY_DIRECTION_UNSPECIFIED" => Ok(FrequencyDirection::Unspecified),
                    "FREQUENCY_DIRECTION_UPLINK" => Ok(FrequencyDirection::Uplink),
                    "FREQUENCY_DIRECTION_DOWNLINK" => Ok(FrequencyDirection::Downlink),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for GetMetadataRequest {
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
        let mut struct_ser = serializer.serialize_struct("metadata.GetMetadataRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetMetadataRequest {
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
            type Value = GetMetadataRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.GetMetadataRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetMetadataRequest, V::Error>
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
                Ok(GetMetadataRequest {
                    identifier: identifier__,
                })
            }
        }
        deserializer.deserialize_struct("metadata.GetMetadataRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetMetadataResponse {
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
        let mut struct_ser = serializer.serialize_struct("metadata.GetMetadataResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetMetadataResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
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
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetMetadataResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.GetMetadataResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetMetadataResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetMetadataResponse {
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("metadata.GetMetadataResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListSatelliteMetadataRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.object_type.is_some() {
            len += 1;
        }
        if self.mission_type.is_some() {
            len += 1;
        }
        if self.operational_status.is_some() {
            len += 1;
        }
        if self.orbit_regime.is_some() {
            len += 1;
        }
        if self.constellation.is_some() {
            len += 1;
        }
        if self.page_size != 0 {
            len += 1;
        }
        if !self.page_token.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metadata.ListSatelliteMetadataRequest", len)?;
        if let Some(v) = self.object_type.as_ref() {
            let v = ObjectType::try_from(*v)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", *v)))?;
            struct_ser.serialize_field("objectType", &v)?;
        }
        if let Some(v) = self.mission_type.as_ref() {
            let v = MissionType::try_from(*v)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", *v)))?;
            struct_ser.serialize_field("missionType", &v)?;
        }
        if let Some(v) = self.operational_status.as_ref() {
            let v = OperationalStatus::try_from(*v)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", *v)))?;
            struct_ser.serialize_field("operationalStatus", &v)?;
        }
        if let Some(v) = self.orbit_regime.as_ref() {
            let v = OrbitRegime::try_from(*v)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", *v)))?;
            struct_ser.serialize_field("orbitRegime", &v)?;
        }
        if let Some(v) = self.constellation.as_ref() {
            struct_ser.serialize_field("constellation", v)?;
        }
        if self.page_size != 0 {
            struct_ser.serialize_field("pageSize", &self.page_size)?;
        }
        if !self.page_token.is_empty() {
            struct_ser.serialize_field("pageToken", &self.page_token)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListSatelliteMetadataRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object_type",
            "objectType",
            "mission_type",
            "missionType",
            "operational_status",
            "operationalStatus",
            "orbit_regime",
            "orbitRegime",
            "constellation",
            "page_size",
            "pageSize",
            "page_token",
            "pageToken",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ObjectType,
            MissionType,
            OperationalStatus,
            OrbitRegime,
            Constellation,
            PageSize,
            PageToken,
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
                            "objectType" | "object_type" => Ok(GeneratedField::ObjectType),
                            "missionType" | "mission_type" => Ok(GeneratedField::MissionType),
                            "operationalStatus" | "operational_status" => Ok(GeneratedField::OperationalStatus),
                            "orbitRegime" | "orbit_regime" => Ok(GeneratedField::OrbitRegime),
                            "constellation" => Ok(GeneratedField::Constellation),
                            "pageSize" | "page_size" => Ok(GeneratedField::PageSize),
                            "pageToken" | "page_token" => Ok(GeneratedField::PageToken),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListSatelliteMetadataRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.ListSatelliteMetadataRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListSatelliteMetadataRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object_type__ = None;
                let mut mission_type__ = None;
                let mut operational_status__ = None;
                let mut orbit_regime__ = None;
                let mut constellation__ = None;
                let mut page_size__ = None;
                let mut page_token__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ObjectType => {
                            if object_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("objectType"));
                            }
                            object_type__ = map_.next_value::<::std::option::Option<ObjectType>>()?.map(|x| x as i32);
                        }
                        GeneratedField::MissionType => {
                            if mission_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("missionType"));
                            }
                            mission_type__ = map_.next_value::<::std::option::Option<MissionType>>()?.map(|x| x as i32);
                        }
                        GeneratedField::OperationalStatus => {
                            if operational_status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("operationalStatus"));
                            }
                            operational_status__ = map_.next_value::<::std::option::Option<OperationalStatus>>()?.map(|x| x as i32);
                        }
                        GeneratedField::OrbitRegime => {
                            if orbit_regime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("orbitRegime"));
                            }
                            orbit_regime__ = map_.next_value::<::std::option::Option<OrbitRegime>>()?.map(|x| x as i32);
                        }
                        GeneratedField::Constellation => {
                            if constellation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("constellation"));
                            }
                            constellation__ = map_.next_value()?;
                        }
                        GeneratedField::PageSize => {
                            if page_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageSize"));
                            }
                            page_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PageToken => {
                            if page_token__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pageToken"));
                            }
                            page_token__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ListSatelliteMetadataRequest {
                    object_type: object_type__,
                    mission_type: mission_type__,
                    operational_status: operational_status__,
                    orbit_regime: orbit_regime__,
                    constellation: constellation__,
                    page_size: page_size__.unwrap_or_default(),
                    page_token: page_token__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metadata.ListSatelliteMetadataRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListSatelliteMetadataResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.items.is_empty() {
            len += 1;
        }
        if !self.next_page_token.is_empty() {
            len += 1;
        }
        if self.total != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metadata.ListSatelliteMetadataResponse", len)?;
        if !self.items.is_empty() {
            struct_ser.serialize_field("items", &self.items)?;
        }
        if !self.next_page_token.is_empty() {
            struct_ser.serialize_field("nextPageToken", &self.next_page_token)?;
        }
        if self.total != 0 {
            struct_ser.serialize_field("total", &self.total)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListSatelliteMetadataResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "items",
            "next_page_token",
            "nextPageToken",
            "total",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Items,
            NextPageToken,
            Total,
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
                            "items" => Ok(GeneratedField::Items),
                            "nextPageToken" | "next_page_token" => Ok(GeneratedField::NextPageToken),
                            "total" => Ok(GeneratedField::Total),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListSatelliteMetadataResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.ListSatelliteMetadataResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListSatelliteMetadataResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut items__ = None;
                let mut next_page_token__ = None;
                let mut total__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Items => {
                            if items__.is_some() {
                                return Err(serde::de::Error::duplicate_field("items"));
                            }
                            items__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NextPageToken => {
                            if next_page_token__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nextPageToken"));
                            }
                            next_page_token__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Total => {
                            if total__.is_some() {
                                return Err(serde::de::Error::duplicate_field("total"));
                            }
                            total__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(ListSatelliteMetadataResponse {
                    items: items__.unwrap_or_default(),
                    next_page_token: next_page_token__.unwrap_or_default(),
                    total: total__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metadata.ListSatelliteMetadataResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MissionType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "MISSION_TYPE_UNSPECIFIED",
            Self::Communications => "MISSION_TYPE_COMMUNICATIONS",
            Self::EarthObservation => "MISSION_TYPE_EARTH_OBSERVATION",
            Self::Navigation => "MISSION_TYPE_NAVIGATION",
            Self::Science => "MISSION_TYPE_SCIENCE",
            Self::Weather => "MISSION_TYPE_WEATHER",
            Self::Amateur => "MISSION_TYPE_AMATEUR",
            Self::TechDemo => "MISSION_TYPE_TECH_DEMO",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for MissionType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "MISSION_TYPE_UNSPECIFIED",
            "MISSION_TYPE_COMMUNICATIONS",
            "MISSION_TYPE_EARTH_OBSERVATION",
            "MISSION_TYPE_NAVIGATION",
            "MISSION_TYPE_SCIENCE",
            "MISSION_TYPE_WEATHER",
            "MISSION_TYPE_AMATEUR",
            "MISSION_TYPE_TECH_DEMO",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MissionType;

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
                    "MISSION_TYPE_UNSPECIFIED" => Ok(MissionType::Unspecified),
                    "MISSION_TYPE_COMMUNICATIONS" => Ok(MissionType::Communications),
                    "MISSION_TYPE_EARTH_OBSERVATION" => Ok(MissionType::EarthObservation),
                    "MISSION_TYPE_NAVIGATION" => Ok(MissionType::Navigation),
                    "MISSION_TYPE_SCIENCE" => Ok(MissionType::Science),
                    "MISSION_TYPE_WEATHER" => Ok(MissionType::Weather),
                    "MISSION_TYPE_AMATEUR" => Ok(MissionType::Amateur),
                    "MISSION_TYPE_TECH_DEMO" => Ok(MissionType::TechDemo),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for ObjectType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "OBJECT_TYPE_UNSPECIFIED",
            Self::Payload => "OBJECT_TYPE_PAYLOAD",
            Self::RocketBody => "OBJECT_TYPE_ROCKET_BODY",
            Self::Debris => "OBJECT_TYPE_DEBRIS",
            Self::Unknown => "OBJECT_TYPE_UNKNOWN",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ObjectType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "OBJECT_TYPE_UNSPECIFIED",
            "OBJECT_TYPE_PAYLOAD",
            "OBJECT_TYPE_ROCKET_BODY",
            "OBJECT_TYPE_DEBRIS",
            "OBJECT_TYPE_UNKNOWN",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ObjectType;

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
                    "OBJECT_TYPE_UNSPECIFIED" => Ok(ObjectType::Unspecified),
                    "OBJECT_TYPE_PAYLOAD" => Ok(ObjectType::Payload),
                    "OBJECT_TYPE_ROCKET_BODY" => Ok(ObjectType::RocketBody),
                    "OBJECT_TYPE_DEBRIS" => Ok(ObjectType::Debris),
                    "OBJECT_TYPE_UNKNOWN" => Ok(ObjectType::Unknown),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for OperationalStatus {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "OPERATIONAL_STATUS_UNSPECIFIED",
            Self::Active => "OPERATIONAL_STATUS_ACTIVE",
            Self::Inactive => "OPERATIONAL_STATUS_INACTIVE",
            Self::Decayed => "OPERATIONAL_STATUS_DECAYED",
            Self::Unknown => "OPERATIONAL_STATUS_UNKNOWN",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for OperationalStatus {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "OPERATIONAL_STATUS_UNSPECIFIED",
            "OPERATIONAL_STATUS_ACTIVE",
            "OPERATIONAL_STATUS_INACTIVE",
            "OPERATIONAL_STATUS_DECAYED",
            "OPERATIONAL_STATUS_UNKNOWN",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = OperationalStatus;

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
                    "OPERATIONAL_STATUS_UNSPECIFIED" => Ok(OperationalStatus::Unspecified),
                    "OPERATIONAL_STATUS_ACTIVE" => Ok(OperationalStatus::Active),
                    "OPERATIONAL_STATUS_INACTIVE" => Ok(OperationalStatus::Inactive),
                    "OPERATIONAL_STATUS_DECAYED" => Ok(OperationalStatus::Decayed),
                    "OPERATIONAL_STATUS_UNKNOWN" => Ok(OperationalStatus::Unknown),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for OrbitRegime {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "ORBIT_REGIME_UNSPECIFIED",
            Self::Leo => "ORBIT_REGIME_LEO",
            Self::Meo => "ORBIT_REGIME_MEO",
            Self::Geo => "ORBIT_REGIME_GEO",
            Self::Heo => "ORBIT_REGIME_HEO",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for OrbitRegime {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ORBIT_REGIME_UNSPECIFIED",
            "ORBIT_REGIME_LEO",
            "ORBIT_REGIME_MEO",
            "ORBIT_REGIME_GEO",
            "ORBIT_REGIME_HEO",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = OrbitRegime;

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
                    "ORBIT_REGIME_UNSPECIFIED" => Ok(OrbitRegime::Unspecified),
                    "ORBIT_REGIME_LEO" => Ok(OrbitRegime::Leo),
                    "ORBIT_REGIME_MEO" => Ok(OrbitRegime::Meo),
                    "ORBIT_REGIME_GEO" => Ok(OrbitRegime::Geo),
                    "ORBIT_REGIME_HEO" => Ok(OrbitRegime::Heo),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for SatelliteMetadata {
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
        if self.cospar_id.is_some() {
            len += 1;
        }
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.aliases.is_empty() {
            len += 1;
        }
        if self.object_type != 0 {
            len += 1;
        }
        if self.mission_type != 0 {
            len += 1;
        }
        if self.orbit_regime != 0 {
            len += 1;
        }
        if self.operator.is_some() {
            len += 1;
        }
        if self.owner.is_some() {
            len += 1;
        }
        if self.constellation.is_some() {
            len += 1;
        }
        if self.launch_date.is_some() {
            len += 1;
        }
        if self.launch_site.is_some() {
            len += 1;
        }
        if self.launch_vehicle.is_some() {
            len += 1;
        }
        if self.operational_status != 0 {
            len += 1;
        }
        if !self.frequencies.is_empty() {
            len += 1;
        }
        if !self.sources.is_empty() {
            len += 1;
        }
        if self.updated_at.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metadata.SatelliteMetadata", len)?;
        if self.norad_id != 0 {
            struct_ser.serialize_field("noradId", &self.norad_id)?;
        }
        if let Some(v) = self.cospar_id.as_ref() {
            struct_ser.serialize_field("cosparId", v)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.aliases.is_empty() {
            struct_ser.serialize_field("aliases", &self.aliases)?;
        }
        if self.object_type != 0 {
            let v = ObjectType::try_from(self.object_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.object_type)))?;
            struct_ser.serialize_field("objectType", &v)?;
        }
        if self.mission_type != 0 {
            let v = MissionType::try_from(self.mission_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.mission_type)))?;
            struct_ser.serialize_field("missionType", &v)?;
        }
        if self.orbit_regime != 0 {
            let v = OrbitRegime::try_from(self.orbit_regime)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.orbit_regime)))?;
            struct_ser.serialize_field("orbitRegime", &v)?;
        }
        if let Some(v) = self.operator.as_ref() {
            struct_ser.serialize_field("operator", v)?;
        }
        if let Some(v) = self.owner.as_ref() {
            struct_ser.serialize_field("owner", v)?;
        }
        if let Some(v) = self.constellation.as_ref() {
            struct_ser.serialize_field("constellation", v)?;
        }
        if let Some(v) = self.launch_date.as_ref() {
            struct_ser.serialize_field("launchDate", v)?;
        }
        if let Some(v) = self.launch_site.as_ref() {
            struct_ser.serialize_field("launchSite", v)?;
        }
        if let Some(v) = self.launch_vehicle.as_ref() {
            struct_ser.serialize_field("launchVehicle", v)?;
        }
        if self.operational_status != 0 {
            let v = OperationalStatus::try_from(self.operational_status)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.operational_status)))?;
            struct_ser.serialize_field("operationalStatus", &v)?;
        }
        if !self.frequencies.is_empty() {
            struct_ser.serialize_field("frequencies", &self.frequencies)?;
        }
        if !self.sources.is_empty() {
            struct_ser.serialize_field("sources", &self.sources)?;
        }
        if let Some(v) = self.updated_at.as_ref() {
            struct_ser.serialize_field("updatedAt", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SatelliteMetadata {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "norad_id",
            "noradId",
            "cospar_id",
            "cosparId",
            "name",
            "aliases",
            "object_type",
            "objectType",
            "mission_type",
            "missionType",
            "orbit_regime",
            "orbitRegime",
            "operator",
            "owner",
            "constellation",
            "launch_date",
            "launchDate",
            "launch_site",
            "launchSite",
            "launch_vehicle",
            "launchVehicle",
            "operational_status",
            "operationalStatus",
            "frequencies",
            "sources",
            "updated_at",
            "updatedAt",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            NoradId,
            CosparId,
            Name,
            Aliases,
            ObjectType,
            MissionType,
            OrbitRegime,
            Operator,
            Owner,
            Constellation,
            LaunchDate,
            LaunchSite,
            LaunchVehicle,
            OperationalStatus,
            Frequencies,
            Sources,
            UpdatedAt,
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
                            "cosparId" | "cospar_id" => Ok(GeneratedField::CosparId),
                            "name" => Ok(GeneratedField::Name),
                            "aliases" => Ok(GeneratedField::Aliases),
                            "objectType" | "object_type" => Ok(GeneratedField::ObjectType),
                            "missionType" | "mission_type" => Ok(GeneratedField::MissionType),
                            "orbitRegime" | "orbit_regime" => Ok(GeneratedField::OrbitRegime),
                            "operator" => Ok(GeneratedField::Operator),
                            "owner" => Ok(GeneratedField::Owner),
                            "constellation" => Ok(GeneratedField::Constellation),
                            "launchDate" | "launch_date" => Ok(GeneratedField::LaunchDate),
                            "launchSite" | "launch_site" => Ok(GeneratedField::LaunchSite),
                            "launchVehicle" | "launch_vehicle" => Ok(GeneratedField::LaunchVehicle),
                            "operationalStatus" | "operational_status" => Ok(GeneratedField::OperationalStatus),
                            "frequencies" => Ok(GeneratedField::Frequencies),
                            "sources" => Ok(GeneratedField::Sources),
                            "updatedAt" | "updated_at" => Ok(GeneratedField::UpdatedAt),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SatelliteMetadata;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.SatelliteMetadata")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SatelliteMetadata, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut norad_id__ = None;
                let mut cospar_id__ = None;
                let mut name__ = None;
                let mut aliases__ = None;
                let mut object_type__ = None;
                let mut mission_type__ = None;
                let mut orbit_regime__ = None;
                let mut operator__ = None;
                let mut owner__ = None;
                let mut constellation__ = None;
                let mut launch_date__ = None;
                let mut launch_site__ = None;
                let mut launch_vehicle__ = None;
                let mut operational_status__ = None;
                let mut frequencies__ = None;
                let mut sources__ = None;
                let mut updated_at__ = None;
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
                        GeneratedField::CosparId => {
                            if cospar_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cosparId"));
                            }
                            cospar_id__ = map_.next_value()?;
                        }
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Aliases => {
                            if aliases__.is_some() {
                                return Err(serde::de::Error::duplicate_field("aliases"));
                            }
                            aliases__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ObjectType => {
                            if object_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("objectType"));
                            }
                            object_type__ = Some(map_.next_value::<ObjectType>()? as i32);
                        }
                        GeneratedField::MissionType => {
                            if mission_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("missionType"));
                            }
                            mission_type__ = Some(map_.next_value::<MissionType>()? as i32);
                        }
                        GeneratedField::OrbitRegime => {
                            if orbit_regime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("orbitRegime"));
                            }
                            orbit_regime__ = Some(map_.next_value::<OrbitRegime>()? as i32);
                        }
                        GeneratedField::Operator => {
                            if operator__.is_some() {
                                return Err(serde::de::Error::duplicate_field("operator"));
                            }
                            operator__ = map_.next_value()?;
                        }
                        GeneratedField::Owner => {
                            if owner__.is_some() {
                                return Err(serde::de::Error::duplicate_field("owner"));
                            }
                            owner__ = map_.next_value()?;
                        }
                        GeneratedField::Constellation => {
                            if constellation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("constellation"));
                            }
                            constellation__ = map_.next_value()?;
                        }
                        GeneratedField::LaunchDate => {
                            if launch_date__.is_some() {
                                return Err(serde::de::Error::duplicate_field("launchDate"));
                            }
                            launch_date__ = map_.next_value()?;
                        }
                        GeneratedField::LaunchSite => {
                            if launch_site__.is_some() {
                                return Err(serde::de::Error::duplicate_field("launchSite"));
                            }
                            launch_site__ = map_.next_value()?;
                        }
                        GeneratedField::LaunchVehicle => {
                            if launch_vehicle__.is_some() {
                                return Err(serde::de::Error::duplicate_field("launchVehicle"));
                            }
                            launch_vehicle__ = map_.next_value()?;
                        }
                        GeneratedField::OperationalStatus => {
                            if operational_status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("operationalStatus"));
                            }
                            operational_status__ = Some(map_.next_value::<OperationalStatus>()? as i32);
                        }
                        GeneratedField::Frequencies => {
                            if frequencies__.is_some() {
                                return Err(serde::de::Error::duplicate_field("frequencies"));
                            }
                            frequencies__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Sources => {
                            if sources__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sources"));
                            }
                            sources__ = Some(map_.next_value()?);
                        }
                        GeneratedField::UpdatedAt => {
                            if updated_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("updatedAt"));
                            }
                            updated_at__ = map_.next_value()?;
                        }
                    }
                }
                Ok(SatelliteMetadata {
                    norad_id: norad_id__.unwrap_or_default(),
                    cospar_id: cospar_id__,
                    name: name__.unwrap_or_default(),
                    aliases: aliases__.unwrap_or_default(),
                    object_type: object_type__.unwrap_or_default(),
                    mission_type: mission_type__.unwrap_or_default(),
                    orbit_regime: orbit_regime__.unwrap_or_default(),
                    operator: operator__,
                    owner: owner__,
                    constellation: constellation__,
                    launch_date: launch_date__,
                    launch_site: launch_site__,
                    launch_vehicle: launch_vehicle__,
                    operational_status: operational_status__.unwrap_or_default(),
                    frequencies: frequencies__.unwrap_or_default(),
                    sources: sources__.unwrap_or_default(),
                    updated_at: updated_at__,
                })
            }
        }
        deserializer.deserialize_struct("metadata.SatelliteMetadata", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Source {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "SOURCE_UNSPECIFIED",
            Self::Celestrak => "SOURCE_CELESTRAK",
            Self::Amsat => "SOURCE_AMSAT",
            Self::Ucs => "SOURCE_UCS",
            Self::Manual => "SOURCE_MANUAL",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for Source {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "SOURCE_UNSPECIFIED",
            "SOURCE_CELESTRAK",
            "SOURCE_AMSAT",
            "SOURCE_UCS",
            "SOURCE_MANUAL",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Source;

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
                    "SOURCE_UNSPECIFIED" => Ok(Source::Unspecified),
                    "SOURCE_CELESTRAK" => Ok(Source::Celestrak),
                    "SOURCE_AMSAT" => Ok(Source::Amsat),
                    "SOURCE_UCS" => Ok(Source::Ucs),
                    "SOURCE_MANUAL" => Ok(Source::Manual),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for SourceAttribution {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.source != 0 {
            len += 1;
        }
        if !self.source_record_id.is_empty() {
            len += 1;
        }
        if self.fetched_at.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metadata.SourceAttribution", len)?;
        if self.source != 0 {
            let v = Source::try_from(self.source)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.source)))?;
            struct_ser.serialize_field("source", &v)?;
        }
        if !self.source_record_id.is_empty() {
            struct_ser.serialize_field("sourceRecordId", &self.source_record_id)?;
        }
        if let Some(v) = self.fetched_at.as_ref() {
            struct_ser.serialize_field("fetchedAt", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SourceAttribution {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "source",
            "source_record_id",
            "sourceRecordId",
            "fetched_at",
            "fetchedAt",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Source,
            SourceRecordId,
            FetchedAt,
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
                            "source" => Ok(GeneratedField::Source),
                            "sourceRecordId" | "source_record_id" => Ok(GeneratedField::SourceRecordId),
                            "fetchedAt" | "fetched_at" => Ok(GeneratedField::FetchedAt),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SourceAttribution;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metadata.SourceAttribution")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SourceAttribution, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut source__ = None;
                let mut source_record_id__ = None;
                let mut fetched_at__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Source => {
                            if source__.is_some() {
                                return Err(serde::de::Error::duplicate_field("source"));
                            }
                            source__ = Some(map_.next_value::<Source>()? as i32);
                        }
                        GeneratedField::SourceRecordId => {
                            if source_record_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sourceRecordId"));
                            }
                            source_record_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::FetchedAt => {
                            if fetched_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fetchedAt"));
                            }
                            fetched_at__ = map_.next_value()?;
                        }
                    }
                }
                Ok(SourceAttribution {
                    source: source__.unwrap_or_default(),
                    source_record_id: source_record_id__.unwrap_or_default(),
                    fetched_at: fetched_at__,
                })
            }
        }
        deserializer.deserialize_struct("metadata.SourceAttribution", FIELDS, GeneratedVisitor)
    }
}
