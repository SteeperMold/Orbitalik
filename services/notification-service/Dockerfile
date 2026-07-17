FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:$PATH"

COPY pyproject.toml uv.lock* /app/

COPY migrations /app/migrations
COPY src /app/src

RUN uv sync --frozen

ENV PYTHONPATH=/app/src

CMD ["uv", "run", "python", "-m", "notification.main"]
