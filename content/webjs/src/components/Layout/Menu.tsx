import { FC } from 'react';
import { Box } from '@material-ui/core';
import {
  MenuItemLink,
  MenuProps,
  DashboardMenuItem,
  usePermissions
} from 'react-admin';
import BookIcon from "@material-ui/icons/Book";
import DomainIcon from "@material-ui/icons/Domain";
import GroupIcon from "@material-ui/icons/Group";
import AccessTime from "@material-ui/icons/AccessTime";
import SupervisedUserCircleIcon from '@material-ui/icons/SupervisedUserCircle';
import NewReleasesIcon from '@material-ui/icons/NewReleases';

export const CustomMenu: FC<MenuProps> = ({ onMenuClick, dense = false }) => {
  const { permissions } = usePermissions()
  return (
    <Box mt={1}>
      <DashboardMenuItem onClick={onMenuClick} sidebarIsOpen={true} />
      <MenuItemLink
        to={`/deployers`}
        primaryText={"Tasks"}
        leftIcon={<AccessTime />}
        onClick={onMenuClick}
        sidebarIsOpen={true}
        dense={dense}
      />
      <MenuItemLink
        to={`/nodes`}
        primaryText={"Nodes"}
        leftIcon={<DomainIcon />}
        onClick={onMenuClick}
        sidebarIsOpen={true}
        dense={dense}
      />
      <MenuItemLink
        to={`/roles`}
        primaryText={"Roles"}
        leftIcon={<GroupIcon />}
        onClick={onMenuClick}
        sidebarIsOpen={true}
        dense={dense}
      />
      <MenuItemLink
        to={`/cookbooks`}
        primaryText={"Cookbooks"}
        leftIcon={<BookIcon />}
        onClick={onMenuClick}
        sidebarIsOpen={true}
        dense={dense}
      />
      {
        permissions === 'Admin' &&
        <MenuItemLink
          to={`/users`}
          primaryText={"Users"}
          leftIcon={<SupervisedUserCircleIcon />}
          onClick={onMenuClick}
          sidebarIsOpen={true}
          dense={dense}
        />
      }
      {        
        permissions === 'Admin' &&
        <MenuItemLink
          to={`/modules`}
          primaryText={"Modules"}
          leftIcon={<NewReleasesIcon />}
          onClick={onMenuClick}
          sidebarIsOpen={true}
          dense={dense}
        />
      }
      {        
        permissions === 'Admin' &&
        <MenuItemLink
          to={`/keys`}
          primaryText={"Keys"}
          leftIcon={<NewReleasesIcon />}
          onClick={onMenuClick}
          sidebarIsOpen={true}
          dense={dense}
        />
      }
    </Box>
  );
};
