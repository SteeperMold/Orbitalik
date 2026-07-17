import dataclasses
import datetime as dt
import enum


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

    scheduled_until: dt.datetime | None = None

    created_at: dt.datetime | None = None
    updated_at: dt.datetime | None = None


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

    created_at: dt.datetime | None = None


class NotificationJobStatus(enum.IntEnum):
    PENDING = 0
    PROCESSING = 1
    SENT = 2
    FAILED = 3
    CANCELLED = 4


@dataclasses.dataclass
class NotificationJob:
    id: int | None

    user_id: int

    subscription_id: int

    scheduled_time: dt.datetime
    status: NotificationJobStatus

    satellite: SatelliteIdentifier
    aos: dt.datetime
    los: dt.datetime
    max_elevation_time: dt.datetime
    max_elevation: float

    attempts: int = 0
    last_error: str | None = None

    created_at: dt.datetime | None = None
    updated_at: dt.datetime | None = None


@dataclasses.dataclass
class Pass:
    satellite: SatelliteIdentifier

    aos: dt.datetime
    aos_azimuth: float

    max_elevation_time: dt.datetime
    max_elevation: float
    max_elevation_azimuth: float

    los: dt.datetime
    los_azimuth: float

    duration_seconds: int
