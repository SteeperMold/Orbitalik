# ruff: noqa: N802

import dataclasses

import grpc.aio
from google.protobuf import empty_pb2

from notification.proto import notification_pb2 as pb2
from notification.services.device_service import DeviceService
from notification.services.subscription_service import SubscriptionService
from notification.transport.grpc import converters


@dataclasses.dataclass
class NotificationServicer:
    subscriptions: SubscriptionService
    devices: DeviceService

    async def CreateSubscription(
        self,
        request: pb2.CreateSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            cmd = converters.create_subscription_command_from_request(request)
            subscription = await self.subscriptions.create(cmd)

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def GetSubscription(
        self,
        request: pb2.GetSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            subscription = await self.subscriptions.get(request.id)

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    async def ListSubscriptions(
        self,
        request: pb2.ListSubscriptionsRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.ListSubscriptionsResponse:
        try:
            items, next_page_token = await self.subscriptions.list(
                user_id=request.user_id,
                enabled=request.enabled if request.HasField("enabled") else None,
                page_size=request.page_size,
                page_token=request.page_token if request.HasField("page_token") else None,
            )

            return pb2.ListSubscriptionsResponse(
                subscriptions=[converters.subscription_to_proto(s) for s in items],
                next_page_token=next_page_token,
            )

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def UpdateSubscription(
        self,
        request: pb2.UpdateSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            cmd = converters.update_subscription_command_from_request(request)

            subscription = await self.subscriptions.update(
                subscription_id=request.id,
                cmd=cmd,
            )

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def DeleteSubscription(
        self,
        request: pb2.DeleteSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        try:
            await self.subscriptions.delete(request.id)
            return empty_pb2.Empty()

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    async def SetSubscriptionStatus(
        self,
        request: pb2.SetSubscriptionStatusRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            subscription = await self.subscriptions.set_status(
                subscription_id=request.id,
                enabled=request.enabled,
            )

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    async def RegisterDevice(
        self,
        request: pb2.RegisterDeviceRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Device:
        try:
            device = await self.devices.register(
                user_id=request.user_id,
                device_type=request.type,
                address=request.address,
            )

            return converters.device_to_proto(device)

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def ListDevices(
        self,
        request: pb2.ListDevicesRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.ListDevicesResponse:
        devices = await self.devices.list(user_id=request.user_id)

        return pb2.ListDevicesResponse(devices=[converters.device_to_proto(d) for d in devices])

    async def DeleteDevice(
        self,
        request: pb2.DeleteDeviceRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        try:
            await self.devices.delete_device(request.id)
            return empty_pb2.Empty()

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))
