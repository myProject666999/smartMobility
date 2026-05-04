import React, { useState, useEffect } from 'react';
import { Card, List, Button, Tag, message, Modal, Form, InputNumber, Rate, Input, Space, Tabs } from 'antd';
import { PayCircleOutlined, CloseCircleOutlined, UndoOutlined, StarOutlined } from '@ant-design/icons';
import { getOrders, payOrder, cancelOrder, refundOrder, createReview } from '../services/api';

const { TextArea } = Input;

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [reviewModalVisible, setReviewModalVisible] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadOrders();
  }, []);

  const loadOrders = async (status) => {
    setLoading(true);
    try {
      const params = { page: 1, page_size: 100 };
      if (status !== undefined && status !== 'all') {
        params.status = status;
      }
      const res = await getOrders(params);
      setOrders(res.data?.list || []);
    } catch (error) {
      console.error('加载订单失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusText = (status) => {
    const statusMap = {
      0: { text: '待支付', color: 'orange' },
      1: { text: '已支付', color: 'green' },
      2: { text: '已取消', color: 'default' },
      3: { text: '已完成', color: 'blue' },
    };
    return statusMap[status] || { text: '未知', color: 'default' };
  };

  const handlePay = async (orderId) => {
    try {
      await payOrder(orderId);
      message.success('支付成功');
      loadOrders();
    } catch (error) {
      console.error('支付失败:', error);
    }
  };

  const handleCancel = async (orderId) => {
    try {
      await cancelOrder(orderId);
      message.success('订单已取消');
      loadOrders();
    } catch (error) {
      console.error('取消订单失败:', error);
    }
  };

  const handleRefund = async (orderId) => {
    try {
      await refundOrder(orderId);
      message.success('退款成功');
      loadOrders();
    } catch (error) {
      console.error('退款失败:', error);
    }
  };

  const handleReview = (order) => {
    setSelectedOrder(order);
    form.setFieldsValue({ rating: 5 });
    setReviewModalVisible(true);
  };

  const handleSubmitReview = async (values) => {
    if (!selectedOrder) return;
    
    try {
      await createReview({
        order_id: selectedOrder.id,
        rating: values.rating,
        content: values.content,
      });
      message.success('评价成功');
      setReviewModalVisible(false);
    } catch (error) {
      console.error('评价失败:', error);
    }
  };

  const tabItems = [
    { key: 'all', label: '全部订单' },
    { key: '0', label: '待支付' },
    { key: '1', label: '已支付' },
    { key: '2', label: '已取消' },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card>
        <Tabs 
          items={tabItems} 
          onChange={(key) => loadOrders(key === 'all' ? undefined : key)}
        />
        
        <List
          loading={loading}
          dataSource={orders}
          renderItem={(order) => {
            const status = getStatusText(order.status);
            return (
              <List.Item
                actions={
                  order.status === 0 ? [
                    <Button type="primary" icon={<PayCircleOutlined />} onClick={() => handlePay(order.id)}>
                      支付
                    </Button>,
                    <Button icon={<CloseCircleOutlined />} onClick={() => handleCancel(order.id)}>
                      取消
                    </Button>
                  ] : order.status === 1 ? [
                    <Button icon={<UndoOutlined />} onClick={() => handleRefund(order.id)}>
                      退款
                    </Button>,
                    <Button icon={<StarOutlined />} onClick={() => handleReview(order)}>
                      评价
                    </Button>
                  ] : []
                }
              >
                <List.Item.Meta
                  title={
                    <Space>
                      <Tag color={status.color}>{status.text}</Tag>
                      <span>订单号: {order.order_no}</span>
                    </Space>
                  }
                  description={
                    <Space direction="vertical" size={0}>
                      <span>路线: {order.ticket?.route?.route_number} - {order.ticket?.route?.start_station?.name} → {order.ticket?.route?.end_station?.name}</span>
                      <span>日期: {order.ticket?.depart_date} {order.ticket?.depart_time}</span>
                      <span>数量: {order.quantity} 张</span>
                      <span>下单时间: {new Date(order.created_at).toLocaleString()}</span>
                    </Space>
                  }
                />
                <div style={{ textAlign: 'right' }}>
                  <div style={{ fontSize: 18, color: '#f5222d', fontWeight: 'bold' }}>
                    ¥{order.total_price}
                  </div>
                </div>
              </List.Item>
            );
          }}
        />
      </Card>

      <Modal
        title="发表评价"
        open={reviewModalVisible}
        onCancel={() => setReviewModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmitReview} layout="vertical">
          <Form.Item name="rating" label="评分" rules={[{ required: true, message: '请选择评分' }]}>
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <TextArea rows={4} placeholder="请输入您的评价内容" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交评价
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Orders;
