import asyncio
import datetime as dt

from notification.adapter.trajectory_client import TrajectoryClient
from notification.domain import models
from notification.repository.notification_job_repo import NotificationJobRepository
from notification.repository.subscription_repo import SubscriptionRepository

REFILL_THRESHOLD = dt.timedelta(hours=1)
SCHEDULING_WINDOW = dt.timedelta(days=2)


class NotificationScheduler:
    def __init__(
        self,
        subscription_repository: SubscriptionRepository,
        job_repository: NotificationJobRepository,
        trajectory_client: TrajectoryClient,
    ):
        self._subscriptions = subscription_repository
        self._jobs = job_repository
        self._trajectory = trajectory_client

    async def run_cycle(self) -> None:
        now = dt.datetime.now(dt.UTC)

        subscriptions = await self._subscriptions.list_needing_schedule(
            refill_before=now + REFILL_THRESHOLD,
        )

        await asyncio.gather(
            *(self._schedule_subscription(subscription, now) for subscription in subscriptions)
        )

    async def _schedule_subscription(
        self,
        subscription: models.Subscription,
        now: dt.datetime,
    ) -> None:
        if subscription.scheduled_until is not None:
            start = max(subscription.scheduled_until, now)
        else:
            start = now

        end = start + SCHEDULING_WINDOW

        passes = await self._trajectory.get_passes(
            subscription=subscription,
            start=start,
            end=end,
        )

        jobs = []

        for p in passes:
            jobs.append(
                models.NotificationJob(
                    id=None,
                    subscription_id=subscription.id,
                    scheduled_time=p.aos - dt.timedelta(seconds=subscription.notify_before_seconds),
                    aos=p.aos,
                    los=p.los,
                    max_elevation_time=p.max_elevation_time,
                    max_elevation=p.max_elevation,
                    status=models.NotificationJobStatus.PENDING,
                )
            )

        await self._jobs.insert_jobs_idempotent(jobs)
        await self._subscriptions.advance_schedule_until(subscription.id, end)
