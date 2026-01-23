import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message, Input, Select, Modal } from 'antd';
import { PlusOutlined, EditOutlined, SyncOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetHumanResources, UpdateHumanResource } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import type { TableRowSelection } from 'antd/es/table/interface';
import { useDatabase } from '../../contexts/DatabaseContext';

const { Title } = Typography;
const { Search } = Input;

export default function HumanResourceList() {
  const navigate = useNavigate();
  const { refreshKey } = useDatabase();
  const [humanResources, setHumanResources] = useState<entities.HumanResource[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Filter states
  const [searchText, setSearchText] = useState('');
  const [statusFilter, setStatusFilter] = useState<number | undefined>(0);

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const loadHumanResources = async () => {
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
        params.title_like = searchText;
        params.level_like = searchText;
      }

      // Apply status filter
      if (statusFilter !== 0 && statusFilter !== undefined) {
        params.status = statusFilter;
      }

      const result = await GetHumanResources(params);
      setHumanResources(result.data || []);
      setTotal(result.total || 0);
    } catch (error) {
      message.error('Failed to load human resources');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadHumanResources();
  }, [currentPage, pageSize, searchText, statusFilter, refreshKey]);

  const handleToggleStatus = async (hr: entities.HumanResource) => {
    const newStatus = hr.status === 2 ? 1 : 2;
    const statusText = newStatus === 2 ? 'activate' : 'deactivate';

    Modal.confirm({
      title: `Confirm ${statusText.charAt(0).toUpperCase() + statusText.slice(1)}`,
      content: `Are you sure you want to ${statusText} this human resource?`,
      okText: statusText.charAt(0).toUpperCase() + statusText.slice(1),
      okType: newStatus === 1 ? 'danger' : 'primary',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          const updatedHR = Object.assign(Object.create(Object.getPrototypeOf(hr)), hr, { status: newStatus });
          await UpdateHumanResource(updatedHR);
          message.success(`Human resource ${statusText}d successfully`);
          loadHumanResources();
        } catch (error) {
          message.error(`Failed to ${statusText} human resource`);
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

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setCurrentPage(1);
  };

  const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
    setSelectedRowKeys(newSelectedRowKeys);
  };

  const rowSelection: TableRowSelection<entities.HumanResource> = {
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

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: entities.HumanResource) => (
        <a onClick={() => navigate(`/human-resources/${record.id}`)} style={{ cursor: 'pointer' }}>
          {name}
        </a>
      ),
    },
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: 'Level',
      dataIndex: 'level',
      key: 'level',
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
      render: (_: any, record: entities.HumanResource) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/human-resources/${record.id}`)}
            size="small"
          />
          <Button
            icon={<SyncOutlined />}
            type={record.status === 2 ? 'default' : 'primary'}
            danger={record.status === 2}
            onClick={() => handleToggleStatus(record)}
            size="small"
            title={record.status === 2 ? 'Deactivate' : 'Activate'}
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
        <Title level={2}>List of Human Resources</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/human-resources/new')}
        >
          Add New
        </Button>
      </div>

      {/* Filter Section */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Search
          placeholder="Search by name, title, level"
          onSearch={handleSearch}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ flex: 1 }}
          allowClear
        />
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
        dataSource={humanResources}
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
