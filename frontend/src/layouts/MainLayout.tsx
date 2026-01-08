import { Layout } from 'antd';
import { Outlet } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import MenuBar from '../components/MenuBar';

const { Content } = Layout;

export default function MainLayout() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <MenuBar />
      <Layout style={{ minHeight: 'calc(100vh - 40px)' }}>
        <Sidebar />
        <Layout>
          <Content style={{ margin: '24px 16px', padding: 24, background: '#fff' }}>
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}
