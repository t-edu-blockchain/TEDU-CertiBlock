"use client"
import React, { useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined,  HomeOutlined, DashboardOutlined, FileSearchOutlined, UserOutlined } from '@ant-design/icons';
import { Button, Layout, Menu, theme } from 'antd';
import {useRouter} from 'next/navigation'

const { Header, Sider, Content } = Layout;

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();
    const router = useRouter();

  return (
    <Layout>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div className="demo-logo-vertical" />
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
        items={[
          {
            key: '1',
            icon: <HomeOutlined />,
            label: 'Home',
            onClick: () => router.push('/home')
          },
          {
            key: '2',
            icon: <DashboardOutlined />,
            label: 'Dashboard',
            onClick: () => router.push('/dashboard')
          },
          {
            key: '3',
            icon: <FileSearchOutlined />,
            label: 'Issue Certificate',
            onClick: () => router.push('/issue')
          },
          {
            key: '4',
            icon: <UserOutlined />,
            label: 'Account',
            onClick: () => router.push('/account')
          },
        ]}

        />
      </Sider>
      <Layout>
        <Header style={{ padding: 0, background: colorBgContainer }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
          }}
        >
        {children}
        </Content>
      </Layout>
    </Layout>
  );
}
