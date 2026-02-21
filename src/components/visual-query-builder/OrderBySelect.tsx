import React, { useEffect, useState } from 'react';
import { Combobox, type ComboboxOption, InlineField, InlineFieldRow, Button, IconButton } from '@grafana/ui';
import { DataSource } from '../../datasource';
import { OrderByClause, ColumnInfo } from '../../types';

interface Props {
  datasource: DataSource;
  table: string;
  value: OrderByClause[];
  onChange: (orderBy: OrderByClause[]) => void;
}

const directionOptions: Array<ComboboxOption<string>> = [
  { label: 'ASC', value: 'ASC' },
  { label: 'DESC', value: 'DESC' },
];

export function OrderBySelect({ datasource, table, value, onChange }: Props) {
  const [availableColumns, setAvailableColumns] = useState<ColumnInfo[]>([]);

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
        }
      })
      .catch(() => {
        if (!cancelled) {
          setAvailableColumns([]);
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

  const addOrderBy = () => {
    onChange([...value, { column: '', direction: 'ASC' }]);
  };

  const removeOrderBy = (idx: number) => {
    onChange(value.filter((_, i) => i !== idx));
  };

  const updateOrderBy = (idx: number, field: keyof OrderByClause, val: string) => {
    const updated = [...value];
    updated[idx] = { ...updated[idx], [field]: val } as OrderByClause;
    onChange(updated);
  };

  return (
    <>
      <InlineField label="ORDER BY" labelWidth={12}>
        <Button variant="secondary" size="sm" onClick={addOrderBy}>
          + Add
        </Button>
      </InlineField>
      {value.map((ob, idx) => (
        <InlineFieldRow key={idx}>
          <InlineField label="" labelWidth={12}>
            <Combobox
              options={columnOptions}
              value={ob.column || null}
              onChange={(option) => updateOrderBy(idx, 'column', option.value)}
              placeholder="Column"
              width={20}
            />
          </InlineField>
          <Combobox
            options={directionOptions}
            value={ob.direction}
            onChange={(option) => updateOrderBy(idx, 'direction', option.value)}
            width={10}
          />
          <IconButton name="trash-alt" tooltip="Remove" onClick={() => removeOrderBy(idx)} />
        </InlineFieldRow>
      ))}
    </>
  );
}
