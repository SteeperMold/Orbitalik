from notification.domain.command import CreateSubscriptionCommand
from notification.domain.models import Subscription
from notification.infrastructure.db.models.subscription import SubscriptionORM
from notification.repository.subscription_repo import SubscriptionRepository

MAX_INT32 = 2_147_483_647


class NotificationService:
    def __init__(self, repository: SubscriptionRepository) -> None:
        self.repository = repository

    async def create_subscription(self, cmd: CreateSubscriptionCommand) -> Subscription:
        if cmd.notify_before_seconds is None:
            raise ValueError("notify_before_seconds is required")

        if cmd.notify_before_seconds < 0:
            raise ValueError("notify_before_seconds must be >= 0")

        if cmd.notify_before_seconds > MAX_INT32:
            raise ValueError("notify_before_seconds exceeds int32 limit")

        if cmd.lookahead_days is None:
            raise ValueError("lookahead_days is required")

        if cmd.lookahead_days <= 0:
            raise ValueError("lookahead_days must be > 0")

        if cmd.lookahead_days > 3650:
            raise ValueError("lookahead_days is unreasonably large")

        obs = cmd.observer

        if obs.lat_deg is not None and not (-90.0 <= obs.lat_deg <= 90.0):
            raise ValueError("observer_lat_deg out of range [-90, 90]")

        if obs.lon_deg is not None and not (-180.0 <= obs.lon_deg <= 180.0):
            raise ValueError("observer_lon_deg out of range [-180, 180]")

        if obs.alt_m is not None and obs.alt_m < 0:
            raise ValueError("observer_alt_m must be >= 0")

        if cmd.satellite is None:
            raise ValueError("satellite is required")

        orm = SubscriptionORM(
            user_id=cmd.user_id,

            norad_id=cmd.satellite.norad_id,
            satellite_name=cmd.satellite.satellite_name,

            observer_lat_deg=obs.lat_deg,
            observer_lat_rad=obs.lat_rad,
            observer_lon_deg=obs.lon_deg,
            observer_lon_rad=obs.lon_rad,
            observer_alt_m=obs.alt_m,
            observer_alt_km=obs.alt_km,

            notify_before_seconds=cmd.notify_before_seconds,

            min_peak_elevation_deg=cmd.min_peak_elevation_deg,
            min_peak_elevation_rad=cmd.min_peak_elevation_rad,
            min_elevation_deg=cmd.min_elevation_deg,
            min_elevation_rad=cmd.min_elevation_rad,

            lookahead_days=cmd.lookahead_days,
        )

        saved = await self.repository.create(orm)

        return Subscription(
            id=str(saved.id),
            user_id=saved.user_id,
            satellite=cmd.satellite,
            observer=cmd.observer,

            enabled=saved.enabled,

            notify_before_seconds=saved.notify_before_seconds,

            min_peak_elevation_deg=saved.min_peak_elevation_deg,
            min_peak_elevation_rad=saved.min_peak_elevation_rad,
            min_elevation_deg=saved.min_elevation_deg,
            min_elevation_rad=saved.min_elevation_rad,

            lookahead_days=saved.lookahead_days,

            created_at=saved.created_at,
            updated_at=saved.updated_at,
        )
