import { DataSourceInstanceSettings, CoreApp, ScopedVars, DataQueryRequest, DataQueryResponse, LiveChannelScope } from '@grafana/data';
import { DataSourceWithBackend, getGrafanaLiveSrv, getTemplateSrv, logInfo } from '@grafana/runtime';
import { Observable, merge } from 'rxjs';
// import { crypto } from 'crypto';

import { QueryInput, MyDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<QueryInput, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<QueryInput> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: QueryInput, scopedVars: ScopedVars) {
    var start = getTemplateSrv().replace("${__from:date:iso}", scopedVars)
    var topic = getTemplateSrv().replace(query.topic, scopedVars)

    const date = new Date(start);

    // Convert to Unix timestamp (in seconds)
    const unixTimestamp = Math.floor(date.getTime() / 1000);

    // Adding time to create unique path
    var pathData = unixTimestamp

    return {
      ...query,
      topic: topic,
      start_time: start,
      path: pathData
    };
  }

  filterQuery(query: QueryInput): boolean {
    // if no query has been provided, prevent the query from being executed
    return true;
  }

  query(request: DataQueryRequest<QueryInput>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query, index) => {

      // To apply scoped vars
      const finalQuery = this.applyTemplateVariables(query, request.scopedVars);

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: finalQuery.path,
          data: {
            ...finalQuery,
          },
        },
      });
    });

    return merge(...observables);
  }
}

