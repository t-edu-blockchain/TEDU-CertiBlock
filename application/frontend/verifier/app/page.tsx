'use client';
import React, { useState } from 'react';
import { Button, Input, Typography, Row, Col, Card, message, Modal, QRCode } from 'antd';

const { Title, Text } = Typography;

const Verifier: React.FC = () => {
  const [secretCode, setSecretCode] = useState<string>('');
  const [qrCode, setQrCode] = useState<string | null>(null);
  const [verificationResult, setVerificationResult] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

  // Handle secret code input change
  const handleSecretCodeChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setSecretCode(e.target.value);
  };

  // Handle QR Code input change (Simulating upload for the demo)
  const handleUploadQrCode = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const file = e.target.files && e.target.files[0];
    if (file) {
      // Simulate QR code upload (you can use a library to decode QR codes)
      setQrCode('validQRCodeData'); // Assume QR code data is valid here
    }
  };

  // Verification logic based on either secret code or QR code
  const handleVerify = (): void => {
    if (!secretCode && !qrCode) {
      message.error('Please either upload a QR code or enter a secret code to verify.');
      return;
    }

    // Simulate certificate verification logic (can be adjusted to your actual logic)
    if (secretCode === 'validSecretCode' || qrCode === 'validQRCodeData') {
      setVerificationResult('Certificate is valid.');
    } else {
      setVerificationResult('Invalid certificate.');
    }
  };

  // Open modal to show certificate details when verified
  const handleModalOpen = (): void => {
    setIsModalOpen(true);
  };

  const handleModalClose = (): void => {
    setIsModalOpen(false);
  };

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>Verifier - Validate Certificate</Title>

      <Row gutter={[16, 16]} style={{ marginBottom: '20px' }}>
        {/* Secret Code Input */}
        <Col span={12}>
          <Card>
            <Title level={4}>Enter Secret Code</Title>
            <Input
              value={secretCode}
              onChange={handleSecretCodeChange}
              placeholder="Enter secret code provided by the issuer"
              style={{ marginBottom: '10px' }}
            />
            <Button type="primary" onClick={handleVerify}>
              Verify Secret Code
            </Button>
          </Card>
        </Col>

        {/* QR Code Upload */}
        <Col span={12}>
          <Card>
            <Title level={4}>Upload QR Code</Title>
            <input type="file" accept="image/*" onChange={handleUploadQrCode} />
            <Button type="primary" onClick={handleVerify} style={{ marginTop: '10px' }}>
              Verify QR Code
            </Button>
          </Card>
        </Col>
      </Row>

      {/* Verification Result */}
      {verificationResult && (
        <Card>
          <Text strong>{verificationResult}</Text>
        </Card>
      )}

      {/* Modal to show certificate details */}
      <Modal
        title="Certificate Details"
        open={isModalOpen}
        onCancel={handleModalClose}
        footer={null}
      >
        <QRCode value="https://verify.com/certificate/123456" />
        <div style={{ marginTop: '10px' }}>
          <Text copyable>https://verify.com/certificate/123456</Text>
        </div>
      </Modal>
    </div>
  );
};

export default Verifier;

