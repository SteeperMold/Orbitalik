import asyncio

import firebase_admin
import httpx

from notification.adapter.trajectory_client import TrajectoryClient
from notification.container import Services
from notification.infrastructure.db.session import create_engine, create_session_factory
from notification.infrastructure.logger import configure_logging
from notification.infrastructure.settings import AppEnv, Settings
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
from notification.worker import push_provider
from notification.worker.notification_worker import NotificationWorker
from notification.worker.push_service import PushService


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
        _subscriptions=sub_repo,
        _jobs=job_repo,
        _trajectory=trajectory_client,
        _scheduling_window=settings.scheduling_window,
        _refill_threshold=settings.scheduling_refill_threshold,
    )

    if settings.app_env == AppEnv.PRODUCTION:
        cred = firebase_admin.credentials.Certificate("service-account.json")
        firebase_admin.initialize_app(cred)

    push_service = PushService(
        _fcm=push_provider.FCMPushProvider(),
        _email=push_provider.EmailPushProvider(settings.smtp_host, settings.smtp_port),
        _webhook=push_provider.WebhookPushProvider(
            httpx.AsyncClient(
                timeout=settings.request_timeout_seconds,
                headers={"User-Agent": "notification-service"},
            )
        ),
        _file=push_provider.FilePushProvider(),
        is_development=settings.app_env == AppEnv.DEVELOPMENT,
    )
    notification_worker = NotificationWorker(
        job_repo,
        device_repo,
        push_service,
        settings.worker_polling_interval_seconds,
    )

    await asyncio.gather(
        grpc_server.start(),
        http_server.start(),
        run_scheduler(notification_scheduler, settings.scheduling_interval_seconds),
        notification_worker.run_worker(),
    )

    await grpc_server.server.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(main())
