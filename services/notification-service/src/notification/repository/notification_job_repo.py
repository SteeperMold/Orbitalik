import dataclasses
import datetime as dt

import sqlalchemy
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

    async def fetch_due_jobs(self, now: dt.datetime) -> list[models.NotificationJob]:
        async with self.session_factory() as session:
            stmt = (
                sqlalchemy.select(NotificationJobORM)
                .where(
                    NotificationJobORM.status == models.NotificationJobStatus.PENDING,
                    NotificationJobORM.scheduled_time <= now,
                )
                .order_by(NotificationJobORM.scheduled_time.asc())
            )

            result = await session.execute(stmt)
            return [self._to_domain(j) for j in result.scalars().all()]

    async def mark_sent(self, job_id: int) -> None:
        async with self.session_factory() as session:
            stmt = (
                sqlalchemy.update(NotificationJobORM)
                .where(NotificationJobORM.id == job_id)
                .values(
                    status=models.NotificationJobStatus.SENT,
                    updated_at=dt.datetime.now(dt.UTC),
                )
            )

            await session.execute(stmt)
            await session.commit()

    async def mark_failed(self, job_id: int, error: str) -> None:
        async with self.session_factory() as session:
            stmt = (
                sqlalchemy.update(NotificationJobORM)
                .where(NotificationJobORM.id == job_id)
                .values(
                    status=models.NotificationJobStatus.FAILED,
                    last_error=error,
                    attempts=NotificationJobORM.attempts + 1,
                    updated_at=dt.datetime.now(dt.UTC),
                )
            )

            await session.execute(stmt)
            await session.commit()

    @staticmethod
    def _to_domain(job_orm: NotificationJobORM) -> models.NotificationJob:
        return models.NotificationJob(
            id=job_orm.id,
            subscription_id=job_orm.subscription_id,
            user_id=job_orm.user_id,
            scheduled_time=job_orm.scheduled_time,
            aos=job_orm.aos,
            los=job_orm.los,
            max_elevation_time=job_orm.max_elevation_time,
            max_elevation=job_orm.max_elevation,
            status=models.NotificationJobStatus(job_orm.status),
        )
