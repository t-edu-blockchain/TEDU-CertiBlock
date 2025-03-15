'use client';
import React, { useCallback, useEffect, useState } from 'react';
import { Button, Modal, QRCode, Table, Typography, Input, notification } from 'antd';
import { BACKEND_URL } from '@/utils/env';
import { useAuth } from '@/components/AuthProvider';
import { decode as base64Decode } from 'base64-arraybuffer';

const { Title, Text } = Typography;
const { Search } = Input;

// Define types for certificate data
interface Certificate {
  "certHash": string,
  "certUUID": string,
  "studentPublicKey": string,
  "universityPublicKey": string,
  "dateOfIssuing": string,
}

const Dashboard: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [selectedCertificate, setSelectedCertificate] = useState<Certificate | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>('');

  const { auth } = useAuth();

  // Demo data for multiple certificates
  const [certificates, setCertificates] = useState<Certificate[]>([]);

  const reloadCertificates = useCallback(() => {
    if (!auth.isAuthenticated) {
      notification.error({ message: 'Error', description: 'You are not authenticated' });
      return;
    }
    fetch(`${BACKEND_URL}/api/students/certificates`, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        privateKey: auth.privateKey,
      }),
    }).then(async res => {
      if (!res.ok) {
        throw new Error('Failed to fetch data');
      }
      const data = await res.json();
      setCertificates(data);
    }).catch(e => {
      notification.error({ message: 'Error', description: e.message });
    });
  }, [auth, setCertificates]);

  useEffect(() => {
    reloadCertificates();
  }, []);

  // Columns definition for the Ant Design Table
  const columns = [
    {
      title: 'UUID',
      dataIndex: 'certUUID',
      key: 'certUUID',
    },
    {
      title: 'Hash',
      dataIndex: 'certHash',
      key: 'certHash',
    },
    {
      title: 'Issue Date',
      dataIndex: 'dateOfIssuing',
      key: 'dateOfIssuing',
    },
    {
      title: 'University Public Key',
      dataIndex: 'universityPublicKey',
      key: 'universityPublicKey',
    },
    {
      title: 'Action',
      key: 'action',
      render: (text: any, record: Certificate) => (
        <Button type="primary" onClick={() => handleDownload(record)}>
          Download
        </Button>
      ),
    },
  ];

  const handleDownload = (certificate: Certificate): void => {
    if (!auth.isAuthenticated) {
      notification.error({ message: 'Error', description: 'You are not authenticated' });
      return;
    }

    fetch(`${BACKEND_URL}/api/students/certificates/${certificate.certUUID}`, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        privateKey: auth.privateKey,
      }),
    }).then(async res => {
      if (!res.ok) {
        throw new Error('Failed to fetch data');
      }
      const data = await res.json();
      const arrayBuffer = base64Decode(data.base64File)
      console.log("base64File", data.base64File);
      console.log("arrayBuffer", arrayBuffer);
      const blob = new Blob([arrayBuffer], { type: 'application/octet-stream' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${certificate.certUUID}.txt`;
      a.click();
    }).catch(e => {
      notification.error({ message: 'Error', description: e.message });
    })
  };

  const handleOk = (): void => {
    setIsModalOpen(false);
  };

  const handleCancel = (): void => {
    setIsModalOpen(false);
  };

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>{auth.isAuthenticated ? auth.fullName : "Unauthenticated"}</Title>
      <Title level={4}>{"Public Key: " + (auth.isAuthenticated ? auth.publicKey : "Unauthorized")}</Title>


      {/* Search Input */}
      <Search
        placeholder="Search by Certificate ID"
        onChange={(e) => setSearchQuery(e.target.value)}
        style={{ marginBottom: '20px', width: '100%' }}
      />

      <Button type="primary" style={{ marginBottom: '20px' }} onClick={() => reloadCertificates()}>
        Reload Certificates
      </Button>

      {/* Table to display certificates */}
      <Table
        columns={columns}
        dataSource={certificates}
        rowKey="certificateId"
        pagination={false} // Disable pagination, you can enable it if needed
      />

      {/* {selectedCertificate && (
        <Modal
          title="Share Certificate"
          open={isModalOpen}
          onOk={handleOk}
          onCancel={handleCancel}
          footer={[
            <Button key="back" onClick={handleCancel}>
              Close
            </Button>,
          ]}
        >
          <QRCode value={`https://example.com/certificate/${selectedCertificate.certUUID}`} />
          <div style={{ marginTop: '10px', display: 'flex', alignItems: 'center' }}>
            <Text copyable>https://verify.com/{selectedCertificate.certificateId}</Text>
          </div>
        </Modal>
      )} */}
    </div>
  );
};

export default Dashboard;

