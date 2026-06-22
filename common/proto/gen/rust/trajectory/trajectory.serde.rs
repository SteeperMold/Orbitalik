// @generated
impl serde::Serialize for LookAnglesRequest {
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
        if self.datetime.is_some() {
            len += 1;
        }
        if self.observer.is_some() {
            len += 1;
        }
        if self.output_mask.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.LookAnglesRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        if let Some(v) = self.datetime.as_ref() {
            struct_ser.serialize_field("datetime", v)?;
        }
        if let Some(v) = self.observer.as_ref() {
            struct_ser.serialize_field("observer", v)?;
        }
        if let Some(v) = self.output_mask.as_ref() {
            struct_ser.serialize_field("outputMask", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for LookAnglesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "identifier",
            "datetime",
            "observer",
            "output_mask",
            "outputMask",
            "units",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Identifier,
            Datetime,
            Observer,
            OutputMask,
            Units,
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
                            "datetime" => Ok(GeneratedField::Datetime),
                            "observer" => Ok(GeneratedField::Observer),
                            "outputMask" | "output_mask" => Ok(GeneratedField::OutputMask),
                            "units" => Ok(GeneratedField::Units),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = LookAnglesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.LookAnglesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<LookAnglesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut identifier__ = None;
                let mut datetime__ = None;
                let mut observer__ = None;
                let mut output_mask__ = None;
                let mut units__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = map_.next_value()?;
                        }
                        GeneratedField::Datetime => {
                            if datetime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("datetime"));
                            }
                            datetime__ = map_.next_value()?;
                        }
                        GeneratedField::Observer => {
                            if observer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("observer"));
                            }
                            observer__ = map_.next_value()?;
                        }
                        GeneratedField::OutputMask => {
                            if output_mask__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outputMask"));
                            }
                            output_mask__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                    }
                }
                Ok(LookAnglesRequest {
                    identifier: identifier__,
                    datetime: datetime__,
                    observer: observer__,
                    output_mask: output_mask__,
                    units: units__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.LookAnglesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for LookAnglesResponse {
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
        if self.time.is_some() {
            len += 1;
        }
        if self.azimuth.is_some() {
            len += 1;
        }
        if self.elevation.is_some() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.LookAnglesResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        if let Some(v) = self.time.as_ref() {
            struct_ser.serialize_field("time", v)?;
        }
        if let Some(v) = self.azimuth.as_ref() {
            struct_ser.serialize_field("azimuth", v)?;
        }
        if let Some(v) = self.elevation.as_ref() {
            struct_ser.serialize_field("elevation", v)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for LookAnglesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
            "time",
            "azimuth",
            "elevation",
            "range",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
            Time,
            Azimuth,
            Elevation,
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
                            "time" => Ok(GeneratedField::Time),
                            "azimuth" => Ok(GeneratedField::Azimuth),
                            "elevation" => Ok(GeneratedField::Elevation),
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
            type Value = LookAnglesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.LookAnglesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<LookAnglesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                let mut time__ = None;
                let mut azimuth__ = None;
                let mut elevation__ = None;
                let mut range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                        GeneratedField::Time => {
                            if time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("time"));
                            }
                            time__ = map_.next_value()?;
                        }
                        GeneratedField::Azimuth => {
                            if azimuth__.is_some() {
                                return Err(serde::de::Error::duplicate_field("azimuth"));
                            }
                            azimuth__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Elevation => {
                            if elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("elevation"));
                            }
                            elevation__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(LookAnglesResponse {
                    metadata: metadata__,
                    time: time__,
                    azimuth: azimuth__,
                    elevation: elevation__,
                    range: range__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.LookAnglesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ObserverTrajectoryPoint {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.datetime.is_some() {
            len += 1;
        }
        if self.azimuth.is_some() {
            len += 1;
        }
        if self.elevation.is_some() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.ObserverTrajectoryPoint", len)?;
        if let Some(v) = self.datetime.as_ref() {
            struct_ser.serialize_field("datetime", v)?;
        }
        if let Some(v) = self.azimuth.as_ref() {
            struct_ser.serialize_field("azimuth", v)?;
        }
        if let Some(v) = self.elevation.as_ref() {
            struct_ser.serialize_field("elevation", v)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ObserverTrajectoryPoint {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "datetime",
            "azimuth",
            "elevation",
            "range",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Datetime,
            Azimuth,
            Elevation,
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
                            "datetime" => Ok(GeneratedField::Datetime),
                            "azimuth" => Ok(GeneratedField::Azimuth),
                            "elevation" => Ok(GeneratedField::Elevation),
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
            type Value = ObserverTrajectoryPoint;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.ObserverTrajectoryPoint")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ObserverTrajectoryPoint, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut datetime__ = None;
                let mut azimuth__ = None;
                let mut elevation__ = None;
                let mut range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Datetime => {
                            if datetime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("datetime"));
                            }
                            datetime__ = map_.next_value()?;
                        }
                        GeneratedField::Azimuth => {
                            if azimuth__.is_some() {
                                return Err(serde::de::Error::duplicate_field("azimuth"));
                            }
                            azimuth__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Elevation => {
                            if elevation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("elevation"));
                            }
                            elevation__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::NumberDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(ObserverTrajectoryPoint {
                    datetime: datetime__,
                    azimuth: azimuth__,
                    elevation: elevation__,
                    range: range__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.ObserverTrajectoryPoint", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ObserverTrajectoryRequest {
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
        if self.range.is_some() {
            len += 1;
        }
        if self.sampling.is_some() {
            len += 1;
        }
        if self.observer.is_some() {
            len += 1;
        }
        if self.output_mask.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.ObserverTrajectoryRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        if let Some(v) = self.sampling.as_ref() {
            struct_ser.serialize_field("sampling", v)?;
        }
        if let Some(v) = self.observer.as_ref() {
            struct_ser.serialize_field("observer", v)?;
        }
        if let Some(v) = self.output_mask.as_ref() {
            struct_ser.serialize_field("outputMask", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ObserverTrajectoryRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "identifier",
            "range",
            "sampling",
            "observer",
            "output_mask",
            "outputMask",
            "units",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Identifier,
            Range,
            Sampling,
            Observer,
            OutputMask,
            Units,
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
                            "range" => Ok(GeneratedField::Range),
                            "sampling" => Ok(GeneratedField::Sampling),
                            "observer" => Ok(GeneratedField::Observer),
                            "outputMask" | "output_mask" => Ok(GeneratedField::OutputMask),
                            "units" => Ok(GeneratedField::Units),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ObserverTrajectoryRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.ObserverTrajectoryRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ObserverTrajectoryRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut identifier__ = None;
                let mut range__ = None;
                let mut sampling__ = None;
                let mut observer__ = None;
                let mut output_mask__ = None;
                let mut units__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = map_.next_value()?;
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                        GeneratedField::Sampling => {
                            if sampling__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sampling"));
                            }
                            sampling__ = map_.next_value()?;
                        }
                        GeneratedField::Observer => {
                            if observer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("observer"));
                            }
                            observer__ = map_.next_value()?;
                        }
                        GeneratedField::OutputMask => {
                            if output_mask__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outputMask"));
                            }
                            output_mask__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                    }
                }
                Ok(ObserverTrajectoryRequest {
                    identifier: identifier__,
                    range: range__,
                    sampling: sampling__,
                    observer: observer__,
                    output_mask: output_mask__,
                    units: units__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.ObserverTrajectoryRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ObserverTrajectoryResponse {
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
        if !self.points.is_empty() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        if self.sampling.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.ObserverTrajectoryResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        if !self.points.is_empty() {
            struct_ser.serialize_field("points", &self.points)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        if let Some(v) = self.sampling.as_ref() {
            struct_ser.serialize_field("sampling", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ObserverTrajectoryResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
            "points",
            "range",
            "sampling",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
            Points,
            Range,
            Sampling,
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
                            "points" => Ok(GeneratedField::Points),
                            "range" => Ok(GeneratedField::Range),
                            "sampling" => Ok(GeneratedField::Sampling),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ObserverTrajectoryResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.ObserverTrajectoryResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ObserverTrajectoryResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                let mut points__ = None;
                let mut range__ = None;
                let mut sampling__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                        GeneratedField::Points => {
                            if points__.is_some() {
                                return Err(serde::de::Error::duplicate_field("points"));
                            }
                            points__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                        GeneratedField::Sampling => {
                            if sampling__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sampling"));
                            }
                            sampling__ = map_.next_value()?;
                        }
                    }
                }
                Ok(ObserverTrajectoryResponse {
                    metadata: metadata__,
                    points: points__.unwrap_or_default(),
                    range: range__,
                    sampling: sampling__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.ObserverTrajectoryResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PositionRequest {
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
        if self.datetime.is_some() {
            len += 1;
        }
        if self.output_mask.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.PositionRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        if let Some(v) = self.datetime.as_ref() {
            struct_ser.serialize_field("datetime", v)?;
        }
        if let Some(v) = self.output_mask.as_ref() {
            struct_ser.serialize_field("outputMask", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PositionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "identifier",
            "datetime",
            "output_mask",
            "outputMask",
            "units",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Identifier,
            Datetime,
            OutputMask,
            Units,
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
                            "datetime" => Ok(GeneratedField::Datetime),
                            "outputMask" | "output_mask" => Ok(GeneratedField::OutputMask),
                            "units" => Ok(GeneratedField::Units),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PositionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.PositionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PositionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut identifier__ = None;
                let mut datetime__ = None;
                let mut output_mask__ = None;
                let mut units__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = map_.next_value()?;
                        }
                        GeneratedField::Datetime => {
                            if datetime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("datetime"));
                            }
                            datetime__ = map_.next_value()?;
                        }
                        GeneratedField::OutputMask => {
                            if output_mask__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outputMask"));
                            }
                            output_mask__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                    }
                }
                Ok(PositionRequest {
                    identifier: identifier__,
                    datetime: datetime__,
                    output_mask: output_mask__,
                    units: units__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.PositionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PositionResponse {
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
        if self.time.is_some() {
            len += 1;
        }
        if self.eci.is_some() {
            len += 1;
        }
        if self.ecef.is_some() {
            len += 1;
        }
        if self.geodetic.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.PositionResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        if let Some(v) = self.time.as_ref() {
            struct_ser.serialize_field("time", v)?;
        }
        if let Some(v) = self.eci.as_ref() {
            struct_ser.serialize_field("eci", v)?;
        }
        if let Some(v) = self.ecef.as_ref() {
            struct_ser.serialize_field("ecef", v)?;
        }
        if let Some(v) = self.geodetic.as_ref() {
            struct_ser.serialize_field("geodetic", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PositionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
            "time",
            "eci",
            "ecef",
            "geodetic",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
            Time,
            Eci,
            Ecef,
            Geodetic,
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
                            "time" => Ok(GeneratedField::Time),
                            "eci" => Ok(GeneratedField::Eci),
                            "ecef" => Ok(GeneratedField::Ecef),
                            "geodetic" => Ok(GeneratedField::Geodetic),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PositionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.PositionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PositionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                let mut time__ = None;
                let mut eci__ = None;
                let mut ecef__ = None;
                let mut geodetic__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                        GeneratedField::Time => {
                            if time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("time"));
                            }
                            time__ = map_.next_value()?;
                        }
                        GeneratedField::Eci => {
                            if eci__.is_some() {
                                return Err(serde::de::Error::duplicate_field("eci"));
                            }
                            eci__ = map_.next_value()?;
                        }
                        GeneratedField::Ecef => {
                            if ecef__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ecef"));
                            }
                            ecef__ = map_.next_value()?;
                        }
                        GeneratedField::Geodetic => {
                            if geodetic__.is_some() {
                                return Err(serde::de::Error::duplicate_field("geodetic"));
                            }
                            geodetic__ = map_.next_value()?;
                        }
                    }
                }
                Ok(PositionResponse {
                    metadata: metadata__,
                    time: time__,
                    eci: eci__,
                    ecef: ecef__,
                    geodetic: geodetic__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.PositionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for StateVector {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.datetime.is_some() {
            len += 1;
        }
        if self.position_eci.is_some() {
            len += 1;
        }
        if self.velocity_eci.is_some() {
            len += 1;
        }
        if self.position_ecef.is_some() {
            len += 1;
        }
        if self.velocity_ecef.is_some() {
            len += 1;
        }
        if self.geodetic.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.StateVector", len)?;
        if let Some(v) = self.datetime.as_ref() {
            struct_ser.serialize_field("datetime", v)?;
        }
        if let Some(v) = self.position_eci.as_ref() {
            struct_ser.serialize_field("positionEci", v)?;
        }
        if let Some(v) = self.velocity_eci.as_ref() {
            struct_ser.serialize_field("velocityEci", v)?;
        }
        if let Some(v) = self.position_ecef.as_ref() {
            struct_ser.serialize_field("positionEcef", v)?;
        }
        if let Some(v) = self.velocity_ecef.as_ref() {
            struct_ser.serialize_field("velocityEcef", v)?;
        }
        if let Some(v) = self.geodetic.as_ref() {
            struct_ser.serialize_field("geodetic", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for StateVector {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "datetime",
            "position_eci",
            "positionEci",
            "velocity_eci",
            "velocityEci",
            "position_ecef",
            "positionEcef",
            "velocity_ecef",
            "velocityEcef",
            "geodetic",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Datetime,
            PositionEci,
            VelocityEci,
            PositionEcef,
            VelocityEcef,
            Geodetic,
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
                            "datetime" => Ok(GeneratedField::Datetime),
                            "positionEci" | "position_eci" => Ok(GeneratedField::PositionEci),
                            "velocityEci" | "velocity_eci" => Ok(GeneratedField::VelocityEci),
                            "positionEcef" | "position_ecef" => Ok(GeneratedField::PositionEcef),
                            "velocityEcef" | "velocity_ecef" => Ok(GeneratedField::VelocityEcef),
                            "geodetic" => Ok(GeneratedField::Geodetic),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = StateVector;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.StateVector")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<StateVector, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut datetime__ = None;
                let mut position_eci__ = None;
                let mut velocity_eci__ = None;
                let mut position_ecef__ = None;
                let mut velocity_ecef__ = None;
                let mut geodetic__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Datetime => {
                            if datetime__.is_some() {
                                return Err(serde::de::Error::duplicate_field("datetime"));
                            }
                            datetime__ = map_.next_value()?;
                        }
                        GeneratedField::PositionEci => {
                            if position_eci__.is_some() {
                                return Err(serde::de::Error::duplicate_field("positionEci"));
                            }
                            position_eci__ = map_.next_value()?;
                        }
                        GeneratedField::VelocityEci => {
                            if velocity_eci__.is_some() {
                                return Err(serde::de::Error::duplicate_field("velocityEci"));
                            }
                            velocity_eci__ = map_.next_value()?;
                        }
                        GeneratedField::PositionEcef => {
                            if position_ecef__.is_some() {
                                return Err(serde::de::Error::duplicate_field("positionEcef"));
                            }
                            position_ecef__ = map_.next_value()?;
                        }
                        GeneratedField::VelocityEcef => {
                            if velocity_ecef__.is_some() {
                                return Err(serde::de::Error::duplicate_field("velocityEcef"));
                            }
                            velocity_ecef__ = map_.next_value()?;
                        }
                        GeneratedField::Geodetic => {
                            if geodetic__.is_some() {
                                return Err(serde::de::Error::duplicate_field("geodetic"));
                            }
                            geodetic__ = map_.next_value()?;
                        }
                    }
                }
                Ok(StateVector {
                    datetime: datetime__,
                    position_eci: position_eci__,
                    velocity_eci: velocity_eci__,
                    position_ecef: position_ecef__,
                    velocity_ecef: velocity_ecef__,
                    geodetic: geodetic__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.StateVector", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TrajectoryComputationMetadata {
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
        if self.norad_id != 0 {
            len += 1;
        }
        if !self.satellite_name.is_empty() {
            len += 1;
        }
        if self.tle_epoch.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.TrajectoryComputationMetadata", len)?;
        if !self.propagation_model.is_empty() {
            struct_ser.serialize_field("propagationModel", &self.propagation_model)?;
        }
        if let Some(v) = self.computation_time.as_ref() {
            struct_ser.serialize_field("computationTime", v)?;
        }
        if self.norad_id != 0 {
            struct_ser.serialize_field("noradId", &self.norad_id)?;
        }
        if !self.satellite_name.is_empty() {
            struct_ser.serialize_field("satelliteName", &self.satellite_name)?;
        }
        if let Some(v) = self.tle_epoch.as_ref() {
            struct_ser.serialize_field("tleEpoch", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TrajectoryComputationMetadata {
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
            "norad_id",
            "noradId",
            "satellite_name",
            "satelliteName",
            "tle_epoch",
            "tleEpoch",
            "units",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PropagationModel,
            ComputationTime,
            NoradId,
            SatelliteName,
            TleEpoch,
            Units,
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
                            "noradId" | "norad_id" => Ok(GeneratedField::NoradId),
                            "satelliteName" | "satellite_name" => Ok(GeneratedField::SatelliteName),
                            "tleEpoch" | "tle_epoch" => Ok(GeneratedField::TleEpoch),
                            "units" => Ok(GeneratedField::Units),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TrajectoryComputationMetadata;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.TrajectoryComputationMetadata")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TrajectoryComputationMetadata, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut propagation_model__ = None;
                let mut computation_time__ = None;
                let mut norad_id__ = None;
                let mut satellite_name__ = None;
                let mut tle_epoch__ = None;
                let mut units__ = None;
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
                    }
                }
                Ok(TrajectoryComputationMetadata {
                    propagation_model: propagation_model__.unwrap_or_default(),
                    computation_time: computation_time__,
                    norad_id: norad_id__.unwrap_or_default(),
                    satellite_name: satellite_name__.unwrap_or_default(),
                    tle_epoch: tle_epoch__,
                    units: units__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.TrajectoryComputationMetadata", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TrajectoryRequest {
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
        if self.range.is_some() {
            len += 1;
        }
        if self.sampling.is_some() {
            len += 1;
        }
        if self.output_mask.is_some() {
            len += 1;
        }
        if self.units.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.TrajectoryRequest", len)?;
        if let Some(v) = self.identifier.as_ref() {
            struct_ser.serialize_field("identifier", v)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        if let Some(v) = self.sampling.as_ref() {
            struct_ser.serialize_field("sampling", v)?;
        }
        if let Some(v) = self.output_mask.as_ref() {
            struct_ser.serialize_field("outputMask", v)?;
        }
        if let Some(v) = self.units.as_ref() {
            struct_ser.serialize_field("units", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TrajectoryRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "identifier",
            "range",
            "sampling",
            "output_mask",
            "outputMask",
            "units",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Identifier,
            Range,
            Sampling,
            OutputMask,
            Units,
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
                            "range" => Ok(GeneratedField::Range),
                            "sampling" => Ok(GeneratedField::Sampling),
                            "outputMask" | "output_mask" => Ok(GeneratedField::OutputMask),
                            "units" => Ok(GeneratedField::Units),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TrajectoryRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.TrajectoryRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TrajectoryRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut identifier__ = None;
                let mut range__ = None;
                let mut sampling__ = None;
                let mut output_mask__ = None;
                let mut units__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Identifier => {
                            if identifier__.is_some() {
                                return Err(serde::de::Error::duplicate_field("identifier"));
                            }
                            identifier__ = map_.next_value()?;
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                        GeneratedField::Sampling => {
                            if sampling__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sampling"));
                            }
                            sampling__ = map_.next_value()?;
                        }
                        GeneratedField::OutputMask => {
                            if output_mask__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outputMask"));
                            }
                            output_mask__ = map_.next_value()?;
                        }
                        GeneratedField::Units => {
                            if units__.is_some() {
                                return Err(serde::de::Error::duplicate_field("units"));
                            }
                            units__ = map_.next_value()?;
                        }
                    }
                }
                Ok(TrajectoryRequest {
                    identifier: identifier__,
                    range: range__,
                    sampling: sampling__,
                    output_mask: output_mask__,
                    units: units__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.TrajectoryRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TrajectoryResponse {
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
        if !self.states.is_empty() {
            len += 1;
        }
        if self.range.is_some() {
            len += 1;
        }
        if self.sampling.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("trajectory.TrajectoryResponse", len)?;
        if let Some(v) = self.metadata.as_ref() {
            struct_ser.serialize_field("metadata", v)?;
        }
        if !self.states.is_empty() {
            struct_ser.serialize_field("states", &self.states)?;
        }
        if let Some(v) = self.range.as_ref() {
            struct_ser.serialize_field("range", v)?;
        }
        if let Some(v) = self.sampling.as_ref() {
            struct_ser.serialize_field("sampling", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TrajectoryResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "metadata",
            "states",
            "range",
            "sampling",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Metadata,
            States,
            Range,
            Sampling,
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
                            "states" => Ok(GeneratedField::States),
                            "range" => Ok(GeneratedField::Range),
                            "sampling" => Ok(GeneratedField::Sampling),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TrajectoryResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct trajectory.TrajectoryResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TrajectoryResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut metadata__ = None;
                let mut states__ = None;
                let mut range__ = None;
                let mut sampling__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = map_.next_value()?;
                        }
                        GeneratedField::States => {
                            if states__.is_some() {
                                return Err(serde::de::Error::duplicate_field("states"));
                            }
                            states__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Range => {
                            if range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("range"));
                            }
                            range__ = map_.next_value()?;
                        }
                        GeneratedField::Sampling => {
                            if sampling__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sampling"));
                            }
                            sampling__ = map_.next_value()?;
                        }
                    }
                }
                Ok(TrajectoryResponse {
                    metadata: metadata__,
                    states: states__.unwrap_or_default(),
                    range: range__,
                    sampling: sampling__,
                })
            }
        }
        deserializer.deserialize_struct("trajectory.TrajectoryResponse", FIELDS, GeneratedVisitor)
    }
}
