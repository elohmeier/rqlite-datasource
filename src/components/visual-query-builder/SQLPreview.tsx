import React from 'react';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { CodeEditor, ClipboardButton, Stack, useStyles2 } from '@grafana/ui';

interface Props {
  sql: string;
}

export function SQLPreview({ sql }: Props) {
  const styles = useStyles2(getStyles);

  if (!sql) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
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

const getStyles = (theme: GrafanaTheme2) => ({
  wrapper: css({
    marginTop: theme.spacing(1),
  }),
});
