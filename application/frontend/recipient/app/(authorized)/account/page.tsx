"use client"
import React, { useState } from 'react';
import { Button, Form, Input, Modal, Typography, notification } from 'antd';
const { Title } = Typography;

export default function Account() {
  const [isChangePasswordModalVisible, setIsChangePasswordModalVisible] = useState(false);
  const [isDeleteAccountModalVisible, setIsDeleteAccountModalVisible] = useState(false);

  // Function to handle password change submission
  const handleChangePassword = (values: any) => {
    console.log('Password Change:', values); // You can add logic to send this data to the backend
    notification.success({
      message: 'Password Changed',
      description: 'Your password has been successfully changed.',
    });
    setIsChangePasswordModalVisible(false);
  };

  // Function to handle account deletion confirmation
  const handleDeleteAccount = () => {
    // You can add logic here to handle the account deletion.
    notification.success({
      message: 'Account Deleted',
      description: 'Your account has been successfully deleted.',
    });
    setIsDeleteAccountModalVisible(false);
  };

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>Account</Title>

      {/* Change Password Section */}
      <Button
        type="primary"
        onClick={() => setIsChangePasswordModalVisible(true)}
        style={{ marginBottom: '20px', marginRight: '10px' }}
      >
        Change Password
      </Button>

      {/* Delete Account Section */}
      <Button
        danger
        onClick={() => setIsDeleteAccountModalVisible(true)}
        style={{ marginBottom: '20px' }}
      >
        Delete Account
      </Button>

      {/* Modal for Change Password */}
      <Modal
        title="Change Password"
        visible={isChangePasswordModalVisible}
        onCancel={() => setIsChangePasswordModalVisible(false)}
        footer={null}
      >
        <Form
          onFinish={handleChangePassword}
          layout="vertical"
          initialValues={{
            currentPassword: '',
            newPassword: '',
            confirmPassword: '',
          }}
        >
          <Form.Item
            label="Current Password"
            name="currentPassword"
            rules={[{ required: true, message: 'Please enter your current password!' }]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item
            label="New Password"
            name="newPassword"
            rules={[{ required: true, message: 'Please enter your new password!' }]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item
            label="Confirm New Password"
            name="confirmPassword"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: 'Please confirm your new password!' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('The two passwords that you entered do not match!'));
                },
              }),
            ]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
              Change Password
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* Modal for Delete Account Confirmation */}
      <Modal
        title="Are you sure you want to delete your account?"
        visible={isDeleteAccountModalVisible}
        onCancel={() => setIsDeleteAccountModalVisible(false)}
        onOk={handleDeleteAccount}
        okText="Yes, Delete"
        cancelText="Cancel"
        okButtonProps={{ danger: true }}
      >
        <p>This action is irreversible. Please make sure you want to permanently delete your account.</p>
      </Modal>
    </div>
  );
}

