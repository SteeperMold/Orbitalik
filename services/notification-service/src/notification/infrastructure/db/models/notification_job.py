import datetime as dt

from sqlalchemy import BigInteger, DateTime, Float, ForeignKey, Integer, SmallInteger, Text
from sqlalchemy.orm import Mapped, mapped_column

from notification.domain.models import NotificationJobStatus
from notification.infrastructure.db.base import Base


class NotificationJobORM(Base):
    __tablename__ = "notification_jobs"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)

    subscription_id: Mapped[int] = mapped_column(
        BigInteger,
        ForeignKey("subscriptions.id", ondelete="CASCADE"),
        index=True,
    )

    status: Mapped[int] = mapped_column(SmallInteger, default=NotificationJobStatus.PENDING)

    scheduled_time: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))

    aos: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    los: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    max_elevation_time: Mapped[dt.datetime] = mapped_column(DateTime(timezone=True))
    max_elevation: Mapped[float] = mapped_column(Float)

    attempts: Mapped[int] = mapped_column(Integer, default=0)

    last_error: Mapped[str | None] = mapped_column(Text, nullable=True)

    created_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True), default=lambda: dt.datetime.now(dt.UTC)
    )

    updated_at: Mapped[dt.datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: dt.datetime.now(dt.UTC),
        onupdate=lambda: dt.datetime.now(dt.UTC),
    )
