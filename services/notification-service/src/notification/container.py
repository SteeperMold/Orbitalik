import dataclasses

from notification.transport.grpc.services import NotificationServicer


@dataclasses.dataclass
class Services:
    notification_servicer: NotificationServicer
