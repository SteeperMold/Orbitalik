import dataclasses

import sqlalchemy.dialects.postgresql as sqlalchemy_postgres
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from notification.domain import models
from notification.infrastructure.db.models.notification_job import NotificationJobORM


@dataclasses.dataclass
class NotificationJobRepository:
    session_factory: async_sessionmaker[AsyncSession]

    async def insert_jobs_idempotent(self, jobs: list[models.NotificationJob]) -> None:
        if not jobs:
            return

        async with self.session_factory() as session:
            stmt = sqlalchemy_postgres.insert(NotificationJobORM).values(
                [
                    {
                        "subscription_id": j.subscription_id,
                        "scheduled_time": j.scheduled_time,
                        "aos": j.aos,
                        "los": j.los,
                        "max_elevation_time": j.max_elevation_time,
                        "max_elevation": j.max_elevation,
                        "status": j.status.value,
                    }
                    for j in jobs
                ]
            )

            stmt = stmt.on_conflict_do_nothing(index_elements=["subscription_id", "aos"])

            await session.execute(stmt)
            await session.commit()
