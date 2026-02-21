import React, { useCallback, useEffect, useMemo, useState } from 'react';
import {
  type MonacoEditor,
  CodeEditor,
  InlineField,
  InlineFieldRow,
  Input,
  RadioButtonGroup,
  Combobox,
  type ComboboxOption,
  Collapse,
  Modal,
  Button,
  ConfirmModal,
  Stack,
} from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { format as formatSQL } from 'sql-formatter';
import { DataSource } from '../datasource';
import {
  RqliteDataSourceOptions,
  RqliteQuery,
  EditorMode,
  QueryFormat,
  ColumnSelection,
  WhereCondition,
  OrderByClause,
} from '../types';
import { TableSelect } from './visual-query-builder/TableSelect';
import { ColumnSelect } from './visual-query-builder/ColumnSelect';
import { WhereEditor } from './visual-query-builder/WhereEditor';
import { OrderBySelect } from './visual-query-builder/OrderBySelect';
import { GroupBySelect } from './visual-query-builder/GroupBySelect';
import { SQLPreview } from './visual-query-builder/SQLPreview';
import { generateSQL } from './visual-query-builder/sqlGenerator';

type Props = QueryEditorProps<DataSource, RqliteQuery, RqliteDataSourceOptions>;

const editorModeOptions: Array<SelectableValue<EditorMode>> = [
  { label: 'Code', value: 'code' },
  { label: 'Builder', value: 'builder' },
];

