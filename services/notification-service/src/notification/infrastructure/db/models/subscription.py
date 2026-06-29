from datetime import datetime

from sqlalchemy import BigInteger, String, Boolean, Integer, Float, DateTime
from sqlalchemy.orm import Mapped, mapped_column

from notification.infrastructure.db.base import Base


class SubscriptionORM(Base):
    __tablename__ = "subscriptions"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)

    user_id: Mapped[int] = mapped_column(BigInteger, index=True)

    norad_id: Mapped[int | None] = mapped_column(BigInteger, nullable=True)
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

    lookahead_days: Mapped[int] = mapped_column(Integer)

    created_at: Mapped[datetime] = mapped_column(DateTime, default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
