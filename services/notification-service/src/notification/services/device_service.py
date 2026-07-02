import dataclasses

from notification.domain import models
from notification.repository.device_repo import DeviceRepository


@dataclasses.dataclass
class DeviceService:
    repository: DeviceRepository

    async def register(
        self,
        user_id: int,
        device_type: int,
        address: str,
    ) -> models.Device:
        if user_id <= 0:
            raise ValueError("invalid user_id")

        if device_type <= 0:
            raise ValueError("invalid device type")

        if not address:
            raise ValueError("address is required")

        device = models.Device(
            id=None,
            user_id=user_id,
            type=self._map_device_type(device_type),
            address=address,
            enabled=True,
        )

        return await self.repository.create(device)

    async def list(self, user_id: int) -> list[models.Device]:
        return await self.repository.list(user_id)

    async def delete_device(self, device_id: int) -> None:
        return await self.repository.delete(device_id)

    @staticmethod
    def _map_device_type(value: int) -> models.DeviceType:
        try:
            return models.DeviceType(value)
        except ValueError as e:
            raise ValueError("invalid device type") from e
