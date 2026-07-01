import asyncio

from notification.adapter.trajectory_client import TrajectoryClient
from notification.container import Services
from notification.infrastructure.db.session import create_engine, create_session_factory
from notification.infrastructure.logger import configure_logging
from notification.infrastructure.settings import Settings
from notification.repository.device_repo import DeviceRepository
from notification.repository.notification_job_repo import NotificationJobRepository
from notification.repository.subscription_repo import SubscriptionRepository
from notification.scheduler.loop import run_scheduler
from notification.scheduler.scheduler import NotificationScheduler
from notification.services.device_service import DeviceService
from notification.services.subscription_service import SubscriptionService
from notification.transport.grpc.server import GRPCServer
from notification.transport.grpc.services import NotificationServicer
from notification.transport.http.server import HTTPServer


async def main() -> None:
    settings = Settings()

    configure_logging(settings.app_env)

    engine = create_engine(settings.database_url, settings.app_env)
    session_factory = create_session_factory(engine)

    sub_repo = SubscriptionRepository(session_factory)
    device_repo = DeviceRepository(session_factory)
    job_repo = NotificationJobRepository(session_factory)

    notification_svc = SubscriptionService(
        sub_repo,
        settings.max_page_size,
    )
    device_svc = DeviceService(device_repo)

    trajectory_client = TrajectoryClient.create(settings.trajectory_service_url)

    servicer = NotificationServicer(notification_svc, device_svc)

    grpc_server = GRPCServer(
        services=Services(notification_servicer=servicer),
        settings=settings,
    )
    http_server = HTTPServer(settings)

    notification_scheduler = NotificationScheduler(
        subscription_repository=sub_repo,
        job_repository=job_repo,
        trajectory_client=trajectory_client,
    )

    await asyncio.gather(
        grpc_server.start(),
        http_server.start(),
        run_scheduler(notification_scheduler),
    )

    await grpc_server.server.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(main())
