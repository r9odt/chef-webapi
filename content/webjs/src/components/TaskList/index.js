import {
  List,
  Datagrid,
  TextField,
  BooleanField,
  Pagination,
  ReferenceField,
  ChipField,
  useRefresh
} from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';
import classnames from 'classnames';
import { RefreshTime, useRecursiveTimeout } from "../../App.js";
import { TaskShow } from "../TaskShow";

const useStyles = makeStyles({
  waiting: { },
  inprogress: { backgroundColor: '#efff00' },
  complete: { backgroundColor: '#74fd74' },
  error: { backgroundColor: '#ff7373' },
});

const TaskPagination = props => <Pagination
  rowsPerPageOptions={[5, 10, 15, 25]}
  {...props} />;

const ColoredChipField = props => {
  const classes = useStyles();

  const isWaiting = status => status === "Waiting";
  const isInProgress = status => status === "InProgress";
  const isComplete = status => status === "Complete";
  const isError = status => status === "Error";

  return (
    <ChipField
      className={classnames({
        [classes.waiting]: isWaiting(props.record[props.source]),
        [classes.inprogress]: isInProgress(props.record[props.source]),
        [classes.complete]: isComplete(props.record[props.source]),
        [classes.error]: isError(props.record[props.source]),
      })}
      {...props}
    />
  );
};

export const TaskList = props => {
  const refresh = useRefresh()
  useRecursiveTimeout(() => refresh(), RefreshTime)
  return (
    <div>
      <List
        perPage={5}
        bulkActionButtons={false}
        pagination={<TaskPagination />}
        title='Tasks list'
        {...props}>
        <Datagrid
          rowClick="expand"
          expand={<TaskShow />}>
          <TextField source='resource' />
          <BooleanField source='onlyResource' />
          <BooleanField source='selectedResource' />
          <TextField source='resources' />
          <TextField source='name' />
          <ColoredChipField source='status' />
          <ReferenceField
            label="Initiator"
            source="initiatorID"
            reference="profiles"
            link="show"
          >
            <TextField source="fullName" />
          </ReferenceField>
          <TextField source='date' />
        </Datagrid>
      </List>
    </div>
  )
};

export default TaskList;