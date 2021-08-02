import { SearchInput, Filter  } from 'react-admin';

const Searcher = (props) => (
  <Filter {...props}>
    <SearchInput source='q' alwaysOn />
  </Filter>
);

export default Searcher;