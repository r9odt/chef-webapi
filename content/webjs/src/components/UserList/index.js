import { Fragment, useCallback } from 'react';
import {
  List,
  Datagrid,
  TextField,
  BooleanField
} from 'react-admin';
import classnames from 'classnames';
import { Route, useHistory } from 'react-router-dom';
import { Drawer } from '@material-ui/core';
import Searcher from '../../searchers/Searcher';
import UserEdit from '../UserEdit';
import { usersResource } from '../../App.js';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  list: {
    flexGrow: 1,
    transition: theme.transitions.create(['all'], {
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginRight: 0,
  },
  listWithDrawer: {
    marginRight: 400,
  },
  drawerPaper: {
    zIndex: 100,
  },
}));

export const UserList = (props) => {
  const classes = useStyles();
  const history = useHistory();

  const handleClose = useCallback(() => {
    history.push(`/${usersResource}`);
  }, [history]);
  return (
    <div>
      <Route path={`/${usersResource}/:id`}>
        {({ match }) => {
          const isMatch = !!(
            match &&
            match.params &&
            match.params.id !== ''
          );
          return (
            <div>
              <Fragment>
                <List bulkActionButtons={false}
                  title='List of users'
                  filters={<Searcher />}
                  className={classnames(classes.list, {
                    [classes.listWithDrawer]: isMatch,
                  })}
                  {...props}>
                  <Datagrid
                    rowClick="edit"
                  >
                    <TextField source='username' />
                    <TextField source='fullName' />
                    <BooleanField source='admin' />
                    <BooleanField source='blocked' />
                    <BooleanField source='needPasswordChange' />
                    <TextField source='lastLogin' />
                    <TextField source='lastSeen' />
                  </Datagrid>
                </List>
                <Drawer
                  variant="persistent"
                  open={isMatch}
                  anchor="right"
                  onClose={handleClose}
                  classes={{
                    paper: classes.drawerPaper,
                  }}
                >
                  {/* To avoid any errors if the route does not match, we don't render at all the component in this case */}
                  {isMatch ? (
                    <UserEdit
                      id={match.params.id}
                      onCancel={handleClose}
                      {...props}
                    />
                  ) : null}
                </Drawer>
              </Fragment>
            </div>
          );
        }}
      </Route>
    </div>
  )
};

export default UserList;