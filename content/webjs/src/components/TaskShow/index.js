import {
  SimpleShowLayout,
  Show,
  RichTextField,
} from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  pre: {
    whiteSpace: "pre-wrap",
  },
}));

export const TaskShow = (props) => {
  const classes = useStyles();
  return (
      <Show title=' ' {...props}>
        <SimpleShowLayout>
          <pre className={classes.pre}>
            <RichTextField source="log" />
          </pre>
        </SimpleShowLayout>
      </Show>
  );
};

export default TaskShow;