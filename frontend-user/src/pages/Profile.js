import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, message, Tabs, Avatar } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { getUserInfo, updateUserInfo, updatePassword } from '../services/api';
import { setUser, getUser } from '../utils/auth';

const Profile = () => {
  const [loading, setLoading] = useState(false);
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [user, setUserState] = useState(null);
  const [infoForm] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    loadUserInfo();
  }, []);

  const loadUserInfo = async () => {
    try {
      const res = await getUserInfo();
      setUserState(res.data);
      infoForm.setFieldsValue(res.data);
    } catch (error) {
      console.error('加载用户信息失败:', error);
    }
  };

  const handleUpdateInfo = async (values) => {
    setLoading(true);
    try {
      await updateUserInfo(values);
      message.success('更新成功');
      loadUserInfo();
    } catch (error) {
      console.error('更新失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdatePassword = async (values) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的密码不一致');
      return;
    }
    
    setPasswordLoading(true);
    try {
      await updatePassword({
        old_password: values.oldPassword,
        new_password: values.newPassword,
      });
      message.success('密码修改成功');
      passwordForm.resetFields();
    } catch (error) {
      console.error('修改密码失败:', error);
    } finally {
      setPasswordLoading(false);
    }
  };

  const tabItems = [
    {
      key: 'info',
      label: '个人信息',
      children: (
        <Form
          form={infoForm}
          layout="vertical"
          onFinish={handleUpdateInfo}
          style={{ maxWidth: 500 }}
        >
          <Form.Item label="头像">
            <Avatar size={80} icon={<UserOutlined />} src={user?.avatar} />
          </Form.Item>
          <Form.Item label="用户名">
            <Input disabled />
          </Form.Item>
          <Form.Item name="real_name" label="真实姓名">
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input placeholder="请输入邮箱" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'password',
      label: '修改密码',
      children: (
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleUpdatePassword}
          style={{ maxWidth: 500 }}
        >
          <Form.Item
            name="oldPassword"
            label="原密码"
            rules={[{ required: true, message: '请输入原密码' }]}
          >
            <Input.Password placeholder="请输入原密码" />
          </Form.Item>
          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码（至少6位）" />
          </Form.Item>
          <Form.Item
            name="confirmPassword"
            label="确认密码"
            rules={[{ required: true, message: '请确认密码' }]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={passwordLoading}>
              修改密码
            </Button>
          </Form.Item>
        </Form>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Card title="个人中心">
        <Tabs items={tabItems} />
      </Card>
    </div>
  );
};

export default Profile;
