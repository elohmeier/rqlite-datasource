import { DataSourceVariableSupport } from '@grafana/data';

import { DEFAULT_QUERY, RqliteQuery } from './types';
import type { DataSource } from './datasource';

export class RqliteVariableSupport extends DataSourceVariableSupport<DataSource, RqliteQuery> {
  getDefaultQuery(): Partial<RqliteQuery> {
    return {
      ...DEFAULT_QUERY,
      format: 'table',
    };
  }
}
