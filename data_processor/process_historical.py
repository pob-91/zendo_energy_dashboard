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


async def process_historical():
    weather_data = await data_service.get_latest_weather()
    energy_data = await data_service.get_latest_energy()

    # iterate through historical energy data and create aggregated correlations
    # there is a flaw here in that we assume that weather and energy data timestamps are quite well matched
    data_points = []

    count = 0
    for d in energy_data:
        if count >= len(weather_data):
            # done
            break

        # each time just knock off the most recent
        energy_slice = energy_data[count:]
        if len(energy_slice) < 2:
            break
        weather_slice = weather_data[count:]
        if len(weather_slice) <= 2:
            break

        sol_vs_prod = (
            correlation_service.calculate_solar_radiance_production_correlation(
                weather_data=weather_slice, energy_data=energy_slice
            )
        )
        temp_vs_consumption = (
            correlation_service.calculate_temp_consumption_correlation(
                weather_data=weather_slice, energy_data=energy_slice
            )
        )
        data_point = create_data_point(
            latest_weather=weather_slice[0],
            latest_energy=energy_slice[0],
            sol_vs_prod=sol_vs_prod,
            temp_vs_consumption=temp_vs_consumption,
        )
        data_points.append(data_point)

        count += 1

    # create docs so have something to work with
    success = await data_service.post_new_data_points(data_points=data_points)
    if not success:
        logger.error("Failed to create new data points, inspect return code")


if __name__ == "__main__":
    asyncio.run(process_historical())
