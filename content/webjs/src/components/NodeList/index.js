import {
  List,
  Datagrid,
  TextField,
  EditButton
} from 'react-admin';
import classnames from 'classnames';
import { Fragment, useCallback } from 'react';
import { Route, useHistory } from 'react-router-dom';
import Searcher from '../../searchers/Searcher';
import DeployButton from '../../buttons/deploy/DeployButton';
import { NodeShow } from '../NodeShow';
import { nodesResource } from "../../App.js";
import { makeStyles } from '@material-ui/core/styles';
import { Drawer } from '@material-ui/core';
import NodeEdit from '../NodeEdit';

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

export const NodeList = (props) => {
  const classes = useStyles();
  const history = useHistory();

  const handleClose = useCallback(() => {
    history.push(`/${nodesResource}`);
  }, [history]);
  return (
    <div>
      <Route path={`/${nodesResource}/:id`}>
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
                  title='List of nodes'
                  filters={<Searcher />}
                  className={classnames(classes.list, {
                    [classes.listWithDrawer]: isMatch,
                  })}
                  {...props}>
                  <Datagrid
                    rowClick="expand"
                    expand={<NodeShow />}>
                    <TextField source='id' />
                    <TextField label="Last Deploy time" source='date' />
                    <DeployButton deployResource={`${nodesResource}`} />
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
                    <NodeEdit
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

export default NodeList;