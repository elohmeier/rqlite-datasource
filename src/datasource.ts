import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { RqliteQuery, RqliteDataSourceOptions, DEFAULT_QUERY, ColumnInfo } from './types';
import { RqliteVariableSupport } from './variables';

export class DataSource extends DataSourceWithBackend<RqliteQuery, RqliteDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RqliteDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new RqliteVariableSupport();
  }

  getDefaultQuery(_: CoreApp): Partial<RqliteQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: RqliteQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      rawSql: getTemplateSrv().replace(query.rawSql, scopedVars),
    };
  }

  filterQuery(query: RqliteQuery): boolean {
    return !!query.rawSql;
  }

  async getTables(): Promise<string[]> {
    return this.getResource('/tables');
  }

  async getColumns(table: string): Promise<ColumnInfo[]> {
    return this.getResource('/columns', { table });
  }
}
