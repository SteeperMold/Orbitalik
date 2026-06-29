import dataclasses

from notification.domain.models import GeodeticInput, SatelliteIdentifier


@dataclasses.dataclass(frozen=True)
class CreateSubscriptionCommand:
    user_id: int

    satellite: SatelliteIdentifier
    observer: GeodeticInput

    notify_before_seconds: int

    min_peak_elevation_deg: float | None = None
    min_peak_elevation_rad: float | None = None

    min_elevation_deg: float | None = None
    min_elevation_rad: float | None = None

    lookahead_days: int = 0
