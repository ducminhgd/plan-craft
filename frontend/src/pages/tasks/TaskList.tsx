import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message, Input, Select, Modal, Tag } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetTasks, DeleteTask, GetProjects, GetMilestones } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import type { TableRowSelection } from 'antd/es/table/interface';

const { Title } = Typography;
const { Search } = Input;

export default function TaskList() {
  const navigate = useNavigate();
  const [tasks, setTasks] = useState<entities.Task[]>([]);
  const [projects, setProjects] = useState<entities.Project[]>([]);
  const [milestones, setMilestones] = useState<entities.Milestone[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Filter states
  const [searchText, setSearchText] = useState('');
  const [statusFilter, setStatusFilter] = useState<number | undefined>(0);
  const [priorityFilter, setPriorityFilter] = useState<number | undefined>(0);
  const [projectFilter, setProjectFilter] = useState<number | undefined>(0);
  const [milestoneFilter, setMilestoneFilter] = useState<number | undefined>(0);

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

  const loadMilestones = async () => {
    try {
      const params: any = {
        status: 2, // Only active milestones
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      // Apply project filter if set
      if (projectFilter !== 0 && projectFilter !== undefined) {
        params.project_id = projectFilter;
      }
      const result = await GetMilestones(params);
      setMilestones(result.data || []);
    } catch (error) {
      console.error('Failed to load milestones', error);
    }
  };

  const loadTasks = async () => {
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

      // Apply priority filter
      if (priorityFilter !== 0 && priorityFilter !== undefined) {
        params.priority = priorityFilter;
      }

      // Apply project filter
      if (projectFilter !== 0 && projectFilter !== undefined) {
        params.project_id = projectFilter;
      }

      // Apply milestone filter
      if (milestoneFilter !== 0 && milestoneFilter !== undefined) {
        params.milestone_id = milestoneFilter;
      }

      const result = await GetTasks(params);
      const newTotal = result.total || 0;
      setTasks(result.data || []);
      setTotal(newTotal);

      // Clamp currentPage if it exceeds total pages after filtering
      const totalPages = Math.ceil(newTotal / pageSize);
      if (totalPages > 0 && currentPage > totalPages) {
        setCurrentPage(totalPages);
      }
    } catch (error) {
      message.error('Failed to load tasks');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProjects();
  }, []);

  useEffect(() => {
    loadMilestones();
  }, [projectFilter]);

  useEffect(() => {
    loadTasks();
  }, [currentPage, pageSize, searchText, statusFilter, priorityFilter, projectFilter, milestoneFilter]);

  const handleDelete = async (task: entities.Task) => {
    Modal.confirm({
      title: 'Confirm Delete',
      content: `Are you sure you want to delete task "${task.name}"?`,
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await DeleteTask(task.id);
          message.success('Task deleted successfully');
          loadTasks();
        } catch (error) {
          message.error('Failed to delete task');
        }
      },
    });
  };

  const handleSearchChange = (value: string) => {
    setSearchText(value);
    setCurrentPage(1);
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

  const handlePriorityChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setPriorityFilter(value);
    setCurrentPage(1);
  };

  const handleProjectChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setProjectFilter(value);
    setMilestoneFilter(0); // Reset milestone filter when project changes
    setCurrentPage(1);
  };

  const handleMilestoneChange = (value: number | undefined) => {
    if (value == undefined) {
      value = 0;
    }
    setMilestoneFilter(value);
    setCurrentPage(1);
  };

  const handlePageSizeChange = (value: number) => {
    setPageSize(value);
    setCurrentPage(1);
  };

  const onSelectChange = (newSelectedRowKeys: React.Key[]) => {
    setSelectedRowKeys(newSelectedRowKeys);
  };

  const rowSelection: TableRowSelection<entities.Task> = {
    selectedRowKeys,
    onChange: onSelectChange,
  };

  const getStatusTag = (status: number) => {
    switch (status) {
      case 1:
        return <Tag color="default">To Do</Tag>;
      case 2:
        return <Tag color="processing">In Progress</Tag>;
      case 3:
        return <Tag color="success">Done</Tag>;
      case 4:
        return <Tag color="error">Cancelled</Tag>;
      default:
        return <Tag>Unknown</Tag>;
    }
  };

  const getPriorityTag = (priority: number) => {
    switch (priority) {
      case 1:
        return <Tag color="green">Low</Tag>;
      case 2:
        return <Tag color="blue">Medium</Tag>;
      case 3:
        return <Tag color="orange">High</Tag>;
      case 4:
        return <Tag color="red">Critical</Tag>;
      default:
        return <Tag>Unknown</Tag>;
    }
  };

  const getProjectName = (projectId: number) => {
    const project = projects.find(p => p.id === projectId);
    return project ? project.name : '-';
  };

  const getMilestoneName = (milestoneId: number | undefined) => {
    if (!milestoneId) return '-';
    const milestone = milestones.find(m => m.id === milestoneId);
    return milestone ? milestone.name : '-';
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: entities.Task) => (
        <a onClick={() => navigate(`/tasks/${record.id}`)} style={{ cursor: 'pointer' }}>
          {record.level > 1 && <span style={{ marginRight: 4 }}>{'—'.repeat(record.level - 1)}</span>}
          {name}
        </a>
      ),
    },
    {
      title: 'Project',
      dataIndex: 'project_id',
      key: 'project_id',
      render: (projectId: number) => getProjectName(projectId),
    },
    {
      title: 'Milestone',
      dataIndex: 'milestone_id',
      key: 'milestone_id',
      render: (milestoneId: number | undefined) => getMilestoneName(milestoneId),
    },
    {
      title: 'Priority',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority: number) => getPriorityTag(priority),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => getStatusTag(status),
    },
    {
      title: 'Effort (h)',
      dataIndex: 'estimated_effort',
      key: 'estimated_effort',
      render: (effort: number) => effort > 0 ? effort.toFixed(1) : '-',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: entities.Task) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/tasks/${record.id}`)}
            size="small"
          />
          <Button
            icon={<DeleteOutlined />}
            danger
            onClick={() => handleDelete(record)}
            size="small"
            title="Delete task"
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
        <Title level={2}>List of Tasks</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/tasks/new')}
        >
          Add New
        </Button>
      </div>

      {/* Filter Section */}
      <div style={{ marginBottom: 16, display: 'flex', gap: 16, flexWrap: 'wrap' }}>
        <Search
          placeholder="Search by name or description"
          onSearch={handleSearch}
          onChange={(e) => handleSearchChange(e.target.value)}
          style={{ flex: 1, minWidth: 200 }}
          allowClear
        />
        <Select
          style={{ width: 180 }}
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
          style={{ width: 180 }}
          onChange={handleMilestoneChange}
          value={milestoneFilter}
          showSearch
          optionFilterProp="children"
          filterOption={(input, option) =>
            (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
          }
        >
          <Select.Option value={0}>-- All Milestones --</Select.Option>
          {milestones.map(milestone => (
            <Select.Option key={milestone.id} value={milestone.id}>{milestone.name}</Select.Option>
          ))}
        </Select>
        <Select
          style={{ width: 130 }}
          onChange={handlePriorityChange}
          value={priorityFilter}
        >
          <Select.Option value={0}>-- Priority --</Select.Option>
          <Select.Option value={1}>Low</Select.Option>
          <Select.Option value={2}>Medium</Select.Option>
          <Select.Option value={3}>High</Select.Option>
          <Select.Option value={4}>Critical</Select.Option>
        </Select>
        <Select
          style={{ width: 140 }}
          onChange={handleStatusChange}
          value={statusFilter}
        >
          <Select.Option value={0}>-- Status --</Select.Option>
          <Select.Option value={1}>To Do</Select.Option>
          <Select.Option value={2}>In Progress</Select.Option>
          <Select.Option value={3}>Done</Select.Option>
          <Select.Option value={4}>Cancelled</Select.Option>
        </Select>
      </div>

      <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={tasks}
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
