import request from '../utils/request';

export const adminLogin = (data) => {
  return request.post('/auth/admin/login', data);
};

export const getStatistics = (params) => {
  return request.get('/admin/statistics', { params });
};

export const getUsers = (params) => {
  return request.get('/admin/users', { params });
};

export const updateUserStatus = (id, data) => {
  return request.put(`/admin/users/${id}/status`, data);
};

export const getStations = (params) => {
  return request.get('/admin/stations', { params });
};

export const getStationById = (id) => {
  return request.get(`/admin/stations/${id}`);
};

export const createStation = (data) => {
  return request.post('/admin/stations', data);
};

export const updateStation = (id, data) => {
  return request.put(`/admin/stations/${id}`, data);
};

export const deleteStation = (id) => {
  return request.delete(`/admin/stations/${id}`);
};

export const getRoutes = (params) => {
  return request.get('/admin/routes', { params });
};

export const getRouteById = (id) => {
  return request.get(`/admin/routes/${id}`);
};

export const createRoute = (data) => {
  return request.post('/admin/routes', data);
};

export const updateRoute = (id, data) => {
  return request.put(`/admin/routes/${id}`, data);
};

export const deleteRoute = (id) => {
  return request.delete(`/admin/routes/${id}`);
};

export const getTickets = (params) => {
  return request.get('/admin/tickets', { params });
};

export const getTicketById = (id) => {
  return request.get(`/admin/tickets/${id}`);
};

export const createTicket = (data) => {
  return request.post('/admin/tickets', data);
};

export const updateTicket = (id, data) => {
  return request.put(`/admin/tickets/${id}`, data);
};

export const deleteTicket = (id) => {
  return request.delete(`/admin/tickets/${id}`);
};

export const getOrders = (params) => {
  return request.get('/admin/orders', { params });
};

export const getOrderById = (id) => {
  return request.get(`/admin/orders/${id}`);
};

export const cancelOrder = (id) => {
  return request.post(`/admin/orders/${id}/cancel`);
};

export const refundOrder = (id) => {
  return request.post(`/admin/orders/${id}/refund`);
};

export const getReviews = (params) => {
  return request.get('/admin/reviews', { params });
};

export const updateReviewStatus = (id, data) => {
  return request.put(`/admin/reviews/${id}/status`, data);
};

export const getAnnouncements = (params) => {
  return request.get('/admin/announcements', { params });
};

export const getAnnouncementById = (id) => {
  return request.get(`/admin/announcements/${id}`);
};

export const createAnnouncement = (data) => {
  return request.post('/admin/announcements', data);
};

export const updateAnnouncement = (id, data) => {
  return request.put(`/admin/announcements/${id}`, data);
};

export const deleteAnnouncement = (id) => {
  return request.delete(`/admin/announcements/${id}`);
};

export const getBanners = () => {
  return request.get('/admin/banners');
};

export const createBanner = (data) => {
  return request.post('/admin/banners', data);
};

export const updateBanner = (id, data) => {
  return request.put(`/admin/banners/${id}`, data);
};

export const deleteBanner = (id) => {
  return request.delete(`/admin/banners/${id}`);
};

export const getMenus = () => {
  return request.get('/admin/menus');
};

export const createMenu = (data) => {
  return request.post('/admin/menus', data);
};

export const updateMenu = (id, data) => {
  return request.put(`/admin/menus/${id}`, data);
};

export const deleteMenu = (id) => {
  return request.delete(`/admin/menus/${id}`);
};

export const updatePassword = (data) => {
  return request.put('/admin/password', data);
};
