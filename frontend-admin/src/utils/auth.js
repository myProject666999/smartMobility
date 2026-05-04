const TOKEN_KEY = 'smart_mobility_admin_token';
const ADMIN_KEY = 'smart_mobility_admin';

export const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

export const setToken = (token) => {
  localStorage.setItem(TOKEN_KEY, token);
};

export const removeToken = () => {
  localStorage.removeItem(TOKEN_KEY);
};

export const getAdmin = () => {
  const adminStr = localStorage.getItem(ADMIN_KEY);
  try {
    return adminStr ? JSON.parse(adminStr) : null;
  } catch {
    return null;
  }
};

export const setAdmin = (admin) => {
  localStorage.setItem(ADMIN_KEY, JSON.stringify(admin));
};

export const removeAdmin = () => {
  localStorage.removeItem(ADMIN_KEY);
};

export const isAuthenticated = () => {
  return !!getToken();
};

export const logout = () => {
  removeToken();
  removeAdmin();
};
