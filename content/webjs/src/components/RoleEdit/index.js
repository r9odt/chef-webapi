import {
  TextField,
  Edit,
  SimpleForm,
  ArrayInput,
  FormDataConsumer,
  SimpleFormIterator,
  BooleanInput,
  SaveButton,
  Toolbar
} from 'react-admin';

import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import { makeStyles } from '@material-ui/core/styles';
import { DeployIcon } from '../../buttons/deploy/DeployButton';

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

export const RoleEdit = (props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <IconButton onClick={props.onCancel}>
        <CloseIcon />
      </IconButton>
      <Edit
        undoable={false}
        title=' ' {...props}>
        <SimpleForm toolbar={<RoleEditToolbar />}
          className={classes.form}
          sanitizeEmptyValues={false}
          {...props}
        >
          <TextField source="id" label="Role" />
          <BooleanInput source="onlyResource" label="Only Role" />
          <ArrayInput source='data' label=''>
            <SimpleFormIterator disableRemove disableAdd>
              {/* <TextInput source="name" label="Node" /> */}
              <FormDataConsumer>
                {({ getSource, scopedFormData }) => {
                  // console.log(scopedFormData)
                  return (
                    <TextField
                      label={getSource("name")}
                      source={"name"}
                      record={scopedFormData}
                    />
                  );
                }}
              </FormDataConsumer>
              <BooleanInput source="selected" label="Select" />
            </SimpleFormIterator>
          </ArrayInput>
        </SimpleForm>
      </Edit>
    </div>
  );
};

const RoleEditToolbar = props => {
  const classes = useStyles();
  return (
    <div>
      <Toolbar className={classes.toolBarRoot} {...props}>
        <SaveButton
          label={"Deploy selected"}
          icon={DeployIcon}
        />
      </Toolbar>
    </div>
  )
};

export default RoleEdit;