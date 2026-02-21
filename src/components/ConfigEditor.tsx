import React, { ChangeEvent } from 'react';
import {
  AdvancedHttpSettings,
  Auth,
  ConfigSection,
  ConfigSubSection,
  ConnectionSettings,
  convertLegacyAuthProps,
} from '@grafana/plugin-ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Combobox, type ComboboxOption, Divider, InlineField, Input, Stack } from '@grafana/ui';
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
      <ConnectionSettings config={options} onChange={onOptionsChange} urlPlaceholder="http://localhost:4001" />

      <Divider spacing={4} />
      <Auth
        {...convertLegacyAuthProps({
          config: options,
          onChange: onOptionsChange,
        })}
      />

      <Divider spacing={4} />
      <ConfigSection
        title="Additional settings"
        description="Additional settings are optional settings that can be configured for more control over your data source."
        isCollapsible={true}
        isInitiallyOpen={false}
      >
        <Stack gap={5} direction="column">
          <AdvancedHttpSettings config={options} onChange={onOptionsChange} />

          <ConfigSubSection title="rqlite Settings">
            <InlineField label="Consistency Level" labelWidth={20} tooltip="Read consistency level for rqlite queries">
              <Combobox
                options={consistencyOptions}
                value={jsonData.consistencyLevel || 'weak'}
                onChange={onConsistencyChange}
                width={30}
              />
            </InlineField>
            <InlineField label="Query Timeout" labelWidth={20} tooltip="Query timeout (e.g. 10s, 30s)">
              <Input value={jsonData.timeout || ''} onChange={onTimeoutChange} placeholder="10s" width={30} />
            </InlineField>
          </ConfigSubSection>
        </Stack>
      </ConfigSection>
    </>
  );
}
