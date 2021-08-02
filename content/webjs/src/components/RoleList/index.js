import {
  List,
  Datagrid,
  EditButton,
  TextField
} from 'react-admin';
import classnames from 'classnames';
import { Fragment, useCallback } from 'react';
import { Route, useHistory } from 'react-router-dom';
import Searcher from '../../searchers/Searcher';
import DeployButton from '../../buttons/deploy/DeployButton';
import { RoleShow } from '../RoleShow';
import { rolesResource } from "../../App.js";
import { makeStyles } from '@material-ui/core/styles';
import { Drawer } from '@material-ui/core';
import RoleEdit from '../RoleEdit';

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
export const RoleList = (props) => {
  const classes = useStyles();
  const history = useHistory();

  const handleClose = useCallback(() => {
    history.push(`/${rolesResource}`);
  }, [history]);
  return (
    <div>
      <Route path={`/${rolesResource}/:id`}>
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
                  title='List of roles'
                  filters={<Searcher />}
                  className={classnames(classes.list, {
                    [classes.listWithDrawer]: isMatch,
                  })}
                  {...props}>
                  <Datagrid
                    rowClick="expand"
                    expand={<RoleShow />}>
                    <TextField source='id' />
                    <TextField label="Last Deploy time" source='date' />
                    <DeployButton deployResource={`${rolesResource}`} />
                    <DeployButton
                      displayLabel={"Role"}
                      deployResource={`${rolesResource}`}
                      onlyResource={true} />
                    <EditButton />
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
                    <RoleEdit
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

export default RoleList;