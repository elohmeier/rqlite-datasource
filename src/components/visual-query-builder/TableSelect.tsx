import React, { useEffect, useState } from 'react';
import { Combobox, type ComboboxOption, InlineField } from '@grafana/ui';
import { DataSource } from '../../datasource';

interface Props {
  datasource: DataSource;
  value: string;
  onChange: (table: string) => void;
}

export function TableSelect({ datasource, value, onChange }: Props) {
  const [tables, setTables] = useState<Array<ComboboxOption<string>>>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    datasource
      .getTables()
      .then((result) => {
        if (!cancelled) {
          setTables(result.map((t) => ({ label: t, value: t })));
        }
      })
      .catch(() => {
        if (!cancelled) {
          setTables([]);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [datasource]);

  return (
    <InlineField label="Table" labelWidth={12}>
      <Combobox
        options={tables}
        value={value || null}
        onChange={(option) => onChange(option?.value || '')}
        loading={loading}
        placeholder="Select table"
        width={30}
        isClearable
      />
    </InlineField>
  );
}
