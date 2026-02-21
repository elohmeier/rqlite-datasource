import React, { ChangeEvent } from 'react';
import { Combobox, type ComboboxOption, DataSourceHttpSettings, FieldSet, InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { RqliteDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<RqliteDataSourceOptions> {}

const consistencyOptions: Array<ComboboxOption<string>> = [
  { label: 'None', value: 'none', description: 'No consistency guarantee' },
  { label: 'Weak (default)', value: 'weak', description: 'Weak consistency - reads from leader' },
  { label: 'Linearizable', value: 'linearizable', description: 'Linearizable reads' },
  { label: 'Strong', value: 'strong', description: 'Strong consistency - Raft consensus' },
];

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData } = options;

  const onConsistencyChange = (option: ComboboxOption<string>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        consistencyLevel: option.value,
      },
    });
  };

  const onTimeoutChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        timeout: event.target.value,
      },
    });
  };

  return (
    <>
      <DataSourceHttpSettings
        defaultUrl="http://localhost:4001"
        dataSourceConfig={options}
        onChange={onOptionsChange}
      />
      <FieldSet label="rqlite Settings">
        <InlineField label="Consistency Level" labelWidth={20} tooltip="Read consistency level for rqlite queries">
          <Combobox
            options={consistencyOptions}
            value={jsonData.consistencyLevel || 'weak'}
            onChange={onConsistencyChange}
            width={30}
          />
        </InlineField>
        <InlineField label="Query Timeout" labelWidth={20} tooltip="Query timeout (e.g. 10s, 30s)">
          <Input
            value={jsonData.timeout || ''}
            onChange={onTimeoutChange}
            placeholder="10s"
            width={30}
          />
        </InlineField>
      </FieldSet>
    </>
  );
}
