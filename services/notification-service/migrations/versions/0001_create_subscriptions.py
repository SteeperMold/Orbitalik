"""create subscriptions table

Revision ID: 0001
Revises:
Create Date: 2026-06-29

"""

import sqlalchemy as sa
from alembic import op

revision = "0001"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "subscriptions",
        sa.Column("id", sa.BigInteger(), primary_key=True, autoincrement=True),
        sa.Column("user_id", sa.BigInteger(), nullable=False),
        sa.Column("norad_id", sa.BigInteger(), nullable=True),
        sa.Column("satellite_name", sa.String(), nullable=True),
        sa.Column("observer_lat_deg", sa.Float(), nullable=True),
        sa.Column("observer_lat_rad", sa.Float(), nullable=True),
        sa.Column("observer_lon_deg", sa.Float(), nullable=True),
        sa.Column("observer_lon_rad", sa.Float(), nullable=True),
        sa.Column("observer_alt_m", sa.Float(), nullable=True),
        sa.Column("observer_alt_km", sa.Float(), nullable=True),
        sa.Column("enabled", sa.Boolean(), nullable=False, server_default=sa.true()),
        sa.Column("notify_before_seconds", sa.Integer(), nullable=False),
        sa.Column("min_peak_elevation_deg", sa.Float(), nullable=True),
        sa.Column("min_peak_elevation_rad", sa.Float(), nullable=True),
        sa.Column("min_elevation_deg", sa.Float(), nullable=True),
        sa.Column("min_elevation_rad", sa.Float(), nullable=True),
        sa.Column("lookahead_days", sa.Integer(), nullable=False),
        sa.Column("created_at", sa.DateTime(), nullable=False),
        sa.Column("updated_at", sa.DateTime(), nullable=False),
    )

    op.create_index(
        "ix_subscriptions_user_id",
        "subscriptions",
        ["user_id"],
        unique=False,
    )


def downgrade() -> None:
    op.drop_index("ix_subscriptions_user_id", table_name="subscriptions")
    op.drop_table("subscriptions")
