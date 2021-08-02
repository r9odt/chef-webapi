import {
  DateField
} from "react-admin";

export const DateTimestampField = props => {
  const recordWithTimestampAsInteger = {
    [props.source]: parseInt(props.record[props.source] * 1000, 10)
  };
  return <DateField
    {...props}
    record={recordWithTimestampAsInteger}
    showTime={true}
  />
}

export default DateTimestampField;