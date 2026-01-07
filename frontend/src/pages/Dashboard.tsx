import { Card, Row, Col, Typography } from 'antd';

const { Title } = Typography;

export default function Dashboard() {
  return (
    <div>
      <Title level={2}>Dashboard</Title>
      <Row gutter={16}>
        <Col span={8}>
          <Card title="Clients" bordered={false}>
            <p>Total: 0</p>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="Projects" bordered={false}>
            <p>Total: 0</p>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="Tasks" bordered={false}>
            <p>Total: 0</p>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
