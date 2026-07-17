import grpc
import structlog

from notification.container import Services
from notification.infrastructure.settings import Settings
from notification.proto import notification_pb2_grpc

logger = structlog.get_logger()


class GRPCServer:
    def __init__(self, services: Services, settings: Settings) -> None:
        self.services = services
        self.settings = settings

        self.server = grpc.aio.server()

        notification_pb2_grpc.add_NotificationServiceServicer_to_server(
            services.notification_servicer,
            self.server,
        )

        self.server.add_insecure_port(f"[::]:{settings.grpc_port}")

    async def start(self) -> None:
        await self.server.start()
        logger.info(
            "grpc_server_started",
            port=self.settings.grpc_port,
        )

    async def stop(self) -> None:
        await self.server.stop(grace=5)
