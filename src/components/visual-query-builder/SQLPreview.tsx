import React from 'react';
import { CodeEditor, ClipboardButton, Stack } from '@grafana/ui';

interface Props {
  sql: string;
}

export function SQLPreview({ sql }: Props) {
  if (!sql) {
    return null;
  }

  return (
    <div style={{ marginTop: 8 }}>
      <Stack direction="row" alignItems="center" gap={1}>
        <label className="gf-form-label">Generated SQL</label>
        <ClipboardButton getText={() => sql} size="sm" variant="secondary">
          Copy
        </ClipboardButton>
      </Stack>
      <CodeEditor value={sql} language="sql" height={100} readOnly showMiniMap={false} showLineNumbers={false} />
    </div>
  );
}
