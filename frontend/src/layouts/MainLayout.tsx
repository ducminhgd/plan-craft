import { Layout } from 'antd';
import { Outlet } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import './MainLayout.css';

const { Content } = Layout;

export default function MainLayout() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Layout className="main-container">
        <Sidebar />
        <Layout className="content-wrapper">
          <Content className="content-container">
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}
