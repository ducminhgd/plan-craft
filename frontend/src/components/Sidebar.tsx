import { useState, useEffect } from 'react';
import { Layout, Menu, Tooltip } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  HomeOutlined,
  TeamOutlined,
  ProjectOutlined,
  PlusOutlined,
  UnorderedListOutlined,
  PushpinOutlined,
  UserOutlined,
  FlagOutlined,
  AppstoreOutlined,
  ApartmentOutlined,
  ScheduleOutlined,
  CheckSquareOutlined,
  DatabaseOutlined,
} from '@ant-design/icons';
import { useDatabase } from '../contexts/DatabaseContext';

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

  const { currentDbPath, isDraft } = useDatabase();

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
      label: 'Client',
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
      key: 'resources',
      icon: <AppstoreOutlined />,
      label: 'Resources',
      children: [
        {
          key: 'human-resources',
          icon: <UserOutlined />,
          label: 'Human',
          children: [
            {
              key: '/human-resources',
              icon: <UnorderedListOutlined />,
              label: 'List',
            },
            {
              key: '/human-resources/new',
              icon: <PlusOutlined />,
              label: 'Add New',
            },
          ],
        },
      ],
    },
    {
      key: 'projects',
      icon: <ProjectOutlined />,
      label: 'Project',
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
    {
      key: 'project-plan',
      icon: <ApartmentOutlined />,
      label: 'Project Plan',
      children: [
        {
          key: 'milestones',
          icon: <FlagOutlined />,
          label: 'Milestones',
          children: [
            {
              key: '/milestones',
              icon: <UnorderedListOutlined />,
              label: 'List',
            },
            {
              key: '/milestones/new',
              icon: <PlusOutlined />,
              label: 'Add New',
            },
          ],
        },
        {
          key: 'tasks',
          icon: <CheckSquareOutlined />,
          label: 'Tasks',
          children: [
            {
              key: '/tasks',
              icon: <UnorderedListOutlined />,
              label: 'List',
            },
            {
              key: '/tasks/new',
              icon: <PlusOutlined />,
              label: 'Add New',
            },
          ],
        },
      ],
    },
    {
      key: 'resource-allocations',
      icon: <ScheduleOutlined />,
      label: 'Project Allocation',
      children: [
        {
          key: '/resource-allocations',
          icon: <UnorderedListOutlined />,
          label: 'List',
        },
        {
          key: '/resource-allocations/new',
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
      style={{ height: '100vh', position: 'sticky', top: 0, left: 0 }}
    >
      <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
        <Tooltip title={isDraft ? 'Unsaved draft - use Save As to persist' : currentDbPath} placement="right">
          <div style={{
            minHeight: 32,
            margin: 16,
            padding: '8px',
            background: 'rgba(255,255,255,.2)',
            flexShrink: 0,
            borderRadius: 4,
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            overflow: 'hidden',
          }}>
            <DatabaseOutlined style={{ color: isDraft ? '#faad14' : '#52c41a', fontSize: 16, flexShrink: 0 }} />
            {!collapsed && (
              <span style={{
                color: isDraft ? '#faad14' : '#fff',
                fontSize: 12,
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
              }}>
                {isDraft ? 'Draft' : currentDbPath.split(/[/\\]/).pop() || 'No database'}
              </span>
            )}
          </div>
        </Tooltip>
        <div style={{ flex: 1, overflowY: 'auto', overflowX: 'hidden' }}>
          <Menu
            theme="dark"
            mode="inline"
            selectedKeys={[location.pathname]}
            defaultOpenKeys={['clients', 'resources', 'human-resources', 'projects', 'project-plan', 'milestones', 'tasks', 'resource-allocations']}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            style={{ textAlign: 'left' }}
          />
        </div>
        <div style={{ padding: 16, flexShrink: 0 }}>
          <PushpinOutlined
            style={{ color: pinned ? '#1890ff' : '#fff', cursor: 'pointer' }}
            onClick={() => setPinned(!pinned)}
          />
        </div>
      </div>
    </Sider>
  );
}
