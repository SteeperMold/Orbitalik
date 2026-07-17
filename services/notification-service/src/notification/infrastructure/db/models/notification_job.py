import datetime as dt
from typing import TYPE_CHECKING

from sqlalchemy import BigInteger, DateTime, Float, ForeignKey, Integer, SmallInteger, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from notification.domain.models import NotificationJobStatus
from notification.infrastructure.db.base import Base

if TYPE_CHECKING:
    from notification.infrastructure.db.models import SubscriptionORM


class NotificationJobORM(Base):
    __tablename__ = "notification_jobs"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)

    subscription_id: Mapped[int] = mapped_column(
        BigInteger,
        ForeignKey("subscriptions.id", ondelete="CASCADE"),
        index=True,
    )

    subscription: Mapped["SubscriptionORM"] = relationship(
        back_populates="jobs",
    )

    user_id: Mapped[int] = mapped_column(
        BigInteger,
        index=True,
        nullable=False,
    )

    scheduled_time: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))

    status: Mapped[int] = mapped_column(SmallInteger, default=NotificationJobStatus.PENDING)
    attempts: Mapped[int] = mapped_column(Integer, default=0)
    last_error: Mapped[str | None] = mapped_column(Text, nullable=True)

    norad_id: Mapped[int | None] = mapped_column(BigInteger, nullable=True)
    satellite_name: Mapped[str | None] = mapped_column(String, nullable=True)
    aos: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    los: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    max_elevation_time: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    max_elevation: Mapped[float] = mapped_column(Float)

    created_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True), default=lambda: dt.datetime.now(dt.UTC)
    )

    updated_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: dt.datetime.now(dt.UTC),
        onupdate=lambda: dt.datetime.now(dt.UTC),
    )
