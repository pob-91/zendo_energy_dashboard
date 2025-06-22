import {
  type FunctionComponent,
  useEffect,
  useState,
  useMemo,
  useCallback,
  useRef,
} from "react";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  LabelList,
} from "recharts";

import { Metric } from "../model";
import { DataService } from "../services";

const Dashboard: FunctionComponent = () => {
  const [latestData, setLatestData] = useState<Metric | undefined>(undefined);
  const [historicalData, setHistoricalData] = useState<Metric[]>([]);

  const timeout = useRef<NodeJS.Timeout | null>(null);

  const correlationData = useMemo(() => {
    return [
      {
        name: "Solar vs Production",
        value:
          latestData?.correlations
            .solarIrradianceVsSolarProductionCorrelation || 0,
      },
      {
        name: "Temp vs Consumption",
        value:
          latestData?.correlations.temperatureVsConsumptionCorrelation || 0,
      },
    ];
  }, [latestData]);

  // callbacks
  const fetchData = useCallback(async () => {
    const result = await DataService.getLatestMetric();
    setLatestData(result.data);

    const historical = await DataService.getTimeSeriesData();
    setHistoricalData(historical.data || []);
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
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // render
  if (!latestData) {
    return <div>Loading...</div>;
  }

  // render
  return (
    <div className="w-svh h-svh grid grid-cols-2 gap-4">
      <div className="flex flex-col items-center justify-center">
        <div>
          {`Total Energy Production: ${latestData?.totalProduction} kw/h`}
        </div>
        <div>
          {`Total Energy Consumption: ${latestData?.totalConsumption} kw/h`}
        </div>
        <div>{`Net Balance: ${latestData?.netBalance} kw/h`}</div>
      </div>
      <div className="flex flex-col items-center justify-center">
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
          <div className="col-span-2">
            <ResponsiveContainer width="99%" height="100%">
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
          <div>
            <ResponsiveContainer width="99%" height="100%">
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
          <div>
            <ResponsiveContainer width="99%" height="100%">
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
          <div className="col-span-2">
            <ResponsiveContainer width="99%" height="100%">
              <BarChart
                data={correlationData}
                margin={{ top: 20, right: 30, left: 20, bottom: 40 }}
                barCategoryGap="30%"
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="name"
                  angle={-30}
                  textAnchor="end"
                  interval={0}
                  height={60}
                />
                <YAxis domain={[-1, 1]} />
                <Tooltip formatter={(value: number) => value.toFixed(2)} />
                <Bar
                  dataKey="value"
                  isAnimationActive={false}
                  radius={[5, 5, 0, 0]}
                  // dynamically set fill per bar
                  fill="#8884d8"
                >
                  <LabelList
                    dataKey="value"
                    position="top"
                    formatter={(value: number) => value.toFixed(2)}
                  />
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </>
      )}
    </div>
  );
};

export default Dashboard;
