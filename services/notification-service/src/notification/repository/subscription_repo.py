import dataclasses

import sqlalchemy
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from notification.domain import command, models
from notification.infrastructure.db.models.subscription import SubscriptionORM


@dataclasses.dataclass
class SubscriptionRepository:
    session_factory: async_sessionmaker[AsyncSession]

    async def create(self, cmd: command.CreateSubscriptionCommand) -> models.Subscription:
        async with self.session_factory() as session:
            obs = cmd.observer

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

            session.add(orm)
            await session.commit()
            await session.refresh(orm)

            return self._to_domain(orm)

    async def get(self, subscription_id: int) -> models.Subscription | None:
        async with self.session_factory() as session:
            orm = await session.get(SubscriptionORM, subscription_id)

            if orm is None:
                return None

            return self._to_domain(orm)

    async def list(
        self,
        cmd: command.ListSubscriptionsCommand,
    ) -> tuple[list[models.Subscription], bool]:
        async with self.session_factory() as session:
            offset = cmd.offset
            limit = cmd.limit + 1  # fetch extra to detect next page

            stmt = sqlalchemy.select(SubscriptionORM).where(SubscriptionORM.user_id == cmd.user_id)

            if cmd.enabled is not None:
                stmt = stmt.where(SubscriptionORM.enabled == cmd.enabled)

            stmt = stmt.offset(offset).limit(limit)

            result = await session.execute(stmt)
            rows = result.scalars().all()

            has_next = len(rows) > cmd.limit
            rows = rows[: cmd.limit]

            return [self._to_domain(r) for r in rows], has_next

    async def update(
        self,
        subscription_id: int,
        cmd: command.UpdateSubscriptionCommand,
    ) -> models.Subscription:
        async with self.session_factory() as session:
            orm = await session.get(SubscriptionORM, subscription_id)

            if orm is None:
                raise ValueError(f"subscription {subscription_id} not found")

            if cmd.enabled is not None:
                orm.enabled = cmd.enabled

            if cmd.notify_before_seconds is not None:
                orm.notify_before_seconds = cmd.notify_before_seconds

            if cmd.min_peak_elevation_deg is not None:
                orm.min_peak_elevation_deg = cmd.min_peak_elevation_deg

            if cmd.min_elevation_deg is not None:
                orm.min_elevation_deg = cmd.min_elevation_deg

            await session.commit()
            await session.refresh(orm)

            return self._to_domain(orm)

    async def delete(self, subscription_id: int) -> bool:
        async with self.session_factory() as session:
            orm = await session.get(SubscriptionORM, subscription_id)

            if orm is None:
                return False

            await session.delete(orm)
            await session.commit()

            return True

    async def set_status(
        self,
        subscription_id: int,
        enabled: bool,
    ) -> models.Subscription | None:
        async with self.session_factory() as session:
            orm = await session.get(SubscriptionORM, subscription_id)

            if orm is None:
                return None

            orm.enabled = enabled

            await session.commit()
            await session.refresh(orm)

            return self._to_domain(orm)

    @staticmethod
    def _to_domain(orm: SubscriptionORM) -> models.Subscription:
        return models.Subscription(
            id=orm.id,
            user_id=orm.user_id,
            satellite=models.SatelliteIdentifier(
                norad_id=orm.norad_id,
                satellite_name=orm.satellite_name,
            ),
            observer=models.GeodeticInput(
                lat_deg=orm.observer_lat_deg,
                lat_rad=orm.observer_lat_rad,
                lon_deg=orm.observer_lon_deg,
                lon_rad=orm.observer_lon_rad,
                alt_m=orm.observer_alt_m,
                alt_km=orm.observer_alt_km,
            ),
            enabled=orm.enabled,
            notify_before_seconds=orm.notify_before_seconds,
            min_peak_elevation_deg=orm.min_peak_elevation_deg,
            min_elevation_deg=orm.min_elevation_deg,
            lookahead_days=orm.lookahead_days,
            created_at=orm.created_at,
            updated_at=orm.updated_at,
        )
