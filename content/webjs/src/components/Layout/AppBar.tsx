import { AppBar, UserMenu, MenuItemLink } from 'react-admin';
import { forwardRef } from 'react';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import SettingsIcon from '@material-ui/icons/Settings';


const useStyles = makeStyles({
    title: {
        flex: 1,
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
        overflow: 'hidden',
    },
    spacer: {
        flex: 1,
    },
});

const Profile = forwardRef<any, any>((props, ref) => {
    return (
        <MenuItemLink
            ref={ref}
            to="/profile"
            primaryText={"Profile"}
            leftIcon={<SettingsIcon />}
            onClick={props.onClick}
            sidebarIsOpen
        />
    );
});

const CustomUserMenu = (props: any) => (
    <UserMenu {...props}>
        <Profile />
    </UserMenu>
);

export const CustomAppBar = (props: any) => {
    const classes = useStyles();
    return (
        <AppBar {...props} elevation={1} userMenu={<CustomUserMenu />}>
            <Typography
                variant="h6"
                color="inherit"
                className={classes.title}
                id="react-admin-title"
            />
            {/* <Logo /> */}
            <span className={classes.spacer} />
        </AppBar>
    );
};
