import request from '../utils/request';

export const userLogin = (data) => {
  return request.post('/auth/user/login', data);
};

export const userRegister = (data) => {
  return request.post('/auth/user/register', data);
};

export const getUserInfo = () => {
  return request.get('/user/info');
};

export const updateUserInfo = (data) => {
  return request.put('/user/info', data);
};

export const updatePassword = (data) => {
  return request.put('/user/password', data);
};

export const getStations = (params) => {
  return request.get('/stations', { params });
};

export const getRoutes = (params) => {
  return request.get('/routes', { params });
};

export const getRouteById = (id) => {
  return request.get(`/routes/${id}`);
};

export const searchRoutes = (params) => {
  return request.get('/routes/search', { params });
};

export const getTickets = (params) => {
  return request.get('/tickets', { params });
};

export const getTicketById = (id) => {
  return request.get(`/tickets/${id}`);
};

export const getAnnouncements = (params) => {
  return request.get('/announcements', { params });
};

export const getAnnouncementById = (id) => {
  return request.get(`/announcements/${id}`);
};

export const getBanners = () => {
  return request.get('/banners');
};

export const getOrders = (params) => {
  return request.get('/user/orders', { params });
};

export const getOrderById = (id) => {
  return request.get(`/user/orders/${id}`);
};

export const createOrder = (data) => {
  return request.post('/user/orders', data);
};

export const payOrder = (id) => {
  return request.post(`/user/orders/${id}/pay`);
};

export const cancelOrder = (id) => {
  return request.post(`/user/orders/${id}/cancel`);
};

export const refundOrder = (id) => {
  return request.post(`/user/orders/${id}/refund`);
};

export const getReviews = (params) => {
  return request.get('/user/reviews', { params });
};

export const createReview = (data) => {
  return request.post('/user/reviews', data);
};

export const getNotifications = (params) => {
  return request.get('/user/notifications', { params });
};

export const markNotificationRead = (id) => {
  return request.put(`/user/notifications/${id}/read`);
};
