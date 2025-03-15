'use client';
import React, { useEffect, useState } from 'react';
import {
  Button,
  DatePicker,
  Form,
  FormProps,
  Input,
  notification,
  Select,
  Typography,
} from 'antd';
import Link from 'next/link';
import { BACKEND_URL } from '@/utils/env';
import { useAuth } from '@/components/AuthProvider';
import { getCountries } from '@/utils/getCountries';
import { useRouter } from 'next/navigation';


const { Title, Text } = Typography;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 8 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 },
  },
};

const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0,
    },
    sm: {
      span: 16,
      offset: 8,
    },
  },
};

type FieldType = {
  countryID?: number;
  NIN?: string;
  fullName?: string;
  dateOfBirth?: string;
  password?: string;
};

const Register: React.FC = () => {
  const [form] = Form.useForm();
  const router = useRouter();

  const { setAuth } = useAuth();

  const onFinish: FormProps<FieldType>['onFinish'] = async (values) => {
    const res = await fetch(`${BACKEND_URL}/api/students`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(values),
    });
    const data = await res.json();
    if (res.status !== 200 && res.status !== 201) {
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
    <div style={{ padding: '40px 0', backgroundColor: '#f9f9f9' }}>
      <div style={{ maxWidth: '600px', margin: '0 auto', backgroundColor: 'white', padding: '30px', borderRadius: '8px', boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)' }}>
        <Title level={2} style={{ textAlign: 'center', marginBottom: '30px' }}>Register</Title>
        <Form
          {...formItemLayout}
          form={form}
          name="register"
          initialValues={{ prefix: '86' }}
          style={{ maxWidth: '100%' }}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          scrollToFirstError
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
            label="National Id Number (NIN)"
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
            name="password"
            label="Password"
            rules={[
              {
                required: true,
                message: 'Please input your password!',
              },
            ]}
            hasFeedback
          >
            <Input.Password />
          </Form.Item>

          <Form.Item
            name="confirm"
            label="Confirm Password"
            dependencies={['password']}
            hasFeedback
            rules={[
              {
                required: true,
                message: 'Please confirm your password!',
              },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('The passwords do not match!'));
                },
              }),
            ]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item {...tailFormItemLayout}>
            <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
              Register
            </Button>
          </Form.Item>

          <Text>
            You already have account?
            <Link href="/login" passHref>
              Login
            </Link>
          </Text>
        </Form>
      </div>
    </div>
  );
};

export default Register;


