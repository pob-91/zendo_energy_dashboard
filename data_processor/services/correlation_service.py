import logging
import pandas as pd

logger = logging.getLogger(__name__)


class CorrelationService:
    def calculate_solar_radiance_production_correlation(
        self, weather_data, energy_data
    ) -> float | None:
        # create data frames
        mapped_weather_data = {
            "timestamp": [d["timestamp"] for d in weather_data],
            "radiation": [d["current"]["direct_radiation"] for d in weather_data],
        }
        mapped_energy_data = {
            "timestamp": [d["timestamp"] for d in energy_data],
            "solar": [d["powerProductionBreakdown"]["solar"] for d in energy_data],
        }
        weather_df = pd.DataFrame(mapped_weather_data)
        energy_df = pd.DataFrame(mapped_energy_data)

        # convert date col and index
        weather_df["timestamp"] = pd.to_datetime(weather_df["timestamp"])
        energy_df["timestamp"] = pd.to_datetime(energy_df["timestamp"])
        weather_df.set_index("timestamp", inplace=True)
        energy_df.set_index("timestamp", inplace=True)

        # there should not be duplicates but handle anyway (take mean)
        weather_df = weather_df.groupby(weather_df.index).mean()
        energy_df = energy_df.groupby(energy_df.index).mean()

        # interpolate to hourly and merge
        resampled_weather = weather_df.resample("1h").interpolate()
        resampled_energy = energy_df.resample("1h").interpolate()
        merged = resampled_weather.join(resampled_energy, how="inner")

        if len(merged) < 2:
            logger.info("Not enough data to calculate correlation")
            return None

        # calculate correlation
        pearson_corr = merged["radiation"].corr(merged["solar"])

        return pearson_corr

    def calculate_temp_consumption_correlation(
        self, weather_data, energy_data
    ) -> float | None:
        # create data frames
        mapped_weather_data = {
            "timestamp": [d["timestamp"] for d in weather_data],
            "temperature": [d["current"]["temperature_2m"] for d in weather_data],
        }
        mapped_energy_data = {
            "timestamp": [d["timestamp"] for d in energy_data],
            "consumption": [d["powerConsumptionTotal"] for d in energy_data],
        }
        weather_df = pd.DataFrame(mapped_weather_data)
        energy_df = pd.DataFrame(mapped_energy_data)

        # convert date col and index
        weather_df["timestamp"] = pd.to_datetime(weather_df["timestamp"])
        energy_df["timestamp"] = pd.to_datetime(energy_df["timestamp"])
        weather_df.set_index("timestamp", inplace=True)
        energy_df.set_index("timestamp", inplace=True)

        # there should not be duplicates but handle anyway (take mean)
        weather_df = weather_df.groupby(weather_df.index).mean()
        energy_df = energy_df.groupby(energy_df.index).mean()

        # interpolate to hourly and merge
        resampled_weather = weather_df.resample("1h").interpolate()
        resampled_energy = energy_df.resample("1h").interpolate()
        merged = resampled_weather.join(resampled_energy, how="inner")

        if len(merged) < 2:
            logger.info("Not enough data to calculate correlation")
            return None

        # calculate correlation
        pearson_corr = merged["temperature"].corr(merged["consumption"])

        return pearson_corr
