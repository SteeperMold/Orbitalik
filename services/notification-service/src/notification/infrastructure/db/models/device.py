import enum
from datetime import datetime

from sqlalchemy import BigInteger, Boolean, DateTime, Integer, String
from sqlalchemy.orm import Mapped, mapped_column

from notification.infrastructure.db.base import Base


class DeviceTypeEnum(enum.IntEnum):
    DEVICE_TYPE_UNSPECIFIED = 0
    DEVICE_TYPE_WEB_PUSH = 1
    DEVICE_TYPE_FCM = 2
    DEVICE_TYPE_EMAIL = 3
    DEVICE_TYPE_WEBHOOK = 4


class DeviceORM(Base):
    __tablename__ = "devices"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)

    user_id: Mapped[int] = mapped_column(BigInteger, index=True, nullable=False)

    type: Mapped[int] = mapped_column(Integer, nullable=False)

    address: Mapped[str] = mapped_column(String, nullable=False)

    enabled: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)

    created_at: Mapped[datetime] = mapped_column(
        DateTime,
        default=datetime.utcnow,
        nullable=False,
    )
