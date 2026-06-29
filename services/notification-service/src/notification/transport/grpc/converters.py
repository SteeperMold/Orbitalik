from google.protobuf.timestamp_pb2 import Timestamp

from notification.domain.command import CreateSubscriptionCommand
from notification.domain.models import SatelliteIdentifier, GeodeticInput, Subscription
from notification.proto import notification_pb2 as pb2


def satellite_from_proto(satellite: pb2.SatelliteIdentifier) -> SatelliteIdentifier:
    kind = satellite.WhichOneof("kind")

    if kind == "norad_id":
        return SatelliteIdentifier(norad_id=satellite.norad_id)
    elif kind == "satellite_name":
        return SatelliteIdentifier(satellite_name=satellite.satellite_name)
    else:
        raise ValueError("satellite identifier is required")


def observer_from_proto(observer: pb2.GeodeticInput) -> GeodeticInput:
    return GeodeticInput(
        lat_deg=observer.lat_deg if observer.HasField("lat_deg") else None,
        lat_rad=observer.lat_rad if observer.HasField("lat_rad") else None,
        lon_deg=observer.lon_deg if observer.HasField("lon_deg") else None,
        lon_rad=observer.lon_rad if observer.HasField("lon_rad") else None,
        alt_m=observer.alt_m if observer.HasField("alt_m") else None,
        alt_km=observer.alt_km if observer.HasField("alt_km") else None,
    )


def create_subscription_command_from_request(request: pb2.CreateSubscriptionRequest) -> CreateSubscriptionCommand:
    return CreateSubscriptionCommand(
        user_id=request.user_id,
        satellite=satellite_from_proto(request.satellite),
        observer=observer_from_proto(request.observer),
        notify_before_seconds=request.notify_before_seconds,

        min_peak_elevation_deg=request.min_peak_elevation_deg
        if request.WhichOneof("min_peak_elevation") == "min_peak_elevation_deg"
        else None,

        min_peak_elevation_rad=request.min_peak_elevation_rad
        if request.WhichOneof("min_peak_elevation") == "min_peak_elevation_rad"
        else None,

        min_elevation_deg=request.min_elevation_deg
        if request.WhichOneof("min_elevation") == "min_elevation_deg"
        else None,

        min_elevation_rad=request.min_elevation_rad
        if request.WhichOneof("min_elevation") == "min_elevation_rad"
        else None,

        lookahead_days=request.lookahead_days,
    )


def subscription_to_proto(subscription: Subscription) -> pb2.Subscription:
    created = Timestamp()
    if subscription.created_at:
        created.FromDatetime(subscription.created_at)

    updated = Timestamp()
    if subscription.updated_at:
        updated.FromDatetime(subscription.updated_at)

    satellite = None
    if subscription.satellite.norad_id is not None:
        satellite = pb2.SatelliteIdentifier(norad_id=subscription.satellite.norad_id)
    elif subscription.satellite.satellite_name is not None:
        satellite = pb2.SatelliteIdentifier(satellite_name=subscription.satellite.satellite_name)

    observer = pb2.GeodeticInput()

    if subscription.observer.lat_deg is not None:
        observer.lat_deg = subscription.observer.lat_deg
    elif subscription.observer.lat_rad is not None:
        observer.lat_rad = subscription.observer.lat_rad

    if subscription.observer.lon_deg is not None:
        observer.lon_deg = subscription.observer.lon_deg
    elif subscription.observer.lon_rad is not None:
        observer.lon_rad = subscription.observer.lon_rad

    if subscription.observer.alt_m is not None:
        observer.alt_m = subscription.observer.alt_m
    elif subscription.observer.alt_km is not None:
        observer.alt_km = subscription.observer.alt_km

    return pb2.Subscription(
        id=subscription.id or "",
        user_id=subscription.user_id,

        satellite=satellite,
        observer=observer,

        enabled=subscription.enabled,

        notify_before_seconds=subscription.notify_before_seconds,

        min_peak_elevation_deg=subscription.min_peak_elevation_deg,
        min_elevation_deg=subscription.min_elevation_deg,

        lookahead_days=subscription.lookahead_days,

        created_at=created,
        updated_at=updated,
    )
