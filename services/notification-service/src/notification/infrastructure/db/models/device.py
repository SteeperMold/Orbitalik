import datetime as dt
import uuid

from sqlalchemy import Boolean, DateTime, Integer, String, Uuid
from sqlalchemy.orm import Mapped, mapped_column

from notification.infrastructure.db.base import Base


class DeviceORM(Base):
    __tablename__ = "devices"

    id: Mapped[uuid.UUID] = mapped_column(
        Uuid(as_uuid=True),
        primary_key=True,
    )

    user_id: Mapped[uuid.UUID] = mapped_column(Uuid(as_uuid=True), index=True, nullable=False)

    type: Mapped[int] = mapped_column(Integer, nullable=False)

    address: Mapped[str] = mapped_column(String, nullable=False)

    enabled: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)

    created_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: dt.datetime.now(dt.UTC),
        nullable=False,
    )
