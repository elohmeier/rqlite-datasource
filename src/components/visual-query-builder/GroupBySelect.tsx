import React, { useEffect, useState } from 'react';
import { MultiCombobox, type ComboboxOption, InlineField } from '@grafana/ui';
import { DataSource } from '../../datasource';
import { ColumnInfo } from '../../types';

interface Props {
  datasource: DataSource;
  table: string;
  value: string[];
  onChange: (groupBy: string[]) => void;
}

export function GroupBySelect({ datasource, table, value, onChange }: Props) {
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
    label: c.name,
    value: c.name,
  }));

  const onSelectionChange = (items: Array<ComboboxOption<string>>) => {
    onChange(items.map((item) => item.value));
  };

  return (
    <InlineField label="GROUP BY" labelWidth={12}>
      <MultiCombobox
        options={columnOptions}
        value={value}
        onChange={onSelectionChange}
        loading={loading}
        placeholder="Select columns"
        width={30}
      />
    </InlineField>
  );
}
