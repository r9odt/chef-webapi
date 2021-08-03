import { ENCRYPT } from '../../config.js';
import {
  TextInput,
  PasswordInput,
  SimpleForm,
  Toolbar,
  SaveButton,
  TextField,
  Edit,
  required,
} from "react-admin";
import { Grid, Card, CardContent, Typography } from '@material-ui/core';

const handleSubmit = (data) => {
  if (data.password) {
    data.password = ENCRYPT(data.password)
  }
  return data
};

export const ProfileEdit = (props) => {
  return (
    <div>
      <Edit
        id={"edit"}
        resource="profiles"
        basePath="/profile"
        redirect={false}
        title="My profile"
        transform={handleSubmit}
        undoable={false}
        {...props}
      >
        <SimpleForm
          toolbar={<ProfileEditToolbar />}
          sanitizeEmptyValues={false}>
          <Card fullWidth>
            <CardContent>
              <Grid container spacing={1}>
                <Grid item xs={12} sm={12} md={6}>
                  <Typography variant="h6" gutterBottom>
                    Identity
                  </Typography>
                </Grid>
                <Grid container spacing={1}>
                  <Grid item xs={12} sm={12} md={6}>
                    <Grid item xs={12} sm={12} md={6}>
                      <TextInput source="fullName" validate={required()}
                        fullWidth />
                    </Grid>
                    <Grid item xs={12} sm={12} md={6}>
                      <TextInput source="avatar" fullWidth />
                    </Grid>
                    <Grid item xs={12} sm={12} md={6}>
                      <PasswordInput source='password' fullWidth />
                    </Grid>
                  </Grid>
                  <Grid item xs={12} sm={12} md={6}>
                    <Grid item xs={12} sm={12} md={6}>
                      <Typography gutterBottom>
                        User ID
                      </Typography>
                      <TextField label="User ID" source="userID" fullWidth />
                    </Grid>
                    <Grid item xs={12} sm={12} md={6}>
                      <Typography gutterBottom>
                        User Name
                      </Typography>
                      <TextField label="User Name" source="username" fullWidth />
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </SimpleForm>
      </Edit>
    </div>
  );
};

export default ProfileEdit;

const ProfileEditToolbar = props => {
  return (
    <div>
      <Toolbar {...props}>
        <SaveButton submitOnEnter={true} />
      </Toolbar>
    </div>
  )
};
