import React, { useState } from 'react';
import { Layout, Menu, Dropdown, Avatar, Button, Space } from 'antd';
import {
  BarChartOutlined,
  UserOutlined,
  CarOutlined,
  EnvironmentOutlined,
  TicketOutlined,
  ShoppingCartOutlined,
  StarOutlined,
  BellOutlined,
  PictureOutlined,
  ReadOutlined,
  MenuOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { logout, getAdmin } from '../utils/auth';

const { Header, Sider, Content, Footer } = Layout;

const AdminLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const admin = getAdmin();

  const menuItems = [
    {
      key: '/',
      icon: <BarChartOutlined />,
      label: '统计分析',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/routes',
      icon: <CarOutlined />,
      label: '路线管理',
    },
    {
      key: '/stations',
      icon: <EnvironmentOutlined />,
      label: '车站管理',
    },
    {
      key: '/tickets',
      icon: <TicketOutlined />,
      label: '车票管理',
    },
    {
      key: '/orders',
      icon: <ShoppingCartOutlined />,
      label: '订单管理',
    },
    {
      key: '/reviews',
      icon: <StarOutlined />,
      label: '评价管理',
    },
    {
      key: '/notifications',
      icon: <BellOutlined />,
      label: '通知管理',
    },
    {
      key: '/banners',
      icon: <PictureOutlined />,
      label: '轮播图管理',
    },
    {
      key: '/announcements',
      icon: <ReadOutlined />,
      label: '公告管理',
    },
    {
      key: '/menus',
      icon: <MenuOutlined />,
      label: '菜单管理',
    },
  ];

  const handleMenuClick = (e) => {
    navigate(e.key);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenu = {
    items: [
      {
        key: '1',
        label: '退出登录',
        icon: <LogoutOutlined />,
        onClick: handleLogout,
      },
    ],
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={(value) => setCollapsed(value)}
      >
        <div style={{
          height: 64,
          margin: 16,
          background: 'rgba(255, 255, 255, 0.1)',
          borderRadius: 6,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'white',
          fontSize: collapsed ? 12 : 16,
          fontWeight: 'bold',
        }}>
          {collapsed ? '🚌' : '🚌 智慧出行'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: '0 24px', background: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}>
          <Space>
            <Dropdown menu={userMenu} placement="bottomRight">
              <Button type="link" style={{ display: 'flex', alignItems: 'center' }}>
                <Avatar icon={<UserOutlined />} />
                <span style={{ marginLeft: 8 }}>{admin?.real_name || admin?.username || '管理员'}</span>
              </Button>
            </Dropdown>
          </Space>
        </Header>
        <Content style={{ margin: '24px 16px', padding: 24, background: '#fff', minHeight: 280 }}>
          <Outlet />
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          智慧出行系统管理后台 ©2024 Created with Golang + React
        </Footer>
      </Layout>
    </Layout>
  );
};

export default AdminLayout;
