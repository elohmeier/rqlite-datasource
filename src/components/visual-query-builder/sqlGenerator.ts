import { ColumnSelection, WhereCondition, OrderByClause } from '../../types';

export interface BuilderState {
  table: string;
  columns: ColumnSelection[];
  whereClause: WhereCondition[];
  groupBy: string[];
  orderBy: OrderByClause[];
  limit: string;
  offset: string;
}

export function generateSQL(state: BuilderState): string {
  if (!state.table) {
    return '';
  }

  const parts: string[] = [];

  // SELECT
  const selectCols = buildSelectColumns(state.columns);
  parts.push(`SELECT ${selectCols}`);

  // FROM
  parts.push(`FROM ${state.table}`);

  // WHERE
  const whereStr = buildWhere(state.whereClause);
  if (whereStr) {
    parts.push(`WHERE ${whereStr}`);
  }

  // GROUP BY
  if (state.groupBy.length > 0) {
    parts.push(`GROUP BY ${state.groupBy.join(', ')}`);
  }

  // ORDER BY
  const orderByStr = buildOrderBy(state.orderBy);
  if (orderByStr) {
    parts.push(`ORDER BY ${orderByStr}`);
  }

  // LIMIT
  if (state.limit) {
    parts.push(`LIMIT ${state.limit}`);
  }

  // OFFSET
  if (state.offset) {
    parts.push(`OFFSET ${state.offset}`);
  }

  return parts.join('\n');
}

function buildSelectColumns(columns: ColumnSelection[]): string {
  if (columns.length === 0) {
    return '*';
  }

  return columns
    .map((col) => {
      if (col.aggregation) {
        return `${col.aggregation}(${col.name})`;
      }
      return col.name;
    })
    .join(', ');
}

function buildWhere(conditions: WhereCondition[]): string {
  const parts = conditions
    .filter((c) => c.column && c.operator)
    .map((c) => {
      if (c.operator === 'IS NULL' || c.operator === 'IS NOT NULL') {
        return `${c.column} ${c.operator}`;
      }
      if (c.operator === 'IN') {
        return `${c.column} IN (${c.value})`;
      }
      return `${c.column} ${c.operator} '${c.value}'`;
    });

  return parts.join(' AND ');
}

function buildOrderBy(orderBy: OrderByClause[]): string {
  return orderBy
    .filter((o) => o.column)
    .map((o) => `${o.column} ${o.direction || 'ASC'}`)
    .join(', ');
}
