import {
  TextField,
  Show,
} from "react-admin";

import { Box, Card, CardContent, Typography } from '@material-ui/core';

export const ProfileShow = (props) => {
  return (
    <div>
      <Show
        title="Profile view"
        {...props}
      >
        <Card>
          <CardContent>
            <Box display={{ md: 'block', lg: 'flex' }}>
              <Box flex={2} mr={{ md: 0, lg: '1em' }}>
                <Typography variant="h6" gutterBottom>
                  Identity
                </Typography>
                <Box display={{ xs: 'block', sm: 'flex' }}>
                  <Box
                    flex={1}
                    mr={{ xs: 0, sm: '0.5em' }}
                  >
                    <Typography gutterBottom>
                      Profile ID
                    </Typography>
                    <TextField label="id" source='id' />
                  </Box>
                  <Box
                    flex={1}
                    ml={{ xs: 0, sm: '0.5em' }}
                  >
                    <Typography gutterBottom>
                      Full Name
                    </Typography>
                    <TextField source="fullName" />
                  </Box>
                </Box>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Show>
    </div>
  );
};

export default ProfileShow;
