'use client';
import React, { useState } from 'react';
import { Button, Modal, QRCode, Table, Typography, Input} from 'antd';

const { Title, Text } = Typography;
const { Search } = Input;

// Define types for certificate data
interface Certificate {
  certificateId: string;
  certificateName: string;
  issueDate: string;
  validity: string;
}

const Dashboard: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [selectedCertificate, setSelectedCertificate] = useState<Certificate | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>('');

  // Demo data for multiple certificates
  const certificates: Certificate[] = [
    {
      certificateId: '123456',
      certificateName: 'A',
      issueDate: 'March 2025',
      validity: 'Lifetime',
    },
    {
      certificateId: '789012',
      certificateName: 'B',
      issueDate: 'January 2025',
      validity: 'Lifetime',
    },
    {
      certificateId: '345678',
      certificateName: 'C',
      issueDate: 'July 2024',
      validity: '1 Year',
    },
    {
      certificateId: '345678',
      certificateName: 'C',
      issueDate: 'July 2024',
      validity: '1 Year',
    },
    {
      certificateId: '345678',
      certificateName: 'C',
      issueDate: 'July 2024',
      validity: '2 Year',
    },
    {
      certificateId: '901234',
      certificateName: 'D',
      issueDate: 'December 2024',
      validity: 'Lifetime',
    },
  ];

  // Columns definition for the Ant Design Table
  const columns = [
      {
          title: 'Certificate ID',
          dataIndex: 'certificateId',
          key: 'certificateId',
      },
      {
          title: 'Certificate Name',
          dataIndex: 'certificateName',
          key: 'certificateName',
      },
      {
          title: 'Issue Date',
          dataIndex: 'issueDate',
          key: 'issueDate',
      },
      {
          title: 'Validity',
          dataIndex: 'validity',
          key: 'validity',
      },
      {
          title: 'Action',
          key: 'action',
          render: (text: any, record: Certificate) => (
              <Button type="primary" onClick={() => handleShare(record)}>
              Share
              </Button>
          ),
      },
  ];

  // Filter certificates based on the search query
  const filteredCertificates = certificates.filter(
      (certificate) =>
      certificate.certificateName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      certificate.certificateId.includes(searchQuery)
  );

  const handleShare = (certificate: Certificate): void => {
      setSelectedCertificate(certificate);
      setIsModalOpen(true);
  };

  const handleOk = (): void => {
      setIsModalOpen(false);
  };

  const handleCancel = (): void => {
      setIsModalOpen(false);
  };

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>Dashboard</Title>


      {/* Search Input */}
      <Search
        placeholder="Search by Certificate ID"
        onChange={(e) => setSearchQuery(e.target.value)}
        style={{ marginBottom: '20px', width: '100%' }}
      />

      {/* Table to display certificates */}
      <Table
        columns={columns}
        dataSource={filteredCertificates}
        rowKey="certificateId"
        pagination={false} // Disable pagination, you can enable it if needed
      />

      {selectedCertificate && (
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
          <QRCode value={`https://example.com/certificate/${selectedCertificate.certificateId}`} />
          <div style={{ marginTop: '10px', display: 'flex', alignItems: 'center' }}>
            <Text copyable>https://verify.com/{selectedCertificate.certificateId}</Text>
          </div>
        </Modal>
      )}
    </div>
  );
};

export default Dashboard;

