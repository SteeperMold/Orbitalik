import grpc.aio

from notification.proto import notification_pb2 as pb2
from notification.services.notification_service import NotificationService
from notification.transport.grpc.converters import create_subscription_command_from_request, subscription_to_proto


class NotificationServicer:
    def __init__(self, service: NotificationService) -> None:
        self.service = service

    async def CreateSubscription(
            self,
            request: pb2.CreateSubscriptionRequest,
            context: grpc.aio.ServicerContext,
    ) -> pb2.Subscription:
        try:
            cmd = create_subscription_command_from_request(request)
            subscription = await self.service.create_subscription(cmd)

            return subscription_to_proto(subscription)
        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def GetSubscription(self, request, context):
        pass

    async def ListSubscriptions(self, request, context):
        pass

    async def UpdateSubscription(self, request, context):
        pass

    async def DeleteSubscription(self, request, context):
        pass

    async def SetSubscriptionStatus(self, request, context):
        pass

    async def RegisterDevice(self, request, context):
        pass

    async def ListDevices(self, request, context):
        pass

    async def DeleteDevice(self, request, context):
        pass
