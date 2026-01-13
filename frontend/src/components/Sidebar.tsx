import { useState, useEffect } from 'react';
import { Layout, Menu } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  HomeOutlined,
  TeamOutlined,
  ProjectOutlined,
  PlusOutlined,
  UnorderedListOutlined,
  PushpinOutlined,
} from '@ant-design/icons';

const { Sider } = Layout;

const STORAGE_KEY = 'sidebar-collapsed';
const PIN_STORAGE_KEY = 'sidebar-pinned';

export default function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  const [collapsed, setCollapsed] = useState(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    return saved ? JSON.parse(saved) : false;
  });

  const [pinned, setPinned] = useState(() => {
    const saved = localStorage.getItem(PIN_STORAGE_KEY);
    return saved ? JSON.parse(saved) : false;
  });

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(collapsed));
  }, [collapsed]);

  useEffect(() => {
    localStorage.setItem(PIN_STORAGE_KEY, JSON.stringify(pinned));
  }, [pinned]);

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: 'Dashboard',
    },
    {
      key: 'clients',
      icon: <TeamOutlined />,
      label: 'Client Management',
      children: [
        {
          key: '/clients',
          icon: <UnorderedListOutlined />,
          label: 'List',
        },
        {
          key: '/clients/new',
          icon: <PlusOutlined />,
          label: 'Add New',
        },
      ],
    },
    {
      key: 'projects',
      icon: <ProjectOutlined />,
      label: 'Project Management',
      children: [
        {
          key: '/projects',
          icon: <UnorderedListOutlined />,
          label: 'List',
        },
        {
          key: '/projects/new',
          icon: <PlusOutlined />,
          label: 'Add New',
        },
      ],
    },
  ];

  return (
    <Sider
      collapsible
      collapsed={collapsed}
      onCollapse={setCollapsed}
      collapsedWidth={80}
      breakpoint="lg"
      trigger={pinned ? null : undefined}
    >
      <div style={{ height: 32, margin: 16, background: 'rgba(255,255,255,.2)' }}>
        {/* Logo placeholder */}
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[location.pathname]}
        defaultOpenKeys={['clients', 'projects']}
        items={menuItems}
        onClick={({ key }) => navigate(key)}
      />
      <div style={{ position: 'absolute', bottom: 16, left: 16 }}>
        <PushpinOutlined
          style={{ color: pinned ? '#1890ff' : '#fff', cursor: 'pointer' }}
          onClick={() => setPinned(!pinned)}
        />
      </div>
    </Sider>
  );
}
