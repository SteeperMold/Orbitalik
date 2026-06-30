import dataclasses
import enum
from datetime import datetime


@dataclasses.dataclass(frozen=True)
class SatelliteIdentifier:
    norad_id: int | None = None
    satellite_name: str | None = None

    def __post_init__(self) -> None:
        if (self.norad_id is None) == (self.satellite_name is None):
            raise ValueError(
                "SatelliteIdentifier must have exactly one of norad_id or satellite_name"
            )


@dataclasses.dataclass(frozen=True)
class GeodeticInput:
    lat_deg: float | None = None
    lat_rad: float | None = None

    lon_deg: float | None = None
    lon_rad: float | None = None

    alt_m: float | None = None
    alt_km: float | None = None

    def __post_init__(self) -> None:
        if (self.lat_deg is None) == (self.lat_rad is None):
            raise ValueError("Exactly one of lat_deg or lat_rad must be set")

        if (self.lon_deg is None) == (self.lon_rad is None):
            raise ValueError("Exactly one of lon_deg or lon_rad must be set")

        if (self.alt_m is None) == (self.alt_km is None):
            raise ValueError("Exactly one of alt_m or alt_km must be set")


@dataclasses.dataclass
class Subscription:
    id: int
    user_id: int

    satellite: SatelliteIdentifier
    observer: GeodeticInput

    enabled: bool

    notify_before_seconds: int

    min_peak_elevation_deg: float | None = None
    min_peak_elevation_rad: float | None = None

    min_elevation_deg: float | None = None
    min_elevation_rad: float | None = None

    lookahead_days: int = 0

    created_at: datetime | None = None
    updated_at: datetime | None = None


class DeviceType(enum.IntEnum):
    WEB_PUSH = 1
    FCM = 2
    EMAIL = 3
    WEBHOOK = 4


@dataclasses.dataclass
class Device:
    id: int | None
    user_id: int

    type: DeviceType

    address: str

    enabled: bool

    created_at: datetime | None = None
