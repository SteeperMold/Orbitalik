import datetime as dt
import uuid
from typing import TYPE_CHECKING

from sqlalchemy import Boolean, DateTime, Float, Integer, String, Uuid
from sqlalchemy.orm import Mapped, mapped_column, relationship

from notification.infrastructure.db.base import Base
from notification.infrastructure.db.models.notification_job import NotificationJobORM

if TYPE_CHECKING:
    from notification.infrastructure.db.models.notification_job import NotificationJobORM


class SubscriptionORM(Base):
    __tablename__ = "subscriptions"

    id: Mapped[uuid.UUID] = mapped_column(
        Uuid(as_uuid=True),
        primary_key=True,
    )

    user_id: Mapped[uuid.UUID] = mapped_column(Uuid(as_uuid=True), index=True)

    norad_id: Mapped[int | None] = mapped_column(nullable=True)
    satellite_name: Mapped[str | None] = mapped_column(String, nullable=True)

    observer_lat_deg: Mapped[float | None] = mapped_column(Float, nullable=True)
    observer_lat_rad: Mapped[float | None] = mapped_column(Float, nullable=True)

    observer_lon_deg: Mapped[float | None] = mapped_column(Float, nullable=True)
    observer_lon_rad: Mapped[float | None] = mapped_column(Float, nullable=True)

    observer_alt_m: Mapped[float | None] = mapped_column(Float, nullable=True)
    observer_alt_km: Mapped[float | None] = mapped_column(Float, nullable=True)

    enabled: Mapped[bool] = mapped_column(Boolean, default=True)

    notify_before_seconds: Mapped[int] = mapped_column(Integer)

    min_peak_elevation_deg: Mapped[float | None] = mapped_column(Float, nullable=True)
    min_peak_elevation_rad: Mapped[float | None] = mapped_column(Float, nullable=True)

    min_elevation_deg: Mapped[float | None] = mapped_column(Float, nullable=True)
    min_elevation_rad: Mapped[float | None] = mapped_column(Float, nullable=True)

    scheduled_until: Mapped[dt.datetime | None] = mapped_column(
        DateTime(timezone=True),
        nullable=True,
    )

    jobs: Mapped[list["NotificationJobORM"]] = relationship(
        back_populates="subscription",
        cascade="all, delete-orphan",
        passive_deletes=True,
    )

    created_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: dt.datetime.now(dt.UTC),
    )

    updated_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: dt.datetime.now(dt.UTC),
        onupdate=lambda: dt.datetime.now(dt.UTC),
    )
