import dataclasses

from notification.domain import models
from notification.worker.push_provider import (
    EmailPushProvider,
    FCMPushProvider,
    FilePushProvider,
    WebhookPushProvider,
)


@dataclasses.dataclass
class PushService:
    _fcm: FCMPushProvider
    _email: EmailPushProvider
    _webhook: WebhookPushProvider
    _file: FilePushProvider
    is_development: bool

    async def send(
        self,
        device: models.Device,
        title: str,
        body: str,
    ) -> None:
        if self.is_development:
            await self._file.send(device, title, body)
            return

        if device.type == models.DeviceType.FCM:
            await self._fcm.send(device, title, body)

        elif device.type == models.DeviceType.EMAIL:
            await self._email.send(device, title, body)

        elif device.type == models.DeviceType.WEBHOOK:
            await self._webhook.send(device, title, body)

        else:
            raise ValueError(f"Unknown device type: {device.type}")
