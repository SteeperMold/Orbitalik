# ruff: noqa: N802

import dataclasses

import grpc.aio
from google.protobuf import empty_pb2

from notification.proto import notification_pb2 as pb2
from notification.services.subscription_service import SubscriptionService
from notification.transport.grpc import converters


@dataclasses.dataclass
class NotificationServicer:
    service: SubscriptionService

    async def CreateSubscription(
        self,
        request: pb2.CreateSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            cmd = converters.create_subscription_command_from_request(request)
            subscription = await self.service.create_subscription(cmd)

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def GetSubscription(
        self,
        request: pb2.GetSubscriptionRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            subscription = await self.service.get_subscription(request.id)

            return converters.subscription_to_proto(subscription)

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    async def ListSubscriptions(
        self,
        request: pb2.ListSubscriptionsRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.ListSubscriptionsResponse:
        try:
            items, next_page_token = await self.service.list_subscriptions(
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

            subscription = await self.service.update_subscription(
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
            await self.service.delete_subscription(request.id)
            return empty_pb2.Empty()

        except ValueError as e:
            await context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    async def SetSubscriptionStatus(
        self,
        request: pb2.SetSubscriptionStatusRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            subscription = await self.service.set_subscription_status(
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
            device = await self.service.register_device(
                user_id=request.user_id,
                token=request.token,
            )

            return pb2.Device(
                id=device.id,
                user_id=device.user_id,
                token=device.token,
            )

        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def ListDevices(
        self,
        request: pb2.ListDevicesRequest,
        context: grpc.aio.ServicerContext,
    ) -> pb2.ListDevicesResponse:
        devices = await self.service.list_devices(user_id=request.user_id)

        return pb2.ListDevicesResponse(
            devices=[
                pb2.Device(
                    id=d.id,
                    user_id=d.user_id,
                    token=d.token,
                )
                for d in devices
            ]
        )

    async def DeleteDevice(
        self,
        request: pb2.DeleteDeviceRequest,
        context: grpc.aio.ServicerContext,
    ) -> empty_pb2.Empty:
        await self.service.delete_device(request.id)
        return empty_pb2.Empty()
