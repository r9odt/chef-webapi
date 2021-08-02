import { fetchUtils } from 'react-admin';
import jsonServerProvider from 'ra-data-json-server';
import { SessionKey } from '../AuthProvider';
import { ApiURL } from '../../config.js';

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' });
  }

  const session = localStorage.getItem(SessionKey)
  options.headers.set('X-Session-Key', session);
  return fetchUtils.fetchJson(url, options);
};

const ApiProvider = jsonServerProvider(ApiURL, httpClient);
export default ApiProvider;