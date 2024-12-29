import { DataSourceInstanceSettings, CoreApp, ScopedVars, DataQueryRequest, DataQueryResponse, LiveChannelScope } from '@grafana/data';
import { DataSourceWithBackend, getGrafanaLiveSrv, getTemplateSrv } from '@grafana/runtime';
import { Observable, merge } from 'rxjs';

import { QueryInput, MyDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<QueryInput, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<QueryInput> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: QueryInput, scopedVars: ScopedVars) {
    return {
      ...query,
      queryText: getTemplateSrv().replace(query.topic, scopedVars),
    };
  }

  filterQuery(query: QueryInput): boolean {
    // if no query has been provided, prevent the query from being executed
    return !!query.topic;
  }

  query(request: DataQueryRequest<QueryInput>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query, index) => {
      console.log(`Query => ${query}`)

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: this.uid,
          data: {
            ...query,
          },
        },
      });
    });

    return merge(...observables);
  }
}

