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

  const onStartTimeChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, start_time: event.target.value });
  };

  return (
    <Stack gap={0}>
      <InlineField label="Topic" labelWidth={16} tooltip="Enter topic to listen from">
        <Input
          id="query-topic"
          onChange={onTopicChange}
          value={query.topic || ''}
          onBlur={onRunQuery}
          placeholder="Enter a topic"
        />
      </InlineField>
    </Stack>
  );
}
