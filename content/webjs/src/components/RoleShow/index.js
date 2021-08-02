import {
  TextField,
  ReferenceField,
  SimpleShowLayout,
  Show,
  Datagrid,
  ArrayField
} from 'react-admin';
import { DateTimestampField } from "../../fields/DateTimestampField"

export const RoleShow = (props) => {
  return (
    <div>
      <Show title=' ' {...props}>
        <SimpleShowLayout>
          <ArrayField source='data' label=''>
            <Datagrid>
              <ReferenceField
                link="show"
                label="Node"
                source='name'
                reference="nodes"
              >
                <TextField source="id" />
              </ReferenceField>
              <TextField source="ipaddress" />
              <DateTimestampField source="ohai_time" />
              <ArrayField source="run_list" >
                <Datagrid>
                  <TextField source="object" label=' ' />
                </Datagrid>
              </ArrayField>
            </Datagrid>
          </ArrayField>
        </SimpleShowLayout>
      </Show>
    </div >
  );
};

export default RoleShow;