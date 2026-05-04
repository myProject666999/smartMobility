import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, DatePicker, Table, Tag } from 'antd';
import {
  UserOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import { getStatistics, getOrders, getUsers } from '../services/api';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

const Dashboard = () => {
  const [statistics, setStatistics] = useState({});
  const [orders, setOrders] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [dateRange, setDateRange] = useState(null);

  useEffect(() => {
    loadStatistics();
    loadRecentOrders();
    loadRecentUsers();
  }, []);

  const loadStatistics = async (startDate, endDate) => {
    setLoading(true);
    try {
      const params = {};
      if (startDate) params.start_date = startDate;
      if (endDate) params.end_date = endDate;
      const res = await getStatistics(params);
      setStatistics(res.data || {});
    } catch (error) {
      console.error('加载统计数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadRecentOrders = async () => {
    try {
      const res = await getOrders({ page: 1, page_size: 10 });
      setOrders(res.data?.list || []);
    } catch (error) {
      console.error('加载订单失败:', error);
    }
  };

  const loadRecentUsers = async () => {
    try {
      const res = await getUsers({ page: 1, page_size: 10 });
      setUsers(res.data?.list || []);
    } catch (error) {
      console.error('加载用户失败:', error);
    }
  };

  const handleDateChange = (dates) => {
    setDateRange(dates);
    if (dates && dates.length === 2) {
      loadStatistics(
        dates[0].format('YYYY-MM-DD'),
        dates[1].format('YYYY-MM-DD')
      );
    } else {
      loadStatistics();
    }
  };

  const getStatusTag = (status) => {
    const statusMap = {
      0: { text: '待支付', color: 'orange' },
      1: { text: '已支付', color: 'green' },
      2: { text: '已取消', color: 'default' },
      3: { text: '已完成', color: 'blue' },
    };
    const s = statusMap[status] || { text: '未知', color: 'default' };
    return <Tag color={s.color}>{s.text}</Tag>;
  };

  const orderColumns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    { title: '用户', dataIndex: ['user', 'username'], key: 'user' },
    { title: '金额', dataIndex: 'total_price', key: 'total_price', render: (val) => `¥${val}` },
    { title: '状态', dataIndex: 'status', key: 'status', render: getStatusTag },
    { title: '下单时间', dataIndex: 'created_at', key: 'created_at', render: (val) => dayjs(val).format('YYYY-MM-DD HH:mm') },
  ];

  const userColumns = [
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '真实姓名', dataIndex: 'real_name', key: 'real_name' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (val) => (
      <Tag color={val === 1 ? 'green' : 'red'}>{val === 1 ? '正常' : '禁用'}</Tag>
    )},
    { title: '注册时间', dataIndex: 'created_at', key: 'created_at', render: (val) => dayjs(val).format('YYYY-MM-DD HH:mm') },
  ];

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <RangePicker onChange={handleDateChange} value={dateRange} />
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="售票数量"
              value={statistics.ticket_count || 0}
              prefix={<TicketOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="售票收入"
              value={statistics.ticket_revenue || 0}
              precision={2}
              prefix={<DollarOutlined />}
              suffix="元"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="用户总数"
              value={statistics.user_count || 0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="待处理订单"
              value={statistics.pending_orders || 0}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Card title="最近订单">
            <Table
              columns={orderColumns}
              dataSource={orders}
              rowKey="id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="最近用户">
            <Table
              columns={userColumns}
              dataSource={users}
              rowKey="id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
