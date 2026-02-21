import React, { useEffect, useState } from 'react';
import { Combobox, type ComboboxOption, InlineField, InlineFieldRow, Input, Button, IconButton } from '@grafana/ui';
import { DataSource } from '../../datasource';
import { WhereCondition, ColumnInfo } from '../../types';

interface Props {
  datasource: DataSource;
  table: string;
  value: WhereCondition[];
  onChange: (conditions: WhereCondition[]) => void;
}

const operatorOptions: Array<ComboboxOption<string>> = [
  { label: '=', value: '=' },
  { label: '!=', value: '!=' },
  { label: '<', value: '<' },
  { label: '>', value: '>' },
  { label: '<=', value: '<=' },
  { label: '>=', value: '>=' },
  { label: 'LIKE', value: 'LIKE' },
  { label: 'IN', value: 'IN' },
  { label: 'IS NULL', value: 'IS NULL' },
  { label: 'IS NOT NULL', value: 'IS NOT NULL' },
];

const noValueOperators = ['IS NULL', 'IS NOT NULL'];

export function WhereEditor({ datasource, table, value, onChange }: Props) {
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

  const addCondition = () => {
    onChange([...value, { column: '', operator: '=', value: '' }]);
  };

  const removeCondition = (idx: number) => {
    onChange(value.filter((_, i) => i !== idx));
  };

  const updateCondition = (idx: number, field: keyof WhereCondition, val: string) => {
    const updated = [...value];
    updated[idx] = { ...updated[idx], [field]: val };
    onChange(updated);
  };

  return (
    <>
      <InlineField label="WHERE" labelWidth={12}>
        <Button variant="secondary" size="sm" onClick={addCondition}>
          + Add condition
        </Button>
      </InlineField>
      {value.map((cond, idx) => (
        <InlineFieldRow key={idx}>
          <InlineField label={idx === 0 ? '' : 'AND'} labelWidth={12}>
            <Combobox
              options={columnOptions}
              value={cond.column || null}
              onChange={(option) => updateCondition(idx, 'column', option.value)}
              placeholder="Column"
              width={20}
            />
          </InlineField>
          <Combobox
            options={operatorOptions}
            value={cond.operator}
            onChange={(option) => updateCondition(idx, 'operator', option.value)}
            width={15}
          />
          {!noValueOperators.includes(cond.operator) && (
            <Input
              value={cond.value}
              onChange={(e) => updateCondition(idx, 'value', e.currentTarget.value)}
              placeholder="Value"
              width={20}
            />
          )}
          <IconButton name="trash-alt" tooltip="Remove condition" onClick={() => removeCondition(idx)} />
        </InlineFieldRow>
      ))}
    </>
  );
}
