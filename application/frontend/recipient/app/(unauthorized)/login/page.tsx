'use client';

import React, { useEffect, useState } from 'react';
import type { FormProps } from 'antd';
import { Button, DatePicker, Form, Input, notification, Select, Typography } from 'antd';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Country, getCountries } from '@/utils/getCountries';
import { BACKEND_URL } from '@/utils/env';
import { useAuth } from '@/components/AuthProvider';

const { Text } = Typography

type FieldType = {
  countryID?: number;
  NIN?: string;
  fullName?: string;
  dateOfBirth?: string;
  password?: string;
};

const Login = () => {
  const router = useRouter();

  const { setAuth } = useAuth();

  const onFinish: FormProps<FieldType>['onFinish'] = async (values) => {
    const res = await fetch(`${BACKEND_URL}/api/students/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(values),
    });
    const data = await res.json();
    if (res.status !== 200) {
      notification.error({ message: 'Error', description: data.error });
      return;
    }
    setAuth({
      isAuthenticated: true,
      privateKey: data.privateKey,
      publicKey: data.publicKey,
      fullName: data.fullName,
    });
    notification.success({ message: 'Success', description: 'You have successfully logged in!' });
    router.push("/home");
  };

  const onFinishFailed: FormProps<FieldType>['onFinishFailed'] = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  const [countries, setCountries] = useState<{ value: number, label: string }[]>([]);

  useEffect(() => {
    getCountries().then(data => {
      setCountries(data.map(c => {
        return {
          value: c.id,
          label: c.name,
        };
      }));
    });
  }, []);

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
        >
          <Form.Item
            name="countryID"
            label="Country"
            rules={[{ required: true, message: 'Please select your country!' }]}
          >
            <Select options={countries} />
          </Form.Item>

          <Form.Item
            name="NIN"
            label="National Identification Number (NIN)"
            rules={[
              {
                type: 'string',
                required: true,
                message: 'Please input your NIN!',
              },
            ]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="fullName"
            label="Full Name"
            rules={[
              {
                type: 'string',
                required: true,
                message: 'Please input your full name!',
              },
            ]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="dateOfBirth"
            label="Date of Birth"
            rules={[
              {
                type: 'date',
                required: true,
                message: 'Please input your date of birth!',
              },
            ]}
          >
            <DatePicker />
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
  )
};

export default Login;


