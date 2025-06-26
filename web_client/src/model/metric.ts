export class WeatherData {
  // properties
  temperature: number;
  radiation: number;
  cloudCoverPercent: number;
  windSpeed: number;

  // constructors
  constructor(
    temperature: number,
    radiation: number,
    cloudCoverPercent: number,
    windSpeed: number
  ) {
    this.temperature = temperature;
    this.radiation = radiation;
    this.cloudCoverPercent = cloudCoverPercent;
    this.windSpeed = windSpeed;
  }

  static fromJson(json: any): WeatherData {
    return new WeatherData(
      json['temperature_2m'],
      json['direct_radiation'],
      json['cloud_cover'],
      json['wind_speed_10m']
    );
  }
}

export class PowerProductionBreakdown {
  nuclear?: number;
  geothermal?: number;
  biomass?: number;
  coal?: number;
  wind?: number;
  solar?: number;
  hydro?: number;
  gas?: number;
  oil?: number;
  unknown?: number;
  hydroDischarge?: number;
  batteryDischarge?: number;

  constructor(data: Partial<PowerProductionBreakdown>) {
    Object.assign(this, data);
  }

  static fromJson(json: any): PowerProductionBreakdown {
    return new PowerProductionBreakdown({
      nuclear: json['nuclear'],
      geothermal: json['geothermal'],
      biomass: json['biomass'],
      coal: json['coal'],
      wind: json['wind'],
      solar: json['solar'],
      hydro: json['hydro'],
      gas: json['gas'],
      oil: json['oil'],
      unknown: json['unknown'],
      hydroDischarge: json['hydro discharge'],
      batteryDischarge: json['battery discharge']
    });
  }
}

export class PowerConsumptionBreakdown {
  nuclear?: number;
  geothermal?: number;
  biomass?: number;
  coal?: number;
  wind?: number;
  solar?: number;
  hydro?: number;
  gas?: number;
  oil?: number;
  unknown?: number;
  hydroDischarge?: number;
  batteryDischarge?: number;

  constructor(data: Partial<PowerConsumptionBreakdown>) {
    Object.assign(this, data);
  }

  static fromJson(json: any): PowerConsumptionBreakdown {
    return new PowerConsumptionBreakdown({
      nuclear: json['nuclear'],
      geothermal: json['geothermal'],
      biomass: json['biomass'],
      coal: json['coal'],
      wind: json['wind'],
      solar: json['solar'],
      hydro: json['hydro'],
      gas: json['gas'],
      oil: json['oil'],
      unknown: json['unknown'],
      hydroDischarge: json['hydro discharge'],
      batteryDischarge: json['battery discharge']
    });
  }
}

export class CorrelationData {
  solarIrradianceVsSolarProductionCorrelation: number | null;
  temperatureVsConsumptionCorrelation: number | null;

  constructor(
    solarIrradianceVsSolarProductionCorrelation: number | null,
    temperatureVsConsumptionCorrelation: number | null
  ) {
    this.solarIrradianceVsSolarProductionCorrelation =
      solarIrradianceVsSolarProductionCorrelation;
    this.temperatureVsConsumptionCorrelation =
      temperatureVsConsumptionCorrelation;
  }

  static fromJson(json: any): CorrelationData {
    return new CorrelationData(
      json['solar_irradiance_vs_solar_production_correlation'],
      json['temperature_vs_consumption_correlation']
    );
  }
}

export class Metric {
  // properties
  timestamp: Date;
  totalProduction: number;
  totalConsumption: number;
  netBalance: number;
  weatherData: WeatherData;
  powerProductionData: PowerProductionBreakdown;
  powerConsumptionData: PowerConsumptionBreakdown;
  correlations: CorrelationData;

  constructor(
    timestamp: Date,
    totalProduction: number,
    totalConsumption: number,
    netBalance: number,
    weatherData: WeatherData,
    powerProductionData: PowerProductionBreakdown,
    powerConsumptionData: PowerConsumptionBreakdown,
    correlations: CorrelationData
  ) {
    this.timestamp = timestamp;
    this.totalProduction = totalProduction;
    this.totalConsumption = totalConsumption;
    this.netBalance = netBalance;
    this.weatherData = weatherData;
    this.powerProductionData = powerProductionData;
    this.powerConsumptionData = powerConsumptionData;
    this.correlations = correlations;
  }

  static fromJson(json: any): Metric {
    return new Metric(
      new Date(json['timestamp']),
      json['totalProduction'],
      json['totalConsumption'],
      json['netBalance'],
      WeatherData.fromJson(json['weatherData']),
      PowerProductionBreakdown.fromJson(json['powerProductionData']),
      PowerConsumptionBreakdown.fromJson(json['powerConsumptionData']),
      CorrelationData.fromJson(json['correlations'])
    );
  }
}
