import { ApiURL, ENCRYPT } from '../../config.js';
export const SessionKey = 'session';
const AuthProvider = {
    // called when the user attempts to log in
    login: ({ username, password }) => {
        try {
            password = ENCRYPT(password)
            const request = new Request(`${ApiURL}/authentication/auth`, {
                method: 'POST',
                body: JSON.stringify({ username, password }),
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'X-Auth-Header': 'X-Auth-Header',
                }),
            });
            return fetch(request)
                .then(response => {
                    if (response.status < 200 || response.status >= 300) {
                        throw new Error(response.statusText);
                    }
                    return response.json();
                })
                .then((response) => {
                    localStorage.setItem(SessionKey, response.session);
                });
        } catch (error) {
            return Promise.reject(error);
        }
    },

    // called when the user clicks on the logout button
    logout: () => {
        const session = localStorage.getItem(SessionKey)
        const request = new Request(`${ApiURL}/authentication/logout`, {
            method: 'DELETE',
            headers: new Headers({
                'X-Session-Key': session,
                'X-Auth-Header': 'X-Auth-Header',
            }),
        });
        return fetch(request)
            .then(response => {
                // if (response.status < 200 || response.status >= 300) {
                //     // return Promise.reject();
                //     throw new Error(response.statusText);
                // }
                return response.json();
            })
            .then((response) => {
                Promise.resolve();
            });
    },

    // called when the API returns an error
    checkError: ({ status }) => {
        if (status === 401 || status === 403) {
            localStorage.removeItem(SessionKey);
            return Promise.reject();
        }
        return Promise.resolve();
    },

    // called when the user navigates to a new location, to check for authentication
    checkAuth: () => {
        try {
            const session = localStorage.getItem(SessionKey)
            const request = new Request(`${ApiURL}/authentication/ping`, {
                method: 'GET',
                headers: new Headers({
                    'X-Session-Key': session,
                    'X-Auth-Header': 'X-Auth-Header',
                }),
            });
            return fetch(request)
                .then(response => {
                    if (response.status < 200 || response.status >= 300) {
                        // return Promise.reject();
                        throw new Error(response.statusText);
                    }
                    return response.json();
                })
                .then((response) => {
                    Promise.resolve();
                });
        } catch (error) {
            return Promise.reject(error);
        }
    },

    // called when the user navigates to a new location,
    // to check for permissions / roles
    getPermissions: () => {
        try {
            const session = localStorage.getItem(SessionKey)
            const request = new Request(`${ApiURL}/authentication/perms`, {
                method: 'GET',
                headers: new Headers({
                    'X-Session-Key': session,
                    'X-Auth-Header': 'X-Auth-Header',
                }),
            });
            return fetch(request)
                .then(response => {
                    if (response.status < 200 || response.status >= 300) {
                        // return Promise.reject();
                        throw new Error(response.statusText);
                    }
                    return response.json();
                })
                .then((response) => {
                    var role = 'User';
                    if (response.admin === true) {
                        role = 'Admin';
                    }
                    return Promise.resolve(role);
                });
        } catch (error) {
            return Promise.reject(error);
        }
    },

    // called for display userdata
    getIdentity: () => {
        try {
            const session = localStorage.getItem(SessionKey);
            const request = new Request(`${ApiURL}/authentication/info`, {
                method: 'GET',
                headers: new Headers({
                    'X-Session-Key': session,
                    'X-Auth-Header': 'X-Auth-Header',
                }),
            });
            return fetch(request)
                .then(response => {
                    if (response.status < 200 || response.status >= 300) {
                        throw new Error(response.statusText);
                    }
                    return response.json();
                })
                .then((response) => {
                    const { id, fullName, avatar } = response;
                    return Promise.resolve({ id, fullName, avatar });
                });
        } catch (error) {
            return Promise.reject(error);
        }
    }
};

export default AuthProvider;