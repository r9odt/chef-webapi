import {
    Layout,
    LayoutProps,
    Sidebar,
    defaultTheme
} from 'react-admin';
import { Ltheme } from "./Theme"
import { CustomAppBar } from "./AppBar"
import { CustomMenu } from "./Menu"

const CustomSidebar = (props: any) => <Sidebar {...props} size={200} />;

const CustomLayout = (props: LayoutProps) => {
    return (
        <Layout
            {...props}
            appBar={CustomAppBar}
            sidebar={CustomSidebar}
            menu={CustomMenu}
            theme={defaultTheme && Ltheme}
        />
    );
};
export default CustomLayout;
