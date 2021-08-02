import {
  Admin,
  Resource
} from 'react-admin';
import {
  useEffect,
  useRef
} from 'react';
import { Route } from 'react-router-dom';
import './App.css';
import { NodeList } from './components/NodeList';
import { RoleList } from './components/RoleList';
import { Dashboard } from './components/Dashboard';
import TaskList from './components/TaskList';
import { UserCreate } from './components/UserCreate';
import { UserList } from './components/UserList';
import CookbookList from './components/Cookbooks';
import AuthProvider from './providers/AuthProvider';
import ApiProvider from './providers/ApiProvider';
import Layout from './components/Layout/Layout';
import ProfileEdit from './components/MyProfileEdit';
import ProfileShow from './components/ProfileShow';
import { AppModuleList } from './components/AppModuleList';
import { AppKeyList } from './components/AppKeyList';
import { AppKeyEdit } from './components/AppKeyEdit';
import { ProfileURL } from './config.js';

export const deployersResource = 'deployers';
export const nodesResource = 'nodes';
export const rolesResource = 'roles';
export const cookbooksResource = 'cookbooks';
export const usersResource = 'users';
export const profileResource = 'profile';
export const profilesResource = 'profiles';
export const appModulesResource = 'modules';
export const appKeysResource = 'keys';
export const RefreshTime = 10000;

export const useRecursiveTimeout = (callback, delay) => {
  const savedCallback = useRef(callback)

  useEffect(() => {
    savedCallback.current = callback
  }, [callback])

  useEffect(() => {
    let id
    function tick() {
      const ret = savedCallback.current()

      if (ret instanceof Promise) {
        ret.then(() => {
          if (delay !== null) {
            id = setTimeout(tick, delay)
          }
        })
      } else {
        if (delay !== null) {
          id = setTimeout(tick, delay)
        }
      }
    }
    if (delay !== null) {
      id = setTimeout(tick, delay)
      return () => id && clearTimeout(id)
    }
  }, [delay])
}


function App() {
  return (
    <Admin
      layout={Layout}
      dashboard={Dashboard}
      authProvider={AuthProvider}
      dataProvider={ApiProvider}
      customRoutes={[
        <Route
          key="my-profile"
          path={ProfileURL}
          render={() => <ProfileEdit />}
        />,
      ]}
    >
      {permissions => [
        <Resource
          name={profilesResource}
          show={ProfileShow}
        />,
        <Resource
          name={deployersResource}
          list={TaskList}
        />,
        <Resource
          name={nodesResource}
          list={NodeList}
        // show={NodeShow}
        />,
        <Resource
          name={rolesResource}
          list={RoleList}
        // show={RoleShow}
        />,
        <Resource
          name={cookbooksResource}
          list={CookbookList}
        />,
        // Only include the users resource for admin users
        permissions === 'Admin'
          ?
          <Resource
            name={usersResource}
            list={UserList}
            create={UserCreate}
          // edit={UserEdit}
          />
          : null,
        permissions === 'Admin'
          ?
          <Resource
            name={appModulesResource}
            list={AppModuleList}
          />
          : null,
        permissions === 'Admin'
          ?
          <Resource
            name={appKeysResource}
            list={AppKeyList}
            edit={AppKeyEdit}
          />
          : null,
        <Resource name={profileResource} />,
      ]}
      {/* <Resource name='test' list={ListGuesser} /> */}
    </Admin>
  );
}

export default App;
