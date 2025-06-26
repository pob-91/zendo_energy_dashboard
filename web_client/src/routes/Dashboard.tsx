import {
  type FunctionComponent,
  useEffect,
  useState,
  useCallback,
  useRef
} from 'react';

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Legend,
  Bar
} from 'recharts';

import { Metric } from '../model';
import { DataService } from '../services';

const Dashboard: FunctionComponent = () => {
  const [latestData, setLatestData] = useState<Metric | undefined>(undefined);
  const [historicalData, setHistoricalData] = useState<Metric[]>([]);

  const timeout = useRef<NodeJS.Timeout | null>(null);

  // callbacks
  const fetchData = useCallback(async () => {
    const result = await DataService.getLatestMetric();
    setLatestData(result.data);

    const historical = await DataService.getTimeSeriesData();
    setHistoricalData((historical.data || []).reverse());
  }, []);

  // lifecycle
  useEffect(() => {
    // TODO: Make this a web socket
    timeout.current = setInterval(fetchData, 60000 * 5); // poll every 5 mins

    fetchData();

    return () => {
      if (!timeout.current) {
        return;
      }
      clearInterval(timeout.current);
    };
  }, [fetchData]);

  // utils
  const formatDate = (tickItem: Date) => {
    if (!(tickItem instanceof Date)) {
      tickItem = new Date(tickItem); // in case tickItem is string or number
    }
    return tickItem.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // render
  if (!latestData) {
    return <div>Loading...</div>;
  }

  // render
  return (
    <div className="w-full h-svh">
      <div className="grid grid-cols-2 gap-4 p-8 bg-[#8C97ED33] w-full h-full overflow-y-auto">
        <div className="col-span-2 text-3xl text-center p-8 bg-[#001744] text-white font-bold rounded-lg shadow-md shadow-gray-400 border border-gray-100">
          Zendo Energy Dashboard
        </div>
        <div className="flex flex-col items-center justify-center p-8 bg-[#001744] rounded-lg shadow-md shadow-gray-400 border border-gray-100 text-[#C4F07D] font-bold text-xl">
          <div>
            {`Total Energy Production: ${latestData?.totalProduction} kw/h`}
          </div>
          <div>
            {`Total Energy Consumption: ${latestData?.totalConsumption} kw/h`}
          </div>
          <div>{`Net Balance: ${latestData?.netBalance} kw/h`}</div>
        </div>
        <div className="flex flex-col items-center justify-center bg-[#001744] p-8 rounded-lg shadow-md shadow-gray-400 border border-gray-100 text-[#C4F07D] font-bold text-xl">
          <div>{`Temperature: ${latestData?.weatherData.temperature} deg C`}</div>
          <div>
            {`Cloud Cover: ${latestData?.weatherData.cloudCoverPercent} %`}
          </div>
          <div>{`Wind Speed: ${latestData?.weatherData.windSpeed} km/h`}</div>
          <div>{`Solar Irradiance: ${latestData?.weatherData.radiation} W/m2`}</div>
        </div>
        {historicalData.length === 0 && <div>No historical data...</div>}
        {historicalData.length > 0 && (
          <>
            <div className="col-span-2 flex flex-col bg-white justify-center p-8 rounded-lg shadow-md shadow-gray-400 border border-gray-100 aspect-5/2">
              <div className="text-xl text-center mb-4">Net Balance</div>
              <ResponsiveContainer width="90%" height="100%">
                <LineChart data={historicalData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="timestamp" tickFormatter={formatDate} />
                  <YAxis />
                  <Tooltip />
                  <Line
                    type="monotone"
                    dataKey="netBalance"
                    stroke="#8884d8"
                    strokeWidth={4}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="flex flex-col bg-white justify-center p-8 rounded-lg shadow-md shadow-gray-400 border border-gray-100 aspect-video">
              <div className="text-xl text-center mb-4">
                Solar Radiation vs Production
              </div>
              <ResponsiveContainer width="90%" height="100%">
                <LineChart data={historicalData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="timestamp" tickFormatter={formatDate} />
                  <YAxis yAxisId="left" />
                  <YAxis yAxisId="right" orientation="right" />
                  <Tooltip />
                  <Line
                    type="monotone"
                    dataKey="weatherData.radiation"
                    stroke="#8884d8"
                    strokeWidth={4}
                    dot={false}
                    yAxisId="left"
                  />
                  <Line
                    type="monotone"
                    dataKey="powerProductionData.solar"
                    stroke="#ff6b35"
                    strokeWidth={4}
                    dot={false}
                    yAxisId="right"
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="flex flex-col justify-center p-8 bg-white rounded-lg shadow-md shadow-gray-400 border border-gray-100 aspect-video">
              <div className="text-xl text-center mb-4">
                Temperature vs Consumption
              </div>
              <ResponsiveContainer width="90%" height="100%">
                <LineChart data={historicalData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="timestamp" tickFormatter={formatDate} />
                  <YAxis yAxisId="left" />
                  <YAxis yAxisId="right" orientation="right" />
                  <Tooltip />
                  <Line
                    type="monotone"
                    dataKey="weatherData.temperature"
                    stroke="#8884d8"
                    strokeWidth={4}
                    dot={false}
                    yAxisId="left"
                  />
                  <Line
                    type="monotone"
                    dataKey="totalConsumption"
                    stroke="#ff6b35"
                    strokeWidth={4}
                    dot={false}
                    yAxisId="right"
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
            <div className="col-span-2 flex flex-col justify-center p-8 rounded-lg shadow-md shadow-gray-400 border border-gray-100 bg-white aspect-5/2">
              <div className="text-xl text-center mb-4">Correlations</div>
              <div className="flex">
                <div className="flex text-md font-bold text-center mb-4">
                  <div
                    className="mr-2"
                    style={{
                      transform: `rotate(${90 - 90 * (latestData?.correlations?.solarIrradianceVsSolarProductionCorrelation || 0)}deg)`
                    }}
                  >
                    &uarr;
                  </div>
                  {`Latest Irradiance VS Solar Production: ${latestData?.correlations?.solarIrradianceVsSolarProductionCorrelation?.toFixed(2)}`}
                </div>
                <div className="flex text-md font-bold text-center mb-4 ml-4">
                  <div
                    className="mr-2"
                    style={{
                      transform: `rotate(${90 - 90 * (latestData?.correlations?.temperatureVsConsumptionCorrelation || 0)}deg)`
                    }}
                  >
                    &uarr;
                  </div>
                  {`Latest Temperature VS Consumption: ${latestData?.correlations?.temperatureVsConsumptionCorrelation?.toFixed(2)}`}
                </div>
              </div>
              <ResponsiveContainer width="90%" height="100%">
                <BarChart data={historicalData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="timestamp" tickFormatter={formatDate} />
                  <YAxis yAxisId="left" />
                  <Tooltip />
                  <Legend
                    formatter={value => {
                      if (
                        value.includes(
                          'solarIrradianceVsSolarProductionCorrelation'
                        )
                      ) {
                        return 'Irradiance VS Solar Prouction';
                      }
                      if (
                        value.includes('temperatureVsConsumptionCorrelation')
                      ) {
                        return 'Temperature VS Consumption';
                      }
                      return value;
                    }}
                  />
                  <Bar
                    dataKey="correlations.solarIrradianceVsSolarProductionCorrelation"
                    fill="#8884d8"
                    yAxisId="left"
                  />
                  <Bar
                    dataKey="correlations.temperatureVsConsumptionCorrelation"
                    fill="#ff6b35"
                    yAxisId="left"
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
