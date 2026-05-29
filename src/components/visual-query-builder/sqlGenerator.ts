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

const safeIdentifierPattern = /^[A-Za-z0-9_]+$/;
const aggregateFunctions = new Set(['COUNT', 'SUM', 'AVG', 'MIN', 'MAX']);
const operators = new Set(['=', '!=', '<', '>', '<=', '>=', 'LIKE', 'IN', 'IS NULL', 'IS NOT NULL']);
const noValueOperators = new Set(['IS NULL', 'IS NOT NULL']);
const orderDirections = new Set(['ASC', 'DESC']);
const positiveIntegerPattern = /^\d+$/;

export function generateSQL(state: BuilderState): string {
  if (!state.table) {
    return '';
  }
  const table = quoteSafeIdentifier(state.table);
  if (!table) {
    return '';
  }

  const parts: string[] = [];

  // SELECT
  const selectCols = buildSelectColumns(state.columns);
  if (!selectCols) {
    return '';
  }
  parts.push(`SELECT ${selectCols}`);

  // FROM
  parts.push(`FROM ${table}`);

  // WHERE
  const whereStr = buildWhere(state.whereClause);
  if (whereStr === null) {
    return '';
  }
  if (whereStr) {
    parts.push(`WHERE ${whereStr}`);
  }

  // GROUP BY
  if (state.groupBy.length > 0) {
    const groupBy = buildIdentifierList(state.groupBy);
    if (!groupBy) {
      return '';
    }
    parts.push(`GROUP BY ${groupBy}`);
  }

  // ORDER BY
  const orderByStr = buildOrderBy(state.orderBy);
  if (orderByStr === null) {
    return '';
  }
  if (orderByStr) {
    parts.push(`ORDER BY ${orderByStr}`);
  }

  // LIMIT
  if (state.limit) {
    const limit = buildPositiveInteger(state.limit);
    if (!limit) {
      return '';
    }
    parts.push(`LIMIT ${limit}`);
  }

  // OFFSET
  if (state.offset) {
    const offset = buildPositiveInteger(state.offset);
    if (!offset) {
      return '';
    }
    parts.push(`OFFSET ${offset}`);
  }

  return parts.join('\n');
}

function buildSelectColumns(columns: ColumnSelection[]): string | null {
  if (columns.length === 0) {
    return '*';
  }

  const parts: string[] = [];

  for (const col of columns) {
    const column = quoteSafeIdentifier(col.name);
    if (!column) {
      return null;
    }

    if (col.aggregation) {
      const aggregation = col.aggregation.trim().toUpperCase();
      if (!aggregateFunctions.has(aggregation)) {
        return null;
      }
      parts.push(`${aggregation}(${column})`);
      continue;
    }

    parts.push(column);
  }

  return parts.join(', ');
}

function buildWhere(conditions: WhereCondition[]): string | null {
  const parts: string[] = [];

  for (const condition of conditions) {
    if (!condition.column || !condition.operator) {
      continue;
    }

    const column = quoteSafeIdentifier(condition.column);
    const operator = normalizeOperator(condition.operator);
    if (!column || !operator) {
      return null;
    }

    if (noValueOperators.has(operator)) {
      parts.push(`${column} ${operator}`);
      continue;
    }

    if (operator === 'IN') {
      parts.push(`${column} IN (${buildInValues(condition.value ?? '')})`);
      continue;
    }

    parts.push(`${column} ${operator} ${quoteStringLiteral(condition.value ?? '')}`);
  }

  return parts.join(' AND ');
}

function buildOrderBy(orderBy: OrderByClause[]): string | null {
  const parts: string[] = [];

  for (const order of orderBy) {
    if (!order.column) {
      continue;
    }

    const column = quoteSafeIdentifier(order.column);
    const direction = (order.direction || 'ASC').trim().toUpperCase();
    if (!column || !orderDirections.has(direction)) {
      return null;
    }

    parts.push(`${column} ${direction}`);
  }

  return parts.join(', ');
}

function buildIdentifierList(identifiers: string[]): string | null {
  const parts: string[] = [];

  for (const identifier of identifiers) {
    const quoted = quoteSafeIdentifier(identifier);
    if (!quoted) {
      return null;
    }
    parts.push(quoted);
  }

  return parts.join(', ');
}

function normalizeOperator(operator: string): string | null {
  const normalized = operator.trim().toUpperCase();
  return operators.has(normalized) ? normalized : null;
}

function buildPositiveInteger(value: string): string | null {
  const normalized = value.trim();
  return positiveIntegerPattern.test(normalized) ? normalized : null;
}

function buildInValues(value: string): string {
  return value.split(',').map((item) => quoteStringLiteral(item.trim())).join(', ');
}

function quoteSafeIdentifier(identifier: string): string | null {
  if (!safeIdentifierPattern.test(identifier)) {
    return null;
  }
  return `"${identifier}"`;
}

function quoteStringLiteral(value: string): string {
  return `'${value.replace(/'/g, "''")}'`;
}
