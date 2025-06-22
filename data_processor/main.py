import json
import aiohttp
import asyncio
import logging

from dotenv import load_dotenv

from services.data_service import DataService
from services.correlation_service import CorrelationService
from utils.data import create_data_point

# load env vars
load_dotenv()

# configure logger
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger(__name__)


data_service = DataService()
correlation_service = CorrelationService()


async def fetch_data_and_correlate():
    weather_data = await data_service.get_latest_weather()
    energy_data = await data_service.get_latest_energy()
    sol_vs_prod = correlation_service.calculate_solar_radiance_production_correlation(
        weather_data=weather_data, energy_data=energy_data
    )
    temp_vs_consumption = correlation_service.calculate_temp_consumption_correlation(
        weather_data=weather_data, energy_data=energy_data
    )

    # create a master document to drive the API
    # we are listening to energy data changes so use that as a master template
    data_point = create_data_point(
        latest_weather=weather_data[0],
        latest_energy=energy_data[0],
        sol_vs_prod=sol_vs_prod,
        temp_vs_consumption=temp_vs_consumption,
    )
    success = await data_service.post_new_data_point(data_point=data_point)
    if not success:
        logger.error("Failed to create new data point, inspect return code")


async def listen_for_energy_changes():
    async with aiohttp.ClientSession() as session:
        async with session.get(DataService.energy_change_feed_url()) as resp:
            async for raw_line in resp.content:
                line = ""
                if isinstance(raw_line, bytes):
                    line = raw_line.decode("utf-8").strip()
                elif isinstance(raw_line, str):
                    line = raw_line.strip()
                else:
                    logger.error(f"Could not handle change of type {type(raw_line)}")
                    continue

                if len(line) == 0:
                    logger.info("Hearbeat ping...")
                    continue

                try:
                    await fetch_data_and_correlate()
                except json.JSONDecodeError as e:
                    logger.error(f"Failed to decode change: {e}")
                    continue


if __name__ == "__main__":
    try:
        asyncio.run(listen_for_energy_changes())
    except KeyboardInterrupt:
        logger.info("Shutting down...")
    except Exception:
        raise
