import {
  TextField,
  Edit,
  SimpleForm,
  ArrayInput,
  BooleanInput,
  SimpleFormIterator,
  SaveButton,
  FormDataConsumer,
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

export const NodeEdit = (props) => {
  const classes = useStyles();

  // console.log(props)
  return (
    <div className={classes.root}>
      <IconButton onClick={props.onCancel}>
        <CloseIcon />
      </IconButton>
      <Edit
        undoable={false}
        title={' '} {...props}>
        <SimpleForm toolbar={<NodeEditToolbar />}
          className={classes.form}
          sanitizeEmptyValues={false}
        >
          <TextField source="id" label="Node" />
          <ArrayInput source="data" label=''>
            <SimpleFormIterator disableRemove disableAdd>
                  <FormDataConsumer>
                    {({ getSource, scopedFormData }) => {
                      return (
                        <TextField
                          label={getSource("fqdn")}
                          source={"fqdn"}
                          record={scopedFormData}
                        />
                      );
                    }}
                  </FormDataConsumer>
              <ArrayInput source="run_list" label=''>
                <SimpleFormIterator disableRemove disableAdd>
                  <TextField source="object" label="Node" />
                  <FormDataConsumer>
                    {({ getSource, scopedFormData }) => {
                      return (
                        <TextField
                          label={getSource("object")}
                          source={"object"}
                          record={scopedFormData}
                        />
                      );
                    }}
                  </FormDataConsumer>
                  <BooleanInput source="selected" label="Select" />
                </SimpleFormIterator>
              </ArrayInput>
            </SimpleFormIterator>
          </ArrayInput>
        </SimpleForm>
      </Edit>
    </div>
  );
};

const NodeEditToolbar = props => {
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

export default NodeEdit;