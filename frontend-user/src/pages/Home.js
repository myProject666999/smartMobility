import React, { useState, useEffect } from 'react';
import { Layout, Menu, Carousel, Form, Input, Button, Card, List, Typography, message, Tag, Space } from 'antd';
import { HomeOutlined, CarOutlined, ShoppingCartOutlined, UserOutlined, SearchOutlined, BellOutlined } from '@ant-design/icons';
import { useNavigate, Outlet } from 'react-router-dom';
import { getBanners, searchRoutes, getAnnouncements } from '../services/api';
import { logout, isAuthenticated } from '../utils/auth';

const { Header, Content, Footer, Sider } = Layout;
const { Title, Text } = Typography;

const Home = () => {
  const [banners, setBanners] = useState([]);
  const [announcements, setAnnouncements] = useState([]);
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  useEffect(() => {
    loadBanners();
    loadAnnouncements();
  }, []);

  const loadBanners = async () => {
    try {
      const res = await getBanners();
      setBanners(res.data || []);
    } catch (error) {
      console.error('加载轮播图失败:', error);
    }
  };

  const loadAnnouncements = async () => {
    try {
      const res = await getAnnouncements({ page: 1, page_size: 5 });
      setAnnouncements(res.data?.list || []);
    } catch (error) {
      console.error('加载公告失败:', error);
    }
  };

  const onSearch = async (values) => {
    if (!values.start_station || !values.end_station) {
      message.warning('请输入起点和终点');
      return;
    }
    
    setLoading(true);
    try {
      const res = await searchRoutes({
        start_station: values.start_station,
        end_station: values.end_station,
        date: values.date,
      });
      navigate('/search', { state: { 
        results: res.data,
        searchParams: values 
      }});
    } catch (error) {
      console.error('搜索失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/search', icon: <SearchOutlined />, label: '路线查询' },
    { key: '/orders', icon: <ShoppingCartOutlined />, label: '我的订单' },
    { key: '/notifications', icon: <BellOutlined />, label: '通知中心' },
    { key: '/profile', icon: <UserOutlined />, label: '个人中心' },
  ];

  const handleMenuClick = (e) => {
    if (!isAuthenticated() && e.key !== '/') {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    navigate(e.key);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', background: '#001529' }}>
        <div style={{ color: 'white', fontSize: 20, fontWeight: 'bold', marginRight: 30 }}>
          🚌 智慧出行系统
        </div>
        <div style={{ flex: 1 }}>
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['/']}
            items={menuItems}
            onClick={handleMenuClick}
            style={{ border: 'none' }}
          />
        </div>
        <div>
          {isAuthenticated() ? (
            <Space>
              <Button type="link" style={{ color: 'white' }} onClick={() => navigate('/profile')}>
                <UserOutlined /> 个人中心
              </Button>
              <Button type="link" style={{ color: 'white' }} onClick={handleLogout}>
                退出登录
              </Button>
            </Space>
          ) : (
            <Space>
              <Button type="link" style={{ color: 'white' }} onClick={() => navigate('/login')}>
                登录
              </Button>
              <Button type="primary" onClick={() => navigate('/register')}>
                注册
              </Button>
            </Space>
          )}
        </div>
      </Header>
      
      <Layout>
        <Content style={{ background: '#f5f5f5' }}>
          <div style={{ padding: '24px 50px' }}>
            {banners.length > 0 && (
              <Carousel autoplay style={{ marginBottom: 24, borderRadius: 8, overflow: 'hidden' }}>
                {banners.map((banner) => (
                  <div key={banner.id}>
                    <div style={{
                      height: 300,
                      background: `linear-gradient(135deg, #667eea 0%, #764ba2 100%)`,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: 'white',
                      fontSize: 24,
                    }}>
                      <div style={{ textAlign: 'center' }}>
                        <Title level={2} style={{ color: 'white' }}>{banner.title}</Title>
                        <Text style={{ color: 'white' }}>点击查看详情</Text>
                      </div>
                    </div>
                  </div>
                ))}
              </Carousel>
            )}

            <Card title="车票查询" style={{ marginBottom: 24 }}>
              <Form
                form={form}
                layout="inline"
                onFinish={onSearch}
                style={{ justifyContent: 'center' }}
              >
                <Form.Item name="start_station" rules={[{ required: true, message: '请输入起点' }]}>
                  <Input placeholder="起点站" size="large" prefix={<CarOutlined />} style={{ width: 200 }} />
                </Form.Item>
                <Form.Item name="end_station" rules={[{ required: true, message: '请输入终点' }]}>
                  <Input placeholder="终点站" size="large" prefix={<CarOutlined />} style={{ width: 200 }} />
                </Form.Item>
                <Form.Item name="date">
                  <Input type="date" size="large" style={{ width: 200 }} />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" size="large" htmlType="submit" loading={loading} icon={<SearchOutlined />}>
                    查询
                  </Button>
                </Form.Item>
              </Form>
            </Card>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24 }}>
              <Card title="最新公告" extra={<a href="/announcements">查看更多</a>}>
                <List
                  dataSource={announcements}
                  renderItem={(item) => (
                    <List.Item>
                      <List.Item.Meta
                        title={
                          <Space>
                            {item.type === 1 && <Tag color="red">重要</Tag>}
                            <a href={`/announcements/${item.id}`}>{item.title}</a>
                          </Space>
                        }
                        description={new Date(item.created_at).toLocaleString()}
                      />
                    </List.Item>
                  )}
                />
              </Card>

              <Card title="热门路线">
                <List
                  dataSource={[
                    { name: '101路', from: '火车站', to: '市中心', price: 2 },
                    { name: '202路', from: '汽车站', to: '大学城', price: 3 },
                    { name: '303路', from: '机场', to: '高新区', price: 5 },
                  ]}
                  renderItem={(item) => (
                    <List.Item>
                      <List.Item.Meta
                        title={<Tag color="blue">{item.name}</Tag>}
                        description={`${item.from} → ${item.to}`}
                      />
                      <div>¥{item.price}</div>
                    </List.Item>
                  )}
                />
              </Card>
            </div>
          </div>
          <Outlet />
        </Content>
      </Layout>

      <Footer style={{ textAlign: 'center', background: '#001529', color: 'rgba(255,255,255,0.65)' }}>
        智慧出行系统 ©2024 Created with Golang + React
      </Footer>
    </Layout>
  );
};

export default Home;
