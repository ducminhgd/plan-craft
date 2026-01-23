import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message, Input, Select, Modal } from 'antd';
import { PlusOutlined, EditOutlined, SyncOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetProjectResources, UpdateProjectResource, DeleteProjectResource, GetProjects, GetHumanResources } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import type { TableRowSelection } from 'antd/es/table/interface';
import { formatDate } from '../../utils/date';
import { useDatabase } from '../../contexts/DatabaseContext';

const { Title } = Typography;
const { Search } = Input;

export default function ResourceAllocationList() {
  const navigate = useNavigate();
  const { refreshKey } = useDatabase();
  const [allocations, setAllocations] = useState<entities.ProjectResource[]>([]);
  const [projects, setProjects] = useState<entities.Project[]>([]);
  const [humanResources, setHumanResources] = useState<entities.HumanResource[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Filter states
  const [searchText, setSearchText] = useState('');
  const [statusFilter, setStatusFilter] = useState<number | undefined>(0);
  const [projectFilter, setProjectFilter] = useState<number | undefined>(0);
  const [humanResourceFilter, setHumanResourceFilter] = useState<number | undefined>(0);

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const loadProjects = async () => {
    try {
      const params: any = {
        status: 2, // Only active projects
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      const result = await GetProjects(params);
      setProjects(result.data || []);
    } catch (error) {
      console.error('Failed to load projects', error);
    }
  };

  const loadHumanResources = async () => {
    try {
      const params: any = {
        status: 2, // Only active human resources
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      const result = await GetHumanResources(params);
      setHumanResources(result.data || []);
    } catch (error) {
      console.error('Failed to load human resources', error);
    }
  };

  const loadAllocations = async () => {
    setLoading(true);
    try {
      const params: any = {
        pagination: {
          page: currentPage,
          page_size: pageSize,
          total: 0,
        },
      };

      // Apply search filter
      if (searchText) {
        params.role_like = searchText;
      }

      // Apply status filter
      if (statusFilter !== 0 && statusFilter !== undefined) {
        params.status = statusFilter;
      }

      // Apply project filter
      if (projectFilter !== 0 && projectFilter !== undefined) {
        params.project_id = projectFilter;
      }

      // Apply human resource filter
      if (humanResourceFilter !== 0 && humanResourceFilter !== undefined) {
        params.human_resource_id = humanResourceFilter;
      }

      const result = await GetProjectResources(params);
      setAllocations(result.data || []);
      setTotal(result.total || 0);
    } catch (error) {
      message.error('Failed to load resource allocations');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProjects();
    loadHumanResources();
  }, [refreshKey]);

  useEffect(() => {
    loadAllocations();
  }, [currentPage, pageSize, searchText, statusFilter, projectFilter, humanResourceFilter, refreshKey]);

  const handleToggleStatus = async (allocation: entities.ProjectResource) => {
    const newStatus = allocation.status === 2 ? 1 : 2;
    const statusText = newStatus === 2 ? 'activate' : 'deactivate';

    Modal.confirm({
      title: `Confirm ${statusText.charAt(0).toUpperCase() + statusText.slice(1)}`,
      content: `Are you sure you want to ${statusText} this allocation?`,
      okText: statusText.charAt(0).toUpperCase() + statusText.slice(1),
      okType: newStatus === 1 ? 'danger' : 'primary',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          const updatedAllocation = Object.assign(Object.create(Object.getPrototypeOf(allocation)), allocation, { status: newStatus });
          await UpdateProjectResource(updatedAllocation);
          message.success(`Allocation ${statusText}d successfully`);
          loadAllocations();
        } catch (error) {
          message.error(`Failed to ${statusText} allocation`);
        }
      },
    });
  };

  const handleDelete = async (allocation: entities.ProjectResource) => {
    Modal.confirm({
      title: 'Confirm Delete',
      content: 'Are you sure you want to delete this allocation? This action cannot be undone.',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await DeleteProjectResource(allocation.id);
          message.success('Allocation deleted successfully');
          loadAllocations();
        } catch (error) {
          message.error('Failed to delete allocation');
        }
      },
    });
  };

  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrentPage(1);
  };

  const handleStatusChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setStatusFilter(value);
    setCurrentPage(1);
  };

  const handleProjectChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setProjectFilter(value);
    setCurrentPage(1);
  };

  const handleHumanResourceChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setHumanResourceFilter(value);
    setCurrentPage(1);
  };

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setCurrentPage(1);
  };

  const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
    setSelectedRowKeys(newSelectedRowKeys);
  };

  const rowSelection: TableRowSelection<entities.ProjectResource> = {
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

  const getProjectName = (projectId: number) => {
    const project = projects.find(p => p.id === projectId);
    return project ? project.name : '-';
  };

  const getHumanResourceName = (humanResourceId: number) => {
    const hr = humanResources.find(h => h.id === humanResourceId);
    return hr ? hr.name : '-';
  };

  const getProjectCurrency = (projectId: number) => {
    const project = projects.find(p => p.id === projectId);
    return project?.currency || '';
  };

  const formatCost = (cost: number, projectId: number) => {
    const currency = getProjectCurrency(projectId);
    const formattedCost = cost?.toLocaleString() ?? '0';
    return currency ? `${formattedCost} ${currency}` : formattedCost;
  };

  const columns = [
    {
      title: 'Project',
      dataIndex: 'project_id',
      key: 'project_id',
      render: (projectId: number) => (
        <a onClick={() => navigate(`/projects/${projectId}`)} style={{ cursor: 'pointer' }}>
          {getProjectName(projectId)}
        </a>
      ),
    },
    {
      title: 'Human Resource',
      dataIndex: 'human_resource_id',
      key: 'human_resource_id',
      render: (hrId: number) => (
        <a onClick={() => navigate(`/human-resources/${hrId}`)} style={{ cursor: 'pointer' }}>
          {getHumanResourceName(hrId)}
        </a>
      ),
    },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
    },
    {
      title: 'Allocation (%)',
      dataIndex: 'allocation',
      key: 'allocation',
      render: (allocation: number) => `${allocation}%`,
    },
    {
      title: 'Cost',
      dataIndex: 'cost',
      key: 'cost',
      render: (_: number, record: entities.ProjectResource) => formatCost(record.cost, record.project_id),
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
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => getStatusText(status),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: entities.ProjectResource) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/resource-allocations/${record.id}`)}
            size="small"
          />
          <Button
            icon={<SyncOutlined />}
            type={record.status === 2 ? 'default' : 'primary'}
            danger={record.status === 2}
            onClick={() => handleToggleStatus(record)}
            size="small"
            title={record.status === 2 ? 'Deactivate allocation' : 'Activate allocation'}
          />
          <Button
            icon={<DeleteOutlined />}
            danger
            onClick={() => handleDelete(record)}
            size="small"
            title="Delete allocation"
          />
        </Space>
      ),
    },
  ];

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
      for (let i = 1; i <= totalPages; i++) {
        pages.push(renderPageButton(i));
      }
    } else {
      pages.push(renderPageButton(1));

      if (currentPage > 3) {
        pages.push(renderEllipsis('ellipsis-start'));
      }

      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(renderPageButton(i));
      }

      if (currentPage < totalPages - 2) {
        pages.push(renderEllipsis('ellipsis-end'));
      }

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
        <Title level={2}>List of Project Allocations</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/resource-allocations/new')}
        >
          Add New
        </Button>
      </div>

      {/* Filter Section */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16, flexWrap: 'wrap' }}>
        <Search
          placeholder="Search by role"
          onSearch={handleSearch}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ flex: 1, minWidth: 200 }}
          allowClear
        />
        <Select
          style={{ width: 200 }}
          onChange={handleProjectChange}
          value={projectFilter}
          showSearch
          optionFilterProp="children"
          filterOption={(input, option) =>
            (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
          }
        >
          <Select.Option value={0}>-- All Projects --</Select.Option>
          {projects.map(project => (
            <Select.Option key={project.id} value={project.id}>{project.name}</Select.Option>
          ))}
        </Select>
        <Select
          style={{ width: 200 }}
          onChange={handleHumanResourceChange}
          value={humanResourceFilter}
          showSearch
          optionFilterProp="children"
          filterOption={(input, option) =>
            (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
          }
        >
          <Select.Option value={0}>-- All Resources --</Select.Option>
          {humanResources.map(hr => (
            <Select.Option key={hr.id} value={hr.id}>{hr.name}</Select.Option>
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
        dataSource={allocations}
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
