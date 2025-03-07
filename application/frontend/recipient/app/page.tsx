import { Button } from 'antd';
import Link from 'next/link';

export default function Home() {
  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      <h1>Wellcome to recipient website, please login or register.</h1>

      <div style={{ marginTop: '20px' }}>
        <Link href="/login" passHref>
          <Button type="primary" style={{ marginRight: '10px' }}>
            Login
          </Button>
        </Link>

        <Link href="/register" passHref>
          <Button type="default">
            Register
          </Button>
        </Link>
      </div>
    </div>
  );
}

