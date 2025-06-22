def create_data_point(
    latest_weather, latest_energy, sol_vs_prod, temp_vs_consumption
) -> dict:
    return {
        "type": "AGGREGATED_DATA",
        "timestamp": latest_energy["timestamp"],
        "totalProduction": latest_energy["powerProductionTotal"],
        "totalConsumption": latest_energy["powerConsumptionTotal"],
        "netBalance": latest_energy["powerProductionTotal"]
        - latest_energy["powerConsumptionTotal"],
        "weatherData": {**latest_weather["current"]},
        "energy_data": {
            "solar_production": latest_energy["powerProductionBreakdown"]["solar"]
        },
        "correlations": {
            "solar_irradiance_vs_solar_production_correlation": sol_vs_prod,
            "temperature_vs_consumption_correlation": temp_vs_consumption,
        },
    }
