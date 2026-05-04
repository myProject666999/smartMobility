import React, { useState } from 'react';
import { Card, List, Button, Tag, Modal, Form, InputNumber, message, Space } from 'antd';
import { ShoppingCartOutlined } from '@ant-design/icons';
import { useLocation, useNavigate } from 'react-router-dom';
import { createOrder } from '../services/api';
import { isAuthenticated } from '../utils/auth';

const Search = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [results, setResults] = useState(location.state?.results || { tickets: [], routes: [] });
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedTicket, setSelectedTicket] = useState(null);
  const [form] = Form.useForm();

  const handleBuy = (ticket) => {
    if (!isAuthenticated()) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    setSelectedTicket(ticket);
    form.setFieldsValue({ quantity: 1 });
    setModalVisible(true);
  };

  const handleSubmit = async (values) => {
    if (!selectedTicket) return;
    
    setLoading(true);
    try {
      const res = await createOrder({
        ticket_id: selectedTicket.id,
        quantity: values.quantity,
      });
      message.success('下单成功');
      setModalVisible(false);
      navigate('/orders');
    } catch (error) {
      console.error('下单失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <Card title="搜索结果" style={{ marginBottom: 24 }}>
        {results.tickets && results.tickets.length > 0 ? (
          <List
            dataSource={results.tickets}
            renderItem={(ticket) => (
              <List.Item
                actions={[
                  <Button 
                    type="primary" 
                    icon={<ShoppingCartOutlined />}
                    onClick={() => handleBuy(ticket)}
                  >
                    购买
                  </Button>
                ]}
              >
                <List.Item.Meta
                  title={
                    <Space>
                      <Tag color="blue">{ticket.route?.route_number}</Tag>
                      <span style={{ fontWeight: 'bold' }}>
                        {ticket.route?.start_station?.name} → {ticket.route?.end_station?.name}
                      </span>
                    </Space>
                  }
                  description={
                    <Space direction="vertical" size={0}>
                      <span>发车日期: {ticket.depart_date}</span>
                      <span>发车时间: {ticket.depart_time}</span>
                      <span>余票: {ticket.seats_total - ticket.seats_sold} / {ticket.seats_total}</span>
                    </Space>
                  }
                />
                <div style={{ textAlign: 'right' }}>
                  <div style={{ fontSize: 24, color: '#1890ff', fontWeight: 'bold' }}>
                    ¥{ticket.price}
                  </div>
                </div>
              </List.Item>
            )}
          />
        ) : (
          <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
            暂无符合条件的车票
          </div>
        )}
      </Card>

      <Modal
        title="确认购买"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
      >
        {selectedTicket && (
          <Form form={form} onFinish={handleSubmit} layout="vertical">
            <div style={{ marginBottom: 16, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
              <p><strong>路线:</strong> {selectedTicket.route?.route_number}</p>
              <p><strong>起点:</strong> {selectedTicket.route?.start_station?.name}</p>
              <p><strong>终点:</strong> {selectedTicket.route?.end_station?.name}</p>
              <p><strong>日期:</strong> {selectedTicket.depart_date} {selectedTicket.depart_time}</p>
              <p><strong>票价:</strong> ¥{selectedTicket.price}</p>
            </div>
            <Form.Item
              name="quantity"
              label="购买数量"
              rules={[{ required: true, message: '请输入数量' }]}
            >
              <InputNumber 
                min={1} 
                max={selectedTicket.seats_total - selectedTicket.seats_sold}
                style={{ width: 200 }}
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading} block>
                确认购买
              </Button>
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  );
};

export default Search;
