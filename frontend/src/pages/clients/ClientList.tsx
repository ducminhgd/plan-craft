import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message, Input, Select, Checkbox, Modal } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetClients, DeleteClient } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import type { CheckboxChangeEvent } from 'antd/es/checkbox';
import type { TableRowSelection } from 'antd/es/table/interface';

const { Title } = Typography;
const { Search } = Input;

export default function ClientList() {
  const navigate = useNavigate();
  const [clients, setClients] = useState<entities.Client[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Filter states
  const [searchText, setSearchText] = useState('');
  const [statusFilter, setStatusFilter] = useState<number | undefined>(undefined);

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const loadClients = async () => {
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
        params.email_like = searchText;
        params.phone_like = searchText;
        params.address_like = searchText;
        params.contact_person_like = searchText;
        params.notes_like = searchText;
      }

      // Apply status filter
      if (statusFilter !== undefined) {
        params.status = statusFilter;
      }

      const result = await GetClients(params);
      setClients(result.data || []);
      setTotal(result.total || 0);
    } catch (error) {
      message.error('Failed to load clients');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClients();
  }, [currentPage, pageSize, searchText, statusFilter]);

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: 'Confirm Delete',
      content: 'Are you sure you want to delete this client?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await DeleteClient(id);
          message.success('Client deleted');
          loadClients();
        } catch (error) {
          message.error('Failed to delete client');
        }
      },
    });
  };

  const handleSearch = (value: string) => {
    setSearchText(value);
    setCurrentPage(1); // Reset to first page on new search
  };

  const handleStatusChange = (value: number | undefined) => {
    setStatusFilter(value);
    setCurrentPage(1); // Reset to first page on filter change
  };

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setCurrentPage(1); // Reset to first page on page size change
  };

  const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
    setSelectedRowKeys(newSelectedRowKeys);
  };

  const rowSelection: TableRowSelection<entities.Client> = {
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
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Phone',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: 'Contact Person',
      dataIndex: 'contact_person',
      key: 'contact_person',
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
      render: (_: any, record: entities.Client) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/clients/${record.id}`)}
            size="small"
          />
          <Button
            icon={<DeleteOutlined />}
            danger
            onClick={() => handleDelete(record.id)}
            size="small"
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
        <Title level={2}>List of Clients</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/clients/new')}
        >
          Add New
        </Button>
      </div>

      {/* Filter Section */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16 }}>
        <Search
          placeholder="Search by name, email, phone, address, contact person, notes"
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
          <Select.Option value={undefined}>-- Status --</Select.Option>
          <Select.Option value={2}>Active</Select.Option>
          <Select.Option value={1}>Inactive</Select.Option>
        </Select>
      </div>

      <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={clients}
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
