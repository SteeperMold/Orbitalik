import uuid

try:
    # Python 3.14+
    from uuid import uuid7 as _uuid7
except ImportError:  # pragma: no cover - exercised on Python < 3.14
    from uuid6 import uuid7 as _uuid7


def new_uuid7() -> uuid.UUID:
    """Generate an internal entity identifier using UUIDv7."""
    return _uuid7()


def parse_uuid(value: str, field: str) -> uuid.UUID:
    """Parse a UUID received at a transport boundary with a clear error message."""
    if not value:
        raise ValueError(f"{field} is required")

    try:
        parsed = uuid.UUID(value)
    except (AttributeError, TypeError, ValueError) as exc:
        raise ValueError(f"invalid {field}") from exc

    if parsed.int == 0:
        raise ValueError(f"invalid {field}")

    return parsed
