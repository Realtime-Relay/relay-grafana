import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface QueryInput extends DataQuery {
  topic?: string
}

export const DEFAULT_QUERY: Partial<QueryInput> = {
  topic: "test-topic",
};

export interface DataPoint {
  Time: object;
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
  username: string;
  password: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
