import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, QueryInput } from '../types';

type Props = QueryEditorProps<DataSource, QueryInput, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTopicChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, topic: event.target.value });
  };

  const { topic } = query;

  return (
    <Stack gap={0}>
      <InlineField label="Topic" labelWidth={16} tooltip="Enter topic to listen from">
        <Input
          id="query-editor-query-text"
          onChange={onTopicChange}
          value={topic || ''}
          required
          placeholder="Enter a topic"
        />
      </InlineField>
    </Stack>
  );
}
