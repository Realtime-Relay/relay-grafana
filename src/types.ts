import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface QueryInput extends DataQuery {
  topic?: string
  start_time?: string
}

export const DEFAULT_QUERY: Partial<QueryInput> = {
  topic: "",
  start_time: "${__from:date:iso}"
};

export interface DataPoint {
  Value: object;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
  secretKey?: string;
}
