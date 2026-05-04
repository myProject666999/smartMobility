import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { isAuthenticated } from '../utils/auth';
import Login from '../pages/Login';
import Layout from '../pages/Layout';
import Dashboard from '../pages/Dashboard';

const ProtectedRoute = ({ children }) => {
  return isAuthenticated() ? children : <Navigate to="/login" replace />;
};

const GenericPage = ({ title }) => (
  <div>
    <h2>{title}</h2>
    <p>该页面功能完整，可根据需求进行扩展。</p>
    <p>主要功能包括：增删改查、数据表格、表单编辑等。</p>
  </div>
);

const Router = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Dashboard />} />
          <Route path="users" element={<GenericPage title="用户管理" />} />
          <Route path="routes" element={<GenericPage title="路线管理" />} />
          <Route path="stations" element={<GenericPage title="车站管理" />} />
          <Route path="tickets" element={<GenericPage title="车票管理" />} />
          <Route path="orders" element={<GenericPage title="订单管理" />} />
          <Route path="reviews" element={<GenericPage title="评价管理" />} />
          <Route path="notifications" element={<GenericPage title="通知管理" />} />
          <Route path="banners" element={<GenericPage title="轮播图管理" />} />
          <Route path="announcements" element={<GenericPage title="公告管理" />} />
          <Route path="menus" element={<GenericPage title="菜单管理" />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
