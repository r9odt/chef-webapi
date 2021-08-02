import {
  TextField,
  Edit,
  SimpleForm,
  BooleanInput,
  SaveButton,
  Toolbar
} from 'react-admin';

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

export const AppModuleEdit = (props) => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <IconButton onClick={props.onCancel}>
        <CloseIcon />
      </IconButton>
      <Edit
        undoable={false}
        title=' '
        {...props}>
        <SimpleForm toolbar={<AppModuleEditToolbar />}
          className={classes.form}
          sanitizeEmptyValues={false}
        >
          <TextField source='name' />
          <BooleanInput label="Active" source='isON' />
        </SimpleForm>
      </Edit>
    </div>
  );
};

const AppModuleEditToolbar = props => {
  const classes = useStyles();
  return (
    <div>
      <Toolbar className={classes.toolBarRoot} {...props}>
        <SaveButton submitOnEnter={true} />
      </Toolbar>
    </div>
  )
};

export default AppModuleEdit;