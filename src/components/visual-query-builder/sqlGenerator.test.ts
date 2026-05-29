import { BuilderState, generateSQL } from './sqlGenerator';

const baseState: BuilderState = {
  table: 'users',
  columns: [],
  whereClause: [],
  groupBy: [],
  orderBy: [],
  limit: '',
  offset: '',
};

const buildState = (overrides: Partial<BuilderState>): BuilderState => ({
  ...baseState,
  ...overrides,
});

describe('generateSQL', () => {
  it('escapes single quotes in WHERE values', () => {
    const sql = generateSQL(
      buildState({
        whereClause: [{ column: 'name', operator: '=', value: "O'Brien" }],
      })
    );

    expect(sql).toBe(`SELECT *
FROM "users"
WHERE "name" = 'O''Brien'`);
  });

  it('quotes safe identifiers and builds supported clauses', () => {
    const sql = generateSQL(
      buildState({
        table: 'sales',
        columns: [
          { name: 'region', aggregation: '' },
          { name: 'amount', aggregation: 'SUM' },
        ],
        whereClause: [{ column: 'customer', operator: 'LIKE', value: "O'Brien%" }],
        groupBy: ['region'],
        orderBy: [{ column: 'amount', direction: 'DESC' }],
        limit: ' 10 ',
        offset: '0',
      })
    );

    expect(sql).toBe(`SELECT "region", SUM("amount")
FROM "sales"
WHERE "customer" LIKE 'O''Brien%'
GROUP BY "region"
ORDER BY "amount" DESC
LIMIT 10
OFFSET 0`);
  });

  it('quotes and escapes IN values as comma-separated literals', () => {
    const sql = generateSQL(
      buildState({
        whereClause: [{ column: 'name', operator: 'IN', value: "O'Brien, Smith" }],
      })
    );

    expect(sql).toBe(`SELECT *
FROM "users"
WHERE "name" IN ('O''Brien', 'Smith')`);
  });

  it.each([
    ['table', { table: 'users;DROP' }],
    ['selected column', { columns: [{ name: 'name;DROP', aggregation: '' }] }],
    ['where column', { whereClause: [{ column: 'name;DROP', operator: '=', value: 'Ada' }] }],
    ['group by column', { groupBy: ['name;DROP'] }],
    ['order by column', { orderBy: [{ column: 'name;DROP', direction: 'ASC' }] }],
  ])('returns an empty query for an unsafe %s identifier', (_case, overrides) => {
    expect(generateSQL(buildState(overrides as Partial<BuilderState>))).toBe('');
  });

  it.each([
    ['aggregation', { columns: [{ name: 'amount', aggregation: 'SUM);DROP' }] }],
    ['where operator', { whereClause: [{ column: 'name', operator: '= 1 OR 1=1', value: 'Ada' }] }],
    ['order direction', { orderBy: [{ column: 'name', direction: 'ASC;DROP' }] }],
    ['limit', { limit: '1;DROP' }],
    ['offset', { offset: '1;DROP' }],
  ])('returns an empty query for an unsafe %s', (_case, overrides) => {
    expect(generateSQL(buildState(overrides as Partial<BuilderState>))).toBe('');
  });
});
