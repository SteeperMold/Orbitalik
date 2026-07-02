import dataclasses
import datetime as dt
import json
import pathlib
from email.message import EmailMessage

import aiosmtplib
import firebase_admin
import httpx

from notification.domain import models


@dataclasses.dataclass
class WebhookPushProvider:
    _client: httpx.AsyncClient

    async def send(self, device, title: str, body: str) -> None:
        payload = {
            "title": title,
            "body": body,
        }

        await self._client.post(device.address, json=payload, timeout=5)


@dataclasses.dataclass
class EmailPushProvider:
    smtp_host: str
    smtp_port: int

    async def send(self, device: models.Device, title: str, body: str) -> None:
        msg = EmailMessage()
        msg["To"] = device.address
        msg["Subject"] = title
        msg.set_content(body)

        await aiosmtplib.send(
            msg,
            hostname=self.smtp_host,
            port=self.smtp_port,
        )


class FCMPushProvider:
    async def send(self, device: models.Device, title: str, body: str) -> None:
        message = firebase_admin.messaging.Message(
            token=device.address,  # FCM token
            notification=firebase_admin.messaging.Notification(
                title=title,
                body=body,
            ),
        )

        firebase_admin.messaging.send(message)


@dataclasses.dataclass
class FilePushProvider:
    file_path: str = "push_log.jsonl"

    def __post_init__(self) -> None:
        self._path = pathlib.Path(self.file_path)
        self._path.parent.mkdir(parents=True, exist_ok=True)

    async def send(
        self,
        device: models.Device,
        title: str,
        body: str,
    ) -> None:
        event = {
            "timestamp": dt.datetime.now(dt.UTC).isoformat(),
            "device_id": device.id,
            "user_id": device.user_id,
            "device_type": device.type,
            "address": device.address,
            "title": title,
            "body": body,
        }

        with self._path.open("a", encoding="utf-8") as f:
            f.write(json.dumps(event) + "\n")
