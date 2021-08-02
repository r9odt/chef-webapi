import {
  TextField,
  SimpleShowLayout,
  Show,
  Datagrid,
  ArrayField
} from 'react-admin';
import { DateTimestampField } from "../../fields/DateTimestampField"

export const NodeShow = (props) => {
  return (
    <div>
      <Show title=' ' {...props}>
        <SimpleShowLayout>
          <ArrayField source='data' label=''>
            <Datagrid>
              <TextField source="fqdn" />
              <TextField source="ipaddress" />
              <ArrayField source="run_list" >
                <Datagrid>
                  <TextField source="object" label=' ' />
                </Datagrid>
              </ArrayField>
              <DateTimestampField source="ohai_time" />
            </Datagrid>
          </ArrayField>
        </SimpleShowLayout>
      </Show>
    </div>
  );
};

export default NodeShow;