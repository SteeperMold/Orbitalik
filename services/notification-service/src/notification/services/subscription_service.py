import base64
import binascii
import dataclasses
import json

from notification.domain import command, limits, models
from notification.repository.subscription_repo import SubscriptionRepository


@dataclasses.dataclass
class SubscriptionService:
    repository: SubscriptionRepository
    max_page_size: int

    async def create_subscription(
        self,
        cmd: command.CreateSubscriptionCommand,
    ) -> models.Subscription:
        if cmd.notify_before_seconds is None:
            raise ValueError("notify_before_seconds is required")

        if cmd.notify_before_seconds < 0:
            raise ValueError("notify_before_seconds must be >= 0")

        if cmd.notify_before_seconds > limits.MAX_INT32:
            raise ValueError("notify_before_seconds exceeds int32 limit")

        if cmd.lookahead_days is None:
            raise ValueError("lookahead_days is required")

        if cmd.lookahead_days < limits.MIN_LOOKAHEAD_DAYS:
            raise ValueError(f"lookahead_days must be >= {limits.MAX_LOOKAHEAD_DAYS}")

        if cmd.lookahead_days > limits.MAX_LOOKAHEAD_DAYS:
            raise ValueError("lookahead_days is unreasonably large")

        obs = cmd.observer

        if obs.lat_deg is not None and not (limits.MIN_LAT <= obs.lat_deg <= limits.MAX_LAT):
            raise ValueError(f"observer_lat_deg out of range [{limits.MIN_LAT}, {limits.MAX_LAT}]")

        if obs.lon_deg is not None and not (limits.MIN_LON <= obs.lon_deg <= limits.MAX_LON):
            raise ValueError(f"observer_lon_deg out of range [{limits.MIN_LON}, {limits.MAX_LON}]")

        if obs.alt_m is not None and obs.alt_m < 0:
            raise ValueError("observer_alt_m must be >= 0")

        if cmd.satellite is None:
            raise ValueError("satellite is required")

        return await self.repository.create(cmd)

    async def get_subscription(self, subscription_id: int) -> models.Subscription:
        subscription = await self.repository.get(subscription_id)

        if subscription is None:
            raise ValueError(f"subscription {subscription_id} not found")

        return subscription

    async def list_subscriptions(
        self,
        user_id: int,
        enabled: bool | None,
        page_size: int,
        page_token: str | None,
    ) -> tuple[list[models.Subscription], str]:
        if page_size <= 0 or page_size > self.max_page_size:
            raise ValueError(f"page_size must be 1..{self.max_page_size}")

        offset = self._decode_page_token(page_token)

        query = command.ListSubscriptionsCommand(
            user_id=user_id,
            enabled=enabled,
            offset=offset,
            limit=page_size,
        )

        items, has_next = await self.repository.list(query)

        next_token = self._encode_page_token(offset + page_size) if has_next else ""

        return items, next_token

    async def update_subscription(
        self,
        subscription_id: int,
        cmd: command.UpdateSubscriptionCommand,
    ) -> models.Subscription:
        if cmd.notify_before_seconds is not None and cmd.notify_before_seconds > limits.MAX_INT32:
            raise ValueError("notify_before_seconds exceeds int32 limit")

        if cmd.notify_before_seconds is not None and cmd.notify_before_seconds < 0:
            raise ValueError("notify_before_seconds must be >= 0")

        return await self.repository.update(subscription_id, cmd)

    async def delete_subscription(self, subscription_id: int) -> None:
        deleted = await self.repository.delete(subscription_id)

        if not deleted:
            raise ValueError(f"subscription {subscription_id} not found")

    async def set_subscription_status(
        self,
        subscription_id: int,
        enabled: bool,
    ) -> models.Subscription:
        subscription = await self.repository.set_status(
            subscription_id=subscription_id,
            enabled=enabled,
        )

        if subscription is None:
            raise ValueError(f"subscription {subscription_id} not found")

        return subscription

    @staticmethod
    def _encode_page_token(offset: int) -> str:
        payload = json.dumps({"offset": offset}).encode()
        return base64.urlsafe_b64encode(payload).decode()

    @staticmethod
    def _decode_page_token(token: str | None) -> int:
        if not token:
            return 0

        try:
            payload = json.loads(base64.urlsafe_b64decode(token.encode()).decode())
            return int(payload["offset"])
        except (binascii.Error, json.JSONDecodeError, KeyError, ValueError) as exc:
            raise ValueError("invalid page_token") from exc
