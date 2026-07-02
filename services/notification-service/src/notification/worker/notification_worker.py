import asyncio
import dataclasses
import datetime as dt

import structlog

from notification.domain import models
from notification.repository.device_repo import DeviceRepository
from notification.repository.notification_job_repo import NotificationJobRepository
from notification.worker.push_service import PushService

logger = structlog.get_logger()


@dataclasses.dataclass
class NotificationWorker:
    _jobs: NotificationJobRepository
    _devices: DeviceRepository
    _push: PushService
    _polling_interval_seconds: int

    async def run_worker(self) -> None:
        while True:
            try:
                await self._process_cycle()
            except Exception:
                logger.exception("worker cycle failed")

            await asyncio.sleep(self._polling_interval_seconds)

    async def _process_cycle(self) -> None:
        now = dt.datetime.now(dt.UTC)

        jobs = await self._jobs.fetch_due_jobs(now)

        await asyncio.gather(*(self._process_job(j) for j in jobs))

    async def _process_job(self, job: models.NotificationJob) -> None:
        try:
            devices = await self._devices.list_enabled(job.user_id)

            if not devices:
                await self._jobs.mark_sent(job.id)
                return

            title = f"{job.satellite.satellite_name} (NORAD ID {job.satellite.norad_id}) пролетает"
            body = f"{job.satellite.satellite_name} станет виден в {job.aos:%H%M UTC}."

            await asyncio.gather(
                *(
                    self._push.send(
                        device=device,
                        title=title,
                        body=body,
                    )
                    for device in devices
                )
            )

            await self._jobs.mark_sent(job.id)

        except Exception as e:
            await self._jobs.mark_failed(job.id, str(e))
