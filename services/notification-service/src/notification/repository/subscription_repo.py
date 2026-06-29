from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from notification.infrastructure.db.models.subscription import SubscriptionORM


class SubscriptionRepository:
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]) -> None:
        self.session_factory = session_factory

    async def create(self, orm: SubscriptionORM) -> SubscriptionORM:
        async with self.session_factory() as session:
            session.add(orm)
            await session.commit()
            await session.refresh(orm)
            return orm
