import { VariableSupportType } from '@grafana/data';

import { DataSource } from './datasource';
import { DEFAULT_QUERY } from './types';

describe('DataSource variable support', () => {
  const createDataSource = () => {
    const settings = {
      id: 1,
      uid: 'rqlite-test',
      type: 'g42-rqlite-datasource',
      name: 'Rqlite',
      url: 'http://localhost:4001',
      jsonData: {},
    };

    return new DataSource(settings as any);
  };

  it('registers datasource variable support', () => {
    const ds = createDataSource();

    expect(ds.variables).toBeDefined();
    expect(ds.variables?.getType()).toBe(VariableSupportType.Datasource);
  });

  it('exposes table defaults for variable queries', () => {
    const ds = createDataSource();
    const query = ds.variables?.getDefaultQuery?.();

    expect(query).toEqual({
      ...DEFAULT_QUERY,
      format: 'table',
    });
  });
});
