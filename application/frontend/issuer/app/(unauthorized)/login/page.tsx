'use client';
import React from 'react';
import type { FormProps } from 'antd';
import { Button, Form, Input, Typography } from 'antd';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

const {Text} = Typography

type FieldType = {
  email?: string;
  password?: string;
};

const Login = () => {
    const router = useRouter()

    const onFinish: FormProps<FieldType>['onFinish'] = (values) => {
      console.log('Success:', values);
      router.push("/home")
    };

    const onFinishFailed: FormProps<FieldType>['onFinishFailed'] = (errorInfo) => {
      console.log('Failed:', errorInfo);
    };
    return (
  <div style={{ padding: '40px 0', backgroundColor: '#f0f2f5' }}>
    <div
      style={{
        maxWidth: '400px',
        margin: '0 auto',
        backgroundColor: 'white',
        padding: '30px',
        borderRadius: '8px',
        boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Typography.Title level={2} style={{ textAlign: 'center', marginBottom: '30px' }}>
        Login
      </Typography.Title>

      <Form
        name="basic"
        labelCol={{ span: 24 }}
        wrapperCol={{ span: 24 }}
        style={{ maxWidth: '100%' }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
      >
        <Form.Item
          name="email"
          label="E-mail"
          rules={[
            {
              type: 'email',
              message: 'The input is not a valid E-mail!',
            },
            {
              required: true,
              message: 'Please input your E-mail!',
            },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[{ required: true, message: 'Please input your password!' }]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item style={{ textAlign: 'center' }}>
          <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
            Submit
          </Button>
        </Form.Item>

        <Text>
        You do not have account? 
        <Link href="/register" passHref>
        Register
        </Link>
        </Text>


      </Form>
    </div>
  </div>
) };

export default Login;


