import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { RqliteQuery, RqliteDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, RqliteQuery, RqliteDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
