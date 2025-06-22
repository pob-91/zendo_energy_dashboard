import os
import json
import aiohttp
import logging

from datetime import datetime, timedelta, timezone

logger = logging.getLogger(__name__)


class DataService:
    async def get_latest_weather(self) -> list[dict]:
        async with aiohttp.ClientSession() as session:
            async with session.get(DataService.latest_day_weather_url()) as resp:
                try:
                    content = await resp.content.read()
                    data = json.loads(content)
                    return [e["doc"] for e in data["rows"]]
                except json.JSONDecodeError as e:
                    logger.error(f"Failed to decode weather data: {e}")
                    return []

    async def get_latest_energy(self) -> list[dict]:
        async with aiohttp.ClientSession() as session:
            async with session.get(DataService.latest_day_energy_url()) as resp:
                try:
                    content = await resp.content.read()
                    data = json.loads(content)
                    return [e["doc"] for e in data["rows"]]
                except json.JSONDecodeError as e:
                    logger.error(f"Failed to decode energy data: {e}")
                    return []

    async def post_new_data_point(self, data_point) -> bool:
        async with aiohttp.ClientSession() as session:
            async with session.post(DataService._base_url(), json=data_point) as resp:
                return resp.status == 201

    async def post_new_data_points(self, data_points) -> bool:
        async with aiohttp.ClientSession() as session:
            async with session.post(
                DataService.bulk_upload_url(), json={"docs": data_points}
            ) as resp:
                logger.info(f"statis: {resp.status}")
                return resp.status == 201

    @staticmethod
    def energy_change_feed_url() -> str:
        return f"{DataService._base_url()}/_changes?feed=continuous&include_docs=false&since=now&filter=filters/only_energy_documents&type=ENERGY_DATA&heartbeat=10000"

    @staticmethod
    def latest_day_weather_url() -> str:
        now = datetime.now(timezone.utc)
        start_time = now - timedelta(hours=24)
        startkey = start_time.isoformat().replace("+00:00", "Z")
        return f'{DataService._base_url()}/_design/views/_view/weather_by_time?include_docs=true&descending=true&startKey="{startkey}"'

    @staticmethod
    def latest_day_energy_url() -> str:
        now = datetime.now(timezone.utc)
        start_time = now - timedelta(hours=24)
        startkey = start_time.isoformat().replace("+00:00", "Z")
        return f'{DataService._base_url()}/_design/views/_view/energy_by_time?include_docs=true&descending=true&startKey="{startkey}"'

    @staticmethod
    def bulk_upload_url() -> str:
        return f"{DataService._base_url()}/_bulk_docs"

    # private
    @staticmethod
    def _base_url() -> str:
        return f"http://{os.getenv('COUCHDB_USER')}:{os.getenv('COUCHDB_PASSWORD')}@{os.getenv('COUCHDB_URL')}/{os.getenv('COUCHDB_DB')}"
