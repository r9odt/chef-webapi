import {
  TextField,
  Edit,
  TextInput,
  SimpleForm,
  PasswordInput,
  BooleanInput,
  required,
  SaveButton,
  DeleteButton,
  Toolbar
} from 'react-admin';

import { ENCRYPT } from '../../config.js';
import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  root: {
    paddingTop: 40,
  },
  title: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    margin: '1em',
  },
  form: {
    width: 400,
  },
  inlineField: {
    display: 'inline-block',
    width: '50%',
  },
  toolBarRoot: {
    display: 'flex',
    justifyContent: 'space-between',
  },
}));

const handleSubmit = (data) => {
  if (data.password) {
    data.password = ENCRYPT(data.password)
  }
  return data
};

export const UserEdit = (props) => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <IconButton onClick={props.onCancel}>
        <CloseIcon />
      </IconButton>
      <Edit
        undoable={false}
        transform={handleSubmit}
        title=' '
        {...props}>
        <SimpleForm toolbar={<UserEditToolbar />}
          className={classes.form}
          sanitizeEmptyValues={false}
        >
          <TextField source='username' />
          <TextInput source='fullName' validate={required()} />
          <BooleanInput source='admin' />
          <BooleanInput source='blocked' />
          <BooleanInput source='needPasswordChange' defaultValue={false} />
          <TextInput source='avatar' />
          <PasswordInput source='password' />
        </SimpleForm>
      </Edit>
    </div>
  );
};

const UserEditToolbar = props => {
  const classes = useStyles();
  return (
    <div>
      <Toolbar className={classes.toolBarRoot} {...props}>
        <SaveButton submitOnEnter={true} />
        <DeleteButton submitOnEnter={true} />
      </Toolbar>
    </div>
  )
};

export default UserEdit;