from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession


def create_engine(db_url: str, app_env: str):
    return create_async_engine(
        db_url,
        echo=app_env == "development",
        pool_pre_ping=True,
    )


def create_session_factory(engine):
    return async_sessionmaker(
        engine,
        class_=AsyncSession,
        expire_on_commit=False,
    )
