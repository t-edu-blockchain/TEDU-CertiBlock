"use client"
import { Typography, Statistic, Row, Col, Tooltip, Button, Card } from 'antd';
import { Bar, BarChart, CartesianGrid, Legend, ResponsiveContainer, XAxis, YAxis } from 'recharts';
import moment from 'moment'; // For date manipulation
const { Title } = Typography;

// Define types for certificate data
interface Certificate {
  certificateId: string;
  certificateName: string;
  issueDate: string;
  expirationDate: string;  // Updated to a specific expiration date
}

// Helper function to check if certificate is expired
const isExpired = (expirationDate: string, currentDate: moment.Moment) => {
  return moment(expirationDate).isBefore(currentDate); // Return true if expired
};

export default function Home() {
  const certificates: Certificate[] = [
    {
      certificateId: '123456',
      certificateName: 'A',
      issueDate: 'March 2025',
      expirationDate: 'March 2025',
    },
    {
      certificateId: '789012',
      certificateName: 'B',
      issueDate: 'January 2025',
      expirationDate: 'January 2025',
    },
    {
      certificateId: '345678',
      certificateName: 'C',
      issueDate: 'July 2024',
      expirationDate: 'July 2024',
    },
    {
      certificateId: '345679',
      certificateName: 'D',
      issueDate: 'July 2023',
      expirationDate: 'July 2024',
    },
    {
      certificateId: '901234',
      certificateName: 'E',
      issueDate: 'December 2024',
      expirationDate: 'December 2024',
    },
  ];

  // Get current date for comparison
  const currentDate = moment();

  // Determine expired certificates
  const expiredCertificates = certificates.filter((certificate) => {
    return isExpired(certificate.expirationDate, currentDate); // Check if the certificate is expired
  });

  // Chart data for number of certificates per expiration date
  const chartData = [
    { name: 'Expired', certificates: expiredCertificates.length },
    { name: 'Valid', certificates: certificates.length - expiredCertificates.length },
  ];

  return (
    <div style={{ padding: '20px' }}>
      <Title level={2}>Certificates Overview</Title>

      {/* General Statistics Section */}
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic title="Total Certificates" value={certificates.length} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="Valid Certificates" value={certificates.length - expiredCertificates.length} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="Expired Certificates" value={expiredCertificates.length} />
          </Card>
        </Col>
      </Row>

      {/* Expiry Status Breakdown */}
      <Title level={4} style={{ marginTop: '20px' }}>Certificate Expiration Breakdown</Title>
      <Row gutter={16}>
        <Col span={24}>
          <Card>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="certificates" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      {/* Expired Certificates List */}
      <Title level={4} style={{ marginTop: '20px' }}>Expired Certificates</Title>
      {expiredCertificates.length > 0 ? (
        <ul>
          {expiredCertificates.map((certificate) => (
            <li key={certificate.certificateId}>
              <Tooltip title={`Issued: ${moment(certificate.issueDate).format('MMMM YYYY')}`}>
                {certificate.certificateName} (Expired on {moment(certificate.expirationDate).format('MMMM YYYY')})
              </Tooltip>
            </li>
          ))}
        </ul>
      ) : (
        <p>No expired certificates.</p>
      )}

      {/* Button for viewing or exporting expired certificates */}
      {expiredCertificates.length > 0 && (
        <Button type="primary" style={{ marginTop: '20px' }} onClick={() => alert('Exporting expired certificates...')}>
          Export Expired Certificates
        </Button>
      )}
    </div>
  );
}

