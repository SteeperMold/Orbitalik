import sys

import logging
import structlog

from notification.infrastructure.settings import AppEnv

log_levels = {
    AppEnv.DEVELOPMENT: logging.DEBUG,
    AppEnv.TEST: logging.INFO,
    AppEnv.PRODUCTION: logging.INFO,
}


def configure_logging(app_env: AppEnv) -> None:
    level = log_levels.get(app_env)

    logging.basicConfig(
        level=level,
        stream=sys.stdout,
        format="%(message)s",
    )

    structlog.configure(
        processors=[
            structlog.contextvars.merge_contextvars,
            structlog.processors.add_log_level,
            structlog.processors.TimeStamper(fmt="iso"),
            structlog.processors.StackInfoRenderer(),
            structlog.processors.format_exc_info,
            structlog.processors.JSONRenderer(),
        ],
        wrapper_class=structlog.make_filtering_bound_logger(level),
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )
