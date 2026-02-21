import React, { useEffect, useState } from 'react';
import { Combobox, type ComboboxOption, InlineField, InlineFieldRow, IconButton } from '@grafana/ui';
import { DataSource } from '../../datasource';
import { ColumnSelection, ColumnInfo } from '../../types';

interface Props {
  datasource: DataSource;
  table: string;
  value: ColumnSelection[];
  onChange: (columns: ColumnSelection[]) => void;
}

const aggregationOptions: Array<ComboboxOption<string>> = [
  { label: 'None', value: '' },
  { label: 'COUNT', value: 'COUNT' },
  { label: 'SUM', value: 'SUM' },
  { label: 'AVG', value: 'AVG' },
  { label: 'MIN', value: 'MIN' },
  { label: 'MAX', value: 'MAX' },
];

export function ColumnSelect({ datasource, table, value, onChange }: Props) {
  const [availableColumns, setAvailableColumns] = useState<ColumnInfo[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!table) {
      return;
    }
    let cancelled = false;
    datasource
      .getColumns(table)
      .then((cols) => {
        if (!cancelled) {
          setAvailableColumns(cols);
          setLoading(false);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setAvailableColumns([]);
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [datasource, table]);

  const columnOptions: Array<ComboboxOption<string>> = availableColumns.map((c) => ({
    label: `${c.name} (${c.type})`,
    value: c.name,
  }));

  const addColumn = (option: ComboboxOption<string>) => {
    if (option.value && !value.find((c) => c.name === option.value)) {
      onChange([...value, { name: option.value, aggregation: '' }]);
    }
  };

  const removeColumn = (idx: number) => {
    onChange(value.filter((_, i) => i !== idx));
  };

  const updateAggregation = (idx: number, agg: string) => {
    const updated = [...value];
    updated[idx] = { ...updated[idx], aggregation: agg };
    onChange(updated);
  };

  return (
    <>
      <InlineField label="Columns" labelWidth={12}>
        <Combobox
          options={columnOptions}
          onChange={addColumn}
          loading={loading}
          placeholder="Add column"
          width={30}
          value={null}
        />
      </InlineField>
      {value.map((col, idx) => (
        <InlineFieldRow key={col.name}>
          <InlineField label={col.name} labelWidth={12}>
            <Combobox
              options={aggregationOptions}
              value={col.aggregation}
              onChange={(option) => updateAggregation(idx, option.value)}
              width={15}
            />
          </InlineField>
          <IconButton name="trash-alt" tooltip="Remove column" onClick={() => removeColumn(idx)} />
        </InlineFieldRow>
      ))}
    </>
  );
}
