import builtins
import dataclasses

import sqlalchemy
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from notification.domain import models
from notification.infrastructure.db.models import DeviceORM


@dataclasses.dataclass
class DeviceRepository:
    session_factory: async_sessionmaker[AsyncSession]

    async def create(self, device: models.Device) -> models.Device:
        async with self.session_factory() as session:
            orm = DeviceORM(
                user_id=device.user_id,
                type=device.type,
                address=device.address,
                enabled=True,
            )

            session.add(orm)
            await session.commit()
            await session.refresh(orm)

            return self._to_domain(orm)

    async def list(self, user_id: int) -> list[models.Device]:
        async with self.session_factory() as session:
            stmt = sqlalchemy.select(DeviceORM).where(DeviceORM.user_id == user_id)
            result = await session.execute(stmt)

            return [self._to_domain(row) for row in result.scalars().all()]

    async def list_enabled(self, user_id: int) -> builtins.list[models.Device]:
        async with self.session_factory() as session:
            stmt = sqlalchemy.select(DeviceORM).where(
                DeviceORM.user_id == user_id,
                DeviceORM.enabled.is_(True),
            )
            result = await session.execute(stmt)

            return [self._to_domain(row) for row in result.scalars().all()]

    async def delete(self, device_id: int) -> None:
        async with self.session_factory() as session:
            orm = await session.get(DeviceORM, device_id)

            if orm is None:
                raise ValueError(f"device {device_id} not found")

            await session.delete(orm)
            await session.commit()

    @staticmethod
    def _to_domain(orm: DeviceORM) -> models.Device:
        try:
            device_type = models.DeviceType(orm.type)
        except ValueError as e:
            raise ValueError("invalid device type") from e

        return models.Device(
            id=orm.id,
            user_id=orm.user_id,
            type=device_type,
            address=orm.address,
            enabled=orm.enabled,
            created_at=orm.created_at,
        )
