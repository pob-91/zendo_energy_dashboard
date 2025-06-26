const apiUrl = import.meta.env.VITE_API_URL;

import { Metric, type HttpResponse } from '../model';

interface IDataService {
  getLatestMetric(): Promise<HttpResponse<Metric>>;
  getTimeSeriesData(): Promise<HttpResponse<Metric[]>>;
}

export const DataService: IDataService = {
  getLatestMetric,
  getTimeSeriesData
};

async function getLatestMetric(): Promise<HttpResponse<Metric>> {
  const url = `${apiUrl}/energy-summary`;

  try {
    const response = await fetch(url, {
      method: 'GET'
    });

    if (response.status !== 200) {
      return {
        success: false,
        code: response.status
      };
    }

    const json = await response.json();
    return {
      success: true,
      code: 200,
      data: Metric.fromJson(json)
    };
  } catch {
    return {
      success: false,
      code: 500
    };
  }
}

async function getTimeSeriesData(): Promise<HttpResponse<Metric[]>> {
  const url = `${apiUrl}/historical-data`;

  try {
    const response = await fetch(url, {
      method: 'GET'
    });

    if (response.status !== 200) {
      return {
        success: false,
        code: response.status
      };
    }

    const json = (await response.json()) as any[];
    return {
      success: true,
      code: 200,
      data: json.map(e => Metric.fromJson(e))
    };
  } catch {
    return {
      success: false,
      code: 500
    };
  }
}
