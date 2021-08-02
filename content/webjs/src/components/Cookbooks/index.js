import React from 'react';
import {
  List,
  Datagrid,
  TextField,
  UrlField,
  ArrayField
} from 'react-admin';
import Searcher from '../../searchers/Searcher';

export const CookbookList = props => (
  <div>
    <List bulkActionButtons={false}
      title='List of cookbooks'
      filters={<Searcher />} {...props}>
      <Datagrid>
        <TextField source='id' />
        <UrlField source='meta.url' />
        <ArrayField source='meta.versions' >
          <Datagrid>
            <UrlField source='url' />
            <TextField source='version' />
          </Datagrid>
        </ArrayField>
      </Datagrid>
    </List>
  </div>
);

export default CookbookList;