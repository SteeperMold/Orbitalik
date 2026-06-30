from collections.abc import Awaitable, Callable

import grpc
import structlog

logger = structlog.get_logger()


class LoggingInterceptor(grpc.aio.ServerInterceptor):
    async def intercept_service(
        self,
        continuation: Callable[[grpc.HandlerCallDetails], Awaitable[grpc.RpcMethodHandler | None]],
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler | None:
        logger.info(
            "grpc_request",
            method=handler_call_details.method,
        )

        return await continuation(handler_call_details)
