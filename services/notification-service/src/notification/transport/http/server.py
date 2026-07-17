import uvicorn
from fastapi import FastAPI, Response
from prometheus_client import CONTENT_TYPE_LATEST, generate_latest

from notification.infrastructure.settings import Settings


class HTTPServer:
    def __init__(self, settings: Settings):
        self.settings = settings

        self.app = FastAPI()

        @self.app.get("/health")
        async def health() -> dict:
            return {"status": "ok"}

        @self.app.get("/metrics")
        async def metrics() -> Response:
            return Response(
                generate_latest(),
                media_type=CONTENT_TYPE_LATEST,
            )

        self._server = uvicorn.Server(
            uvicorn.Config(
                self.app,
                host="0.0.0.0",
                port=settings.http_port,
                log_level="info",
            )
        )

    async def start(self) -> None:
        await self._server.serve()
