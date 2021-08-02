import {
  TextField,
  Edit,
  SimpleForm,
  TextInput,
  SaveButton,
  Toolbar
} from 'react-admin';

export const AppKeyEdit = (props) => {
  return (
    <div>
      <Edit
        undoable={false}
        title=' '
        {...props}>
        <SimpleForm toolbar={<AppKeyEditToolbar />}
          sanitizeEmptyValues={false}
        >
          <TextField source='id' />
          <TextInput label="Value"
            source='value'
            fullWidth
            rowsMax={15}
            multiline
          />
        </SimpleForm>
      </Edit>
    </div>
  );
};

const AppKeyEditToolbar = props => {
  return (
    <div>
      <Toolbar {...props}>
        <SaveButton submitOnEnter={true} />
      </Toolbar>
    </div>
  )
};

export default AppKeyEdit;