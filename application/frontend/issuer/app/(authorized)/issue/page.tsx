'use client';
import React, { useState } from 'react';
import { Button, Input, Card, Typography, Form, message} from 'antd';
import { v4 as uuidv4 } from 'uuid'; // To generate unique certificate IDs

const { Title} = Typography;

// Define types for certificate data
interface Certificate {
  certificateId: string;
  recipientName: string;
  issueDate: string;
  validity: string;
}

const IssueCertificate: React.FC = () => {
  const [newCertificate, setNewCertificate] = useState<Certificate>({
    certificateId: '',
    recipientName: '',
    issueDate: '',
    validity: '',
  });
  const [certificates, setCertificates] = useState<Certificate[]>([]);

  // Handle form input change
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>, field: string) => {
    setNewCertificate({ ...newCertificate, [field]: e.target.value });
  };

  // Generate a new certificate and add it to the list
  const handleIssueCertificate = () => {
    if (!newCertificate.recipientName || !newCertificate.issueDate || !newCertificate.validity) {
      message.error('Please fill out all the fields.');
      return;
    }

    const newCert: Certificate = {
      ...newCertificate,
      certificateId: uuidv4(), // Generate a unique certificate ID
    };

    setCertificates([...certificates, newCert]);
    setNewCertificate({
      certificateId: '',
      recipientName: '',
      issueDate: '',
      validity: '',
    });
    message.success('Certificate issued successfully!');
  };

  return (
    <div style={{ padding: '30px', maxWidth: '800px', margin: 'auto' }}>
      <Title level={2}>Issue a New Certificate</Title>

      <Card title="Create New Certificate" style={{ marginBottom: '20px' }}>
        <Form layout="vertical">
          <Form.Item label="Recipient Name">
            <Input
              value={newCertificate.recipientName}
              onChange={(e) => handleInputChange(e, 'recipientName')}
              placeholder="Enter the recipient's name"
            />
          </Form.Item>

          <Form.Item label="Issue Date">
            <Input
              value={newCertificate.issueDate}
              onChange={(e) => handleInputChange(e, 'issueDate')}
              placeholder="Enter the issue date"
            />
          </Form.Item>

          <Form.Item label="Validity">
            <Input
              value={newCertificate.validity}
              onChange={(e) => handleInputChange(e, 'validity')}
              placeholder="Enter validity period (e.g., Lifetime, 1 Year)"
            />
          </Form.Item>

          <Button
            type="primary"
            onClick={handleIssueCertificate}
            style={{ width: '100%' }}
          >
            Issue Certificate
          </Button>
        </Form>
      </Card>
    </div>
  );
};

export default IssueCertificate;

