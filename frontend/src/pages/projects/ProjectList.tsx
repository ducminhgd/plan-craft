import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message, Input, Select, Modal } from 'antd';
import { PlusOutlined, EditOutlined, SyncOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetProjects, UpdateProject, GetClients } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import type { TableRowSelection } from 'antd/es/table/interface';
import { formatDate } from '../../utils/date';
import { useDatabase } from '../../contexts/DatabaseContext';

const { Title } = Typography;
const { Search } = Input;

export default function ProjectList() {
  const navigate = useNavigate();
  const { refreshKey } = useDatabase();
  const [projects, setProjects] = useState<entities.Project[]>([]);
  const [clients, setClients] = useState<entities.Client[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Filter states
  const [searchText, setSearchText] = useState('');
  const [statusFilter, setStatusFilter] = useState<number | undefined>(0);
  const [clientFilter, setClientFilter] = useState<number | undefined>(0);

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const loadClients = async () => {
    try {
      const params: any = {
        status: 2, // Only active clients
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      const result = await GetClients(params);
      setClients(result.data || []);
    } catch (error) {
      console.error('Failed to load clients', error);
    }
  };

  const loadProjects = async () => {
    setLoading(true);
    try {
      const params: any = {
        pagination: {
          page: currentPage,
          page_size: pageSize,
          total: 0,
        },
      };

      // Apply search filter across multiple fields
      if (searchText) {
        params.name_like = searchText;
        params.description_like = searchText;
      }

      // Apply status filter
      if (statusFilter !== 0 && statusFilter !== undefined) {
        params.status = statusFilter;
      }

      // Apply client filter
      if (clientFilter !== 0 && clientFilter !== undefined) {
        params.client_id = clientFilter;
      }

      const result = await GetProjects(params);
      setProjects(result.data || []);
      setTotal(result.total || 0);
    } catch (error) {
      message.error('Failed to load projects');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClients();
  }, [refreshKey]);

  useEffect(() => {
    loadProjects();
  }, [currentPage, pageSize, searchText, statusFilter, clientFilter, refreshKey]);

  const handleToggleStatus = async (project: entities.Project) => {
    const newStatus = project.status === 2 ? 1 : 2;
    const statusText = newStatus === 2 ? 'activate' : 'deactivate';

    Modal.confirm({
      title: `Confirm ${statusText.charAt(0).toUpperCase() + statusText.slice(1)}`,
      content: `Are you sure you want to ${statusText} this project?`,
      okText: statusText.charAt(0).toUpperCase() + statusText.slice(1),
      okType: newStatus === 1 ? 'danger' : 'primary',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          // Create a new project instance with updated status
          const updatedProject = Object.assign(Object.create(Object.getPrototypeOf(project)), project, { status: newStatus });
          await UpdateProject(updatedProject);
          message.success(`Project ${statusText}d successfully`);
          loadProjects();
        } catch (error) {
          message.error(`Failed to ${statusText} project`);
        }
      },
    });
  };

  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrentPage(1); // Reset to first page on new search
  };

  const handleStatusChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setStatusFilter(value);
    setCurrentPage(1); // Reset to first page on filter change
  };

  const handleClientChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setClientFilter(value);
    setCurrentPage(1); // Reset to first page on filter change
  };

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setCurrentPage(1); // Reset to first page on page size change
  };

  const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
    setSelectedRowKeys(newSelectedRowKeys);
  };

  const rowSelection: TableRowSelection<entities.Project> = {
    selectedRowKeys,
    onChange: onSelectChange,
  };

  const getStatusText = (status: number) => {
    switch (status) {
      case 2:
        return <span style={{ color: '#52c41a' }}>Active</span>;
      case 1:
        return <span style={{ color: '#ff4d4f' }}>Inactive</span>;
      default:
        return <span style={{ color: '#999' }}>Unknown</span>;
    }
  };

  const getClientName = (clientId: number) => {
    const client = clients.find(c => c.id === clientId);
    return client ? client.name : '-';
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: entities.Project) => (
        <a onClick={() => navigate(`/projects/${record.id}`)} style={{ cursor: 'pointer' }}>
          {name}
        </a>
      ),
    },
    {
      title: 'Client',
      dataIndex: 'client_id',
      key: 'client_id',
      render: (clientId: number) => getClientName(clientId),
    },
    {
      title: 'Start Date',
      dataIndex: 'start_date',
      key: 'start_date',
      render: (date: any) => formatDate(date),
    },
    {
      title: 'End Date',
      dataIndex: 'end_date',
      key: 'end_date',
      render: (date: any) => formatDate(date),
    },
    {
      title: 'Currency',
      dataIndex: 'currency',
      key: 'currency',
      render: (currency: string) => currency || '-',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => getStatusText(status),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: entities.Project) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/projects/${record.id}`)}
            size="small"
          />
          <Button
            icon={<SyncOutlined />}
            type={record.status === 2 ? 'default' : 'primary'}
            danger={record.status === 2}
            onClick={() => handleToggleStatus(record)}
            size="small"
            title={record.status === 2 ? 'Deactivate project' : 'Activate project'}
          />
        </Space>
      ),
    },
  ];

  // Custom pagination renderer
  const renderPagination = () => {
    const totalPages = Math.ceil(total / pageSize);
    if (totalPages <= 1) return null;

    const renderPageButton = (page: number, label?: string) => (
      <Button
        key={page}
        type={currentPage === page ? 'primary' : 'default'}
        size="small"
        onClick={() => setCurrentPage(page)}
        style={{ margin: '0 4px' }}
      >
        {label || page}
      </Button>
    );

    const renderEllipsis = (key: string) => (
      <span key={key} style={{ margin: '0 4px' }}>
        ..
      </span>
    );

    const pages: React.ReactNode[] = [];

    if (totalPages <= 7) {
      // Show all pages if 7 or less
      for (let i = 1; i <= totalPages; i++) {
        pages.push(renderPageButton(i));
      }
    } else {
      // Always show first page
      pages.push(renderPageButton(1));

      if (currentPage > 3) {
        pages.push(renderEllipsis('ellipsis-start'));
      }

      // Show pages around current page
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(renderPageButton(i));
      }

      if (currentPage < totalPages - 2) {
        pages.push(renderEllipsis('ellipsis-end'));
      }

      // Always show last page
      pages.push(renderPageButton(totalPages));
    }

    return (
      <Space>
        <Button
          size="small"
          disabled={currentPage === 1}
          onClick={() => setCurrentPage(currentPage - 1)}
        >
          Previous
        </Button>
        {pages}
        <Button
          size="small"
          disabled={currentPage === totalPages}
          onClick={() => setCurrentPage(currentPage + 1)}
        >
          Next
        </Button>
      </Space>
    );
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={2}>List of Projects</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/projects/new')}
        >
          Add New
        </Button>
      </div>

      {/* Filter Section */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Search
          placeholder="Search by name or description"
          onSearch={handleSearch}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ flex: 1 }}
          allowClear
        />
        <Select
          style={{ width: 200 }}
          onChange={handleClientChange}
          value={clientFilter}
          showSearch
          optionFilterProp="children"
          filterOption={(input, option) =>
            (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
          }
        >
          <Select.Option value={0}>-- All Clients --</Select.Option>
          {clients.map(client => (
            <Select.Option key={client.id} value={client.id}>{client.name}</Select.Option>
          ))}
        </Select>
        <Select
          style={{ width: 150 }}
          onChange={handleStatusChange}
          value={statusFilter}
        >
          <Select.Option value={0}>-- Status --</Select.Option>
          <Select.Option value={2}>Active</Select.Option>
          <Select.Option value={1}>Inactive</Select.Option>
        </Select>
      </div>

      <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={projects}
        loading={loading}
        rowKey="id"
        pagination={false}
      />

      {/* Page Size Selector and Pagination */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <Select
            value={pageSize}
            onChange={handlePageSizeChange}
            style={{ width: 80 }}
          >
            <Select.Option value={10}>10</Select.Option>
            <Select.Option value={20}>20</Select.Option>
            <Select.Option value={50}>50</Select.Option>
            <Select.Option value={100}>100</Select.Option>
          </Select>
          <span>records per page</span>
        </div>
        {renderPagination()}
      </div>
    </div>
  );
}
