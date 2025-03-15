'use client';
import React, { useState } from 'react';
import { Button, Modal, QRCode, Table, Typography, Input } from 'antd';

const { Title, Text } = Typography;
const { Search } = Input;

// Define types for certificate data
interface Certificate {
  certificateId: string;
  recipientName: string;
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
      recipientName: 'John Doe',
      issueDate: 'March 2025',
      validity: 'Lifetime',
    },
    {
      certificateId: '789012',
      recipientName: 'Emily Johnson',
      issueDate: 'January 2025',
      validity: 'Lifetime',
    },
    {
      certificateId: '345678',
      recipientName: 'Michael Brown',
      issueDate: 'July 2024',
      validity: '1 Year',
    },
    {
      certificateId: '901234',
      recipientName: 'Sophia Lee',
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
        title: 'Recipient Name',
        dataIndex: 'recipientName',
        key: 'recipientName',
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
      certificate.recipientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
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
      <Title level={2}>Certificates Issued</Title>

      {/* Search Input */}
      <Search
        placeholder="Search by Recipient Name, or Certificate ID"
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
            </Button>
          ]}
        >
          <QRCode value={`https://example.com/certificate/${selectedCertificate.certificateId}`} />
          <Text>Scan this QR code to view or share your certificate.</Text>
        </Modal>
      )}
    </div>
  );
};

export default Dashboard;

