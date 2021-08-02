import { defaultTheme } from 'react-admin';
import merge from 'lodash/merge';

export const Ltheme = merge({}, defaultTheme, {
  overrides: {
    MuiAppBar: {
      colorSecondary: {
        backgroundColor: '#616161e6',
      },
    },
  },
});