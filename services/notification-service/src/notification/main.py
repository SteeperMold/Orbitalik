import asyncio

from notification.container import Services
from notification.infrastructure.db.session import create_engine, create_session_factory
from notification.infrastructure.logger import configure_logging
from notification.infrastructure.settings import Settings
from notification.repository.device_repo import DeviceRepository
from notification.repository.subscription_repo import SubscriptionRepository
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
    notification_svc = SubscriptionService(
        sub_repo,
        settings.max_page_size,
    )
    device_repo = DeviceRepository(session_factory)
    device_svc = DeviceService(device_repo)
    servicer = NotificationServicer(notification_svc, device_svc)

    grpc_server = GRPCServer(
        services=Services(notification_servicer=servicer),
        settings=settings,
    )
    http_server = HTTPServer(settings)

    await asyncio.gather(
        grpc_server.start(),
        http_server.start(),
    )

    await grpc_server.server.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(main())
