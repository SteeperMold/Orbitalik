import datetime as dt

import grpc

from notification.domain.models import Pass, Subscription
from notification.proto.trajectory_pb2 import (
    GeodeticInput,
    PassPredictionRequest,
    SatelliteIdentifier,
    TimeRange,
    UnitSettings,
)
from notification.proto.trajectory_pb2_grpc import TrajectoryServiceStub


class TrajectoryClient:
    def __init__(self, stub: TrajectoryServiceStub) -> None:
        self._stub = stub

    @classmethod
    def create(cls, target: str):
        channel = grpc.aio.insecure_channel(target)
        stub = TrajectoryServiceStub(channel)
        return cls(stub)

    async def get_passes(
        self,
        subscription: Subscription,
        start: dt.datetime,
        end: dt.datetime,
    ) -> list[Pass]:
        request = self._build_request(subscription, start, end)

        response = await self._stub.GetPasses(request)

        return [self._map_pass(p) for p in response.passes]

    def _build_request(
        self,
        sub: Subscription,
        start: dt.datetime,
        end: dt.datetime,
    ) -> PassPredictionRequest:
        return PassPredictionRequest(
            satellites=[self._build_satellite(sub)],
            observer=self._build_observer(sub),
            range=TimeRange(
                start=start,
                end=end,
            ),
            units=UnitSettings(
                distance_unit=UnitSettings.DISTANCE_UNIT_KILOMETERS,
                angle_unit=UnitSettings.ANGLE_UNIT_DEGREES,
            ),
            min_peak_elevation_deg=sub.min_peak_elevation_deg,
            min_elevation_deg=sub.min_elevation_deg,
        )

    def _map_pass(self, p) -> Pass:
        return Pass(
            satellite=p.satellite,
            aos=p.aos.ToDatetime().replace(tzinfo=dt.UTC),
            aos_azimuth=p.aos_azimuth,
            los=p.los.ToDatetime().replace(tzinfo=dt.UTC),
            los_azimuth=p.los_azimuth,
            max_elevation_time=p.max_elevation_time.ToDatetime().replace(tzinfo=dt.UTC),
            max_elevation=p.max_elevation,
            max_elevation_azimuth=p.max_elevation_azimuth,
            duration_seconds=p.duration_seconds,
        )

    @staticmethod
    def _build_satellite(sub: Subscription) -> SatelliteIdentifier:
        if sub.satellite.norad_id is not None:
            return SatelliteIdentifier(norad_id=sub.satellite.norad_id)

        if sub.satellite.satellite_name is not None:
            return SatelliteIdentifier(satellite_name=sub.satellite.satellite_name)

        raise ValueError("Invalid satellite: missing norad_id and name")

    @staticmethod
    def _build_observer(sub: Subscription) -> GeodeticInput:
        return GeodeticInput(
            lat_deg=sub.observer.lat_deg,
            lat_rad=sub.observer.lat_rad,
            lon_deg=sub.observer.lon_deg,
            lon_rad=sub.observer.lon_rad,
            alt_m=sub.observer.alt_m,
            alt_km=sub.observer.alt_km,
        )
