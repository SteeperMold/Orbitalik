"""create notification jobs table

Revision ID: 0003
Revises: 0002
Create Date: 2026-07-01
"""

import sqlalchemy as sa
from alembic import op

revision = "0003"
down_revision = "0002"
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "notification_jobs",
        sa.Column("id", sa.BigInteger(), primary_key=True, autoincrement=True),
        sa.Column(
            "subscription_id",
            sa.BigInteger(),
            sa.ForeignKey("subscriptions.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("status", sa.SmallInteger(), nullable=False, server_default="0"),  # pending
        sa.Column("scheduled_time", sa.DateTime(timezone=True), nullable=False),
        sa.Column("aos", sa.DateTime(timezone=True), nullable=False),
        sa.Column("los", sa.DateTime(timezone=True), nullable=False),
        sa.Column("max_elevation_time", sa.DateTime(timezone=True), nullable=False),
        sa.Column("max_elevation", sa.Float(), nullable=False),
        sa.Column("attempts", sa.Integer(), nullable=False, server_default="0"),
        sa.Column("last_error", sa.Text(), nullable=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
    )

    op.create_index(
        "ix_notification_jobs_status_scheduled",
        "notification_jobs",
        ["status", "scheduled_time"],
        if_not_exists=True,
    )

    op.create_index(
        "ix_notification_jobs_subscription_id",
        "notification_jobs",
        ["subscription_id"],
        if_not_exists=True,
    )

    op.create_unique_constraint(
        "uq_notification_jobs_subscription_aos",
        "notification_jobs",
        ["subscription_id", "aos"],
    )


def downgrade() -> None:
    op.drop_constraint(
        "uq_notification_jobs_subscription_aos",
        "notification_jobs",
        type_="unique",
    )

    op.drop_index("ix_notification_jobs_subscription_id")
    op.drop_index("ix_notification_jobs_status_scheduled")

    op.drop_table("notification_jobs")