const formatOptions: Array<ComboboxOption<string>> = [
  { label: 'Table', value: 'table' },
  { label: 'Time series', value: 'time_series' },
];

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const {
    rawSql = '',
    format = 'table',
    timeColumns = ['time'],
    editorMode = 'code',
    table = '',
    columns = [],
    whereClause = [],
    groupBy = [],
    orderBy = [],
    limit = '',
    offset = '',
  } = query;

  const [macroRefOpen, setMacroRefOpen] = useState(false);
  const [whereOpen, setWhereOpen] = useState(whereClause.length > 0);
  const [groupByOpen, setGroupByOpen] = useState(groupBy.length > 0);
  const [orderByOpen, setOrderByOpen] = useState(orderBy.length > 0);
  const [previewOpen, setPreviewOpen] = useState(true);
  const [expandedEditor, setExpandedEditor] = useState(false);
  const [expandedSql, setExpandedSql] = useState(rawSql);
  const [confirmSwitchOpen, setConfirmSwitchOpen] = useState(false);

  // Generate SQL from builder state and sync to rawSql
  const generatedSQL = useMemo(
    () => generateSQL({ table, columns, whereClause, groupBy, orderBy, limit, offset }),
    [table, columns, whereClause, groupBy, orderBy, limit, offset]
  );

  useEffect(() => {
    if (editorMode === 'builder' && generatedSQL && generatedSQL !== rawSql) {
      onChange({ ...query, rawSql: generatedSQL });
    }
  }, [editorMode, generatedSQL]); // eslint-disable-line react-hooks/exhaustive-deps

  const onEditorModeChange = useCallback(
    (mode: EditorMode) => {
      if (mode === 'builder' && editorMode === 'code' && rawSql.trim()) {
        setConfirmSwitchOpen(true);
        return;
      }
      onChange({ ...query, editorMode: mode });
    },
    [onChange, query, editorMode, rawSql]
  );

  const onConfirmSwitch = useCallback(() => {
    setConfirmSwitchOpen(false);
    onChange({ ...query, editorMode: 'builder' });
  }, [onChange, query]);

  const onFormatChange = useCallback(
    (option: ComboboxOption<string>) => {
      onChange({ ...query, format: (option.value as QueryFormat) || 'table' });
      onRunQuery();
    },
    [onChange, onRunQuery, query]
  );

  const onTimeColumnsChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      const cols = event.target.value
        .split(',')
        .map((s) => s.trim())
        .filter(Boolean);
      onChange({ ...query, timeColumns: cols });
    },
    [onChange, query]
  );

  const onRawSqlChange = useCallback(
    (sql: string) => {
      onChange({ ...query, rawSql: sql });
    },
    [onChange, query]
  );

  const onTableChange = useCallback(
    (t: string) => {
      onChange({ ...query, table: t, columns: [], whereClause: [], groupBy: [], orderBy: [] });
    },
    [onChange, query]
  );

  const onColumnsChange = useCallback(
    (cols: ColumnSelection[]) => {
      onChange({ ...query, columns: cols });
    },
    [onChange, query]
  );

  const onWhereChange = useCallback(
    (conditions: WhereCondition[]) => {
      onChange({ ...query, whereClause: conditions });
    },
    [onChange, query]
  );

  const onGroupByChange = useCallback(
    (groups: string[]) => {
      onChange({ ...query, groupBy: groups });
    },
    [onChange, query]
  );

  const onOrderByChange = useCallback(
    (ob: OrderByClause[]) => {
      onChange({ ...query, orderBy: ob });
    },
    [onChange, query]
  );

  const onLimitChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      onChange({ ...query, limit: event.target.value });
    },
    [onChange, query]
  );

  const onOffsetChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      onChange({ ...query, offset: event.target.value });
    },
    [onChange, query]
  );

  const onFormatSQL = useCallback(() => {
    try {
      const formatted = formatSQL(rawSql, { language: 'sqlite' });
      onChange({ ...query, rawSql: formatted });
    } catch {
      // If formatting fails, leave SQL unchanged
    }
  }, [rawSql, onChange, query]);

  const onEditorDidMount = useCallback(
    (editor: MonacoEditor) => {
      editor.addAction({
        id: 'run-query',
        label: 'Run Query',
        keybindings: [
          // Monaco KeyMod.CtrlCmd | Monaco KeyCode.Enter
          // CtrlCmd = 2048, Enter = 3
          2048 | 3,
        ],
        run: () => {
          onRunQuery();
        },
      });
    },
    [onRunQuery]
  );

  const onExpandedEditorOpen = useCallback(() => {
    setExpandedSql(rawSql);
    setExpandedEditor(true);
  }, [rawSql]);

  const onExpandedEditorApply = useCallback(() => {
    onChange({ ...query, rawSql: expandedSql });
    setExpandedEditor(false);
  }, [onChange, query, expandedSql]);

  return (
    <div>
      <InlineFieldRow>
        <InlineField label="Mode" labelWidth={12}>
          <RadioButtonGroup options={editorModeOptions} value={editorMode} onChange={onEditorModeChange} />
        </InlineField>
        <InlineField label="Format" labelWidth={12}>
          <Combobox options={formatOptions} value={format} onChange={onFormatChange} width={20} />
        </InlineField>
        <InlineField label="Time columns" labelWidth={14} tooltip="Comma-separated list of columns to parse as time">
          <Input
            value={timeColumns.join(', ')}
            onChange={onTimeColumnsChange}
            placeholder="time"
            width={30}
          />
        </InlineField>
      </InlineFieldRow>

      {editorMode === 'code' && (
        <>
          <Stack direction="row" gap={1} alignItems="center" wrap="wrap">
            <Button variant="secondary" size="sm" icon="brackets-curly" onClick={onFormatSQL}>
              Format
            </Button>
            <Button variant="secondary" size="sm" icon="expand-arrows" onClick={onExpandedEditorOpen}>
              Expand
            </Button>
            <span style={{ fontSize: 12, color: '#8e8e8e' }}>Ctrl+Enter to run</span>
          </Stack>
          <CodeEditor
            value={rawSql}
            language="sql"
            height={200}
            onBlur={onRawSqlChange}
            onSave={onRawSqlChange}
            onEditorDidMount={onEditorDidMount}
            showMiniMap={false}
            showLineNumbers
          />
          <Collapse label="Macro Reference" isOpen={macroRefOpen} onToggle={() => setMacroRefOpen(!macroRefOpen)}>
            <pre style={{ fontSize: 12, padding: 8 }}>
              {`$__timeFilter(column)  → column >= <from> AND column <= <to>
$__timeFrom           → Unix epoch seconds (from)
$__timeTo             → Unix epoch seconds (to)
$__timeGroup(col, 5m) → (CAST(col / 300 AS INTEGER) * 300)
$__unixEpochFilter(c) → alias for $__timeFilter`}
            </pre>
          </Collapse>
          <Modal
            title="Edit SQL"
            isOpen={expandedEditor}
            onDismiss={() => setExpandedEditor(false)}
          >
            <CodeEditor
              value={expandedSql}
              language="sql"
              height="60vh"
              onBlur={setExpandedSql}
              onSave={setExpandedSql}
              showMiniMap={false}
              showLineNumbers
            />
            <Modal.ButtonRow>
              <Button variant="secondary" onClick={() => setExpandedEditor(false)}>
                Cancel
              </Button>
              <Button variant="primary" onClick={onExpandedEditorApply}>
                Apply
              </Button>
            </Modal.ButtonRow>
          </Modal>
        </>
      )}

      {editorMode === 'builder' && (
        <div>
          <TableSelect datasource={datasource} value={table} onChange={onTableChange} />
          <ColumnSelect datasource={datasource} table={table} value={columns} onChange={onColumnsChange} />
          <Collapse label="WHERE" isOpen={whereOpen} onToggle={() => setWhereOpen(!whereOpen)}>
            <WhereEditor datasource={datasource} table={table} value={whereClause} onChange={onWhereChange} />
          </Collapse>
          <Collapse label="GROUP BY" isOpen={groupByOpen} onToggle={() => setGroupByOpen(!groupByOpen)}>
            <GroupBySelect datasource={datasource} table={table} value={groupBy} onChange={onGroupByChange} />
          </Collapse>
          <Collapse label="ORDER BY" isOpen={orderByOpen} onToggle={() => setOrderByOpen(!orderByOpen)}>
            <OrderBySelect datasource={datasource} table={table} value={orderBy} onChange={onOrderByChange} />
          </Collapse>
          <InlineFieldRow>
            <InlineField label="LIMIT" labelWidth={12}>
              <Input value={limit} onChange={onLimitChange} placeholder="e.g. 100" width={15} />
            </InlineField>
            <InlineField label="OFFSET" labelWidth={12}>
              <Input value={offset} onChange={onOffsetChange} placeholder="e.g. 0" width={15} />
            </InlineField>
          </InlineFieldRow>
          <Collapse label="SQL Preview" isOpen={previewOpen} onToggle={() => setPreviewOpen(!previewOpen)}>
            <SQLPreview sql={generatedSQL} />
          </Collapse>
        </div>
      )}

      <ConfirmModal
        isOpen={confirmSwitchOpen}
        title="Switch to Builder mode?"
        body="Switching to Builder mode will discard any manual SQL edits. The query will be rebuilt from the builder fields."
        confirmText="Switch"
        onConfirm={onConfirmSwitch}
        onDismiss={() => setConfirmSwitchOpen(false)}
      />
    </div>
  );
}
