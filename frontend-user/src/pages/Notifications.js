import React, { useState, useEffect } from 'react';
import { Card, List, Button, Badge, Empty } from 'antd';
import { getNotifications, markNotificationRead } from '../services/api';

const Notifications = () => {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadNotifications();
  }, []);

  const loadNotifications = async () => {
    setLoading(true);
    try {
      const res = await getNotifications({ page: 1, page_size: 100 });
      setNotifications(res.data?.list || []);
    } catch (error) {
      console.error('加载通知失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRead = async (id) => {
    try {
      await markNotificationRead(id);
      loadNotifications();
    } catch (error) {
      console.error('标记已读失败:', error);
    }
  };

  const getTypeText = (type) => {
    const typeMap = {
      0: '系统通知',
      1: '订单通知',
      2: '活动通知',
    };
    return typeMap[type] || '其他';
  };

  return (
    <div style={{ padding: 24 }}>
      <Card title="通知中心">
        {notifications.length > 0 ? (
          <List
            loading={loading}
            dataSource={notifications}
            renderItem={(item) => (
              <List.Item
                actions={[
                  item.is_read !== 1 && (
                    <Button type="link" onClick={() => handleRead(item.id)}>
                      标记已读
                    </Button>
                  ),
                ]}
              >
                <List.Item.Meta
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      {item.is_read !== 1 && <Badge status="processing" />}
                      <span style={{ fontWeight: item.is_read !== 1 ? 'bold' : 'normal' }}>
                        {item.title}
                      </span>
                      <span style={{ fontSize: 12, color: '#999' }}>
                        [{getTypeText(item.type)}]
                      </span>
                    </div>
                  }
                  description={
                    <div>
                      <p>{item.content}</p>
                      <span style={{ fontSize: 12, color: '#999' }}>
                        {new Date(item.created_at).toLocaleString()}
                      </span>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无通知" />
        )}
      </Card>
    </div>
  );
};

export default Notifications;
