import {
  List,
  Datagrid,
  TextField
} from 'react-admin';
import Searcher from '../../searchers/Searcher';


export const AppKeyList = (props) => {

  return (
    <div>
      <List bulkActionButtons={false}
        title='List of application keys'
        filters={<Searcher />}
        {...props}>
        <Datagrid
          rowClick="edit"
        >
          <TextField label="Name" source='id' />
        </Datagrid>
      </List>
    </div>
  );
};

export default AppKeyList;
