import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export type EditorMode = 'code' | 'builder';
export type QueryFormat = 'table' | 'time_series';

export interface ColumnSelection {
  name: string;
  aggregation: string; // '', 'COUNT', 'SUM', 'AVG', 'MIN', 'MAX'
}

export interface WhereCondition {
  column: string;
  operator: string; // '=', '!=', '<', '>', '<=', '>=', 'LIKE', 'IN', 'IS NULL', 'IS NOT NULL'
  value: string;
}

export interface OrderByClause {
  column: string;
  direction: 'ASC' | 'DESC';
}

export interface RqliteQuery extends DataQuery {
  rawSql: string;
  format: QueryFormat;
  timeColumns: string[];
  editorMode: EditorMode;

  // Visual builder fields
  table: string;
  columns: ColumnSelection[];
  whereClause: WhereCondition[];
  groupBy: string[];
  orderBy: OrderByClause[];
  limit: string;
  offset: string;
}

export const DEFAULT_QUERY: Partial<RqliteQuery> = {
  rawSql: '',
  format: 'table',
  timeColumns: ['time'],
  editorMode: 'code',
  table: '',
  columns: [],
  whereClause: [],
  groupBy: [],
  orderBy: [],
  limit: '',
  offset: '',
};

export interface RqliteDataSourceOptions extends DataSourceJsonData {
  consistencyLevel?: string;
  timeout?: string;
}

export interface ColumnInfo {
  name: string;
  type: string;
}
