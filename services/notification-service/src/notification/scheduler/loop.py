import asyncio

import structlog

from notification.scheduler.scheduler import NotificationScheduler

logger = structlog.get_logger()


async def run_scheduler(service: NotificationScheduler, interval_seconds: int) -> None:
    while True:
        try:
            await service.run_cycle()
        except Exception:
            logger.exception("Scheduler cycle failed")

        await asyncio.sleep(interval_seconds)
