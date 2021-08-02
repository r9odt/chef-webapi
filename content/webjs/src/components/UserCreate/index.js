import {
  Create,
  TextInput,
  SimpleForm,
  PasswordInput,
  BooleanInput,
  required,
  useNotify,
  useRefresh,
  useRedirect
} from 'react-admin';
import { ENCRYPT } from '../../config.js';
import { usersResource } from '../../App.js';
// import { makeStyles } from '@material-ui/core/styles';

// const useStyles = makeStyles(theme => ({

// }));

const handleSubmit = (data) => {
  if (data.password) {
    data.password = ENCRYPT(data.password)
  }
  return data
};

export const UserCreate = (props) => {
  // const classes = useStyles();

  const notify = useNotify();
  const refresh = useRefresh();
  const redirect = useRedirect();
  const onSuccess = () => {
    notify('User created')
    redirect(`/${usersResource}`)
    refresh();
  };

  return (
    <div>
      <Create onSuccess={onSuccess} transform={handleSubmit}
        title='Create user' warnWhenUnsavedChanges {...props}>
        <SimpleForm>
          <TextInput label='User Name' source='username' validate={required()} />
          <TextInput source='fullName' validate={required()} />
          <BooleanInput source='admin' />
          <BooleanInput source='blocked' />
          <PasswordInput source='password' validate={required()} />
        </SimpleForm>
      </Create >
    </div>
  );
};

export default UserCreate;