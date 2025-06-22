FROM python:3.12-slim-bookworm

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# Copy the project into the image
COPY services/* /app/services/
COPY utils/* /app/utils/
COPY main.py /app
COPY pyproject.toml /app
COPY uv.lock /app

# Sync the project into a new environment, asserting the lockfile is up to date
WORKDIR /app

RUN touch .env # empty env file
RUN uv sync --locked

ENTRYPOINT ["uv", "run", "main.py"]
