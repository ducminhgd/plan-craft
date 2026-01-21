import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select, DatePicker, InputNumber, Checkbox, Divider, Table, Modal } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { GetProject, CreateProject, UpdateProject, GetClients, GetProjectRolesByProject, CreateProjectRole, UpdateProjectRole, DeleteProjectRole } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import { DATE_FORMAT, parseDate, toISOString } from '../../utils/date';

const { Title } = Typography;
const { TextArea } = Input;

const WEEKDAYS = [
  { value: 1, label: 'Monday' },
  { value: 2, label: 'Tuesday' },
  { value: 3, label: 'Wednesday' },
  { value: 4, label: 'Thursday' },
  { value: 5, label: 'Friday' },
  { value: 6, label: 'Saturday' },
  { value: 0, label: 'Sunday' },
];

const COMMON_CURRENCIES = [
  'USD', 'EUR', 'GBP', 'JPY', 'CNY', 'AUD', 'CAD', 'CHF', 'HKD', 'SGD', 'VND',
];

const COMMON_TIMEZONES = [
  'UTC',
  'America/New_York',
  'America/Los_Angeles',
  'America/Chicago',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Asia/Singapore',
  'Asia/Ho_Chi_Minh',
  'Australia/Sydney',
];

const ROLE_LEVELS = [
  { value: 1, label: 'Junior' },
  { value: 2, label: 'Mid' },
  { value: 3, label: 'Senior' },
  { value: 4, label: 'Lead' },
  { value: 5, label: 'Manager' },
  { value: 6, label: 'Director' },
  { value: 7, label: 'VP' },
  { value: 8, label: 'C-Level' },
];

export default function ProjectForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [roleForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [clients, setClients] = useState<entities.Client[]>([]);
  const [projectRoles, setProjectRoles] = useState<entities.ProjectRole[]>([]);
  const [roleModalVisible, setRoleModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<entities.ProjectRole | null>(null);
  const [roleLoading, setRoleLoading] = useState(false);

  const isEdit = !!id;

  useEffect(() => {
    loadClients();
  }, []);

  useEffect(() => {
    if (isEdit) {
      loadProject();
      loadProjectRoles();
    } else {
      // Reset form and set defaults for new projects
      form.resetFields();
      form.setFieldsValue({
        status: 2,
        hours_per_day: 8,
        days_per_week: 5,
        working_days_per_week: [1, 2, 3, 4, 5], // Monday to Friday
      });
      setProjectRoles([]);
    }
  }, [id]);

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

  const loadProject = async () => {
    try {
      const project = await GetProject(Number(id));
      // Convert date strings to dayjs objects for DatePicker
      const formValues: any = {
        ...project,
        start_date: parseDate(project.start_date),
        end_date: parseDate(project.end_date),
      };
      form.setFieldsValue(formValues);
    } catch (error) {
      message.error('Failed to load project');
    }
  };

  const loadProjectRoles = async () => {
    if (!id) return;
    try {
      const result = await GetProjectRolesByProject(Number(id));
      setProjectRoles(result.data || []);
    } catch (error) {
      console.error('Failed to load project roles', error);
    }
  };

  const handleAddRole = () => {
    setEditingRole(null);
    roleForm.resetFields();
    roleForm.setFieldsValue({ level: 2, headcount: 1 });
    setRoleModalVisible(true);
  };

  const handleEditRole = (role: entities.ProjectRole) => {
    setEditingRole(role);
    roleForm.setFieldsValue({
      name: role.name,
      level: role.level,
      headcount: role.headcount,
    });
    setRoleModalVisible(true);
  };

  const handleDeleteRole = (roleId: number) => {
    Modal.confirm({
      title: 'Delete Role',
      content: 'Are you sure you want to delete this role?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await DeleteProjectRole(roleId);
          message.success('Role deleted');
          loadProjectRoles();
        } catch (error) {
          message.error('Failed to delete role');
        }
      },
    });
  };

  const handleRoleSubmit = async (values: any) => {
    if (!id) return;
    setRoleLoading(true);
    try {
      const roleData: any = {
        ...values,
        project_id: Number(id),
      };

      if (editingRole) {
        await UpdateProjectRole({ ...roleData, id: editingRole.id });
        message.success('Role updated');
      } else {
        await CreateProjectRole(roleData);
        message.success('Role created');
      }
      setRoleModalVisible(false);
      loadProjectRoles();
    } catch (error) {
      message.error(`Failed to ${editingRole ? 'update' : 'create'} role`);
    } finally {
      setRoleLoading(false);
    }
  };

  const roleColumns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Level',
      dataIndex: 'level',
      key: 'level',
      render: (level: number) => ROLE_LEVELS.find(l => l.value === level)?.label || 'Unknown',
    },
    {
      title: 'Headcount',
      dataIndex: 'headcount',
      key: 'headcount',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: entities.ProjectRole) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEditRole(record)}
          />
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteRole(record.id)}
          />
        </Space>
      ),
    },
  ];

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      // Convert dayjs objects back to ISO strings for the backend
      const projectData: any = {
        ...values,
        start_date: toISOString(values.start_date),
        end_date: toISOString(values.end_date),
      };

      if (isEdit) {
        await UpdateProject({ ...projectData, id: Number(id) });
        message.success('Project updated');
      } else {
        await CreateProject(projectData);
        message.success('Project created');
      }
      navigate('/projects');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} project`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Project' : 'Add New Project'}</Title>
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        style={{ maxWidth: 600 }}
      >
        <Form.Item
          label="Name"
          name="name"
          rules={[
            { required: true, message: 'Please input project name' },
            { max: 255, message: 'Name cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter project name" />
        </Form.Item>

        <Form.Item
          label="Client"
          name="client_id"
          rules={[{ required: true, message: 'Please select a client' }]}
        >
          <Select
            placeholder="Select a client"
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {clients.map(client => (
              <Select.Option key={client.id} value={client.id}>{client.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Description"
          name="description"
        >
          <TextArea rows={4} placeholder="Enter project description" />
        </Form.Item>

        <Form.Item
          label="Start Date"
          name="start_date"
        >
          <DatePicker style={{ width: '100%' }} placeholder="Select start date" format={DATE_FORMAT} />
        </Form.Item>

        <Form.Item
          label="End Date"
          name="end_date"
          dependencies={['start_date']}
          rules={[
            ({ getFieldValue }) => ({
              validator(_, value) {
                const startDate = getFieldValue('start_date');
                if (!value || !startDate) {
                  return Promise.resolve();
                }
                if (value.isSame(startDate) || value.isAfter(startDate)) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('End date must be greater than or equal to start date'));
              },
            }),
          ]}
        >
          <DatePicker style={{ width: '100%' }} placeholder="Select end date" format={DATE_FORMAT} />
        </Form.Item>

        <Divider>Configuration</Divider>

        <Form.Item
          label="Hours per Day"
          name="hours_per_day"
          rules={[
            { type: 'number', min: 1, max: 24, message: 'Hours per day must be between 1 and 24' }
          ]}
        >
          <InputNumber
            style={{ width: '100%' }}
            placeholder="Enter hours per day (default: 8)"
            min={1}
            max={24}
            precision={0}
          />
        </Form.Item>

        <Form.Item
          label="Days per Week"
          name="days_per_week"
          rules={[
            { type: 'number', min: 1, max: 7, message: 'Days per week must be between 1 and 7' }
          ]}
        >
          <InputNumber
            style={{ width: '100%' }}
            placeholder="Enter days per week (default: 5)"
            min={1}
            max={7}
            precision={0}
          />
        </Form.Item>

        <Form.Item
          label="Working Days"
          name="working_days_per_week"
        >
          <Checkbox.Group options={WEEKDAYS} />
        </Form.Item>

        <Form.Item
          label="Timezone"
          name="timezone"
        >
          <Select
            placeholder="Select timezone"
            allowClear
            showSearch
          >
            {COMMON_TIMEZONES.map(tz => (
              <Select.Option key={tz} value={tz}>{tz}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Currency"
          name="currency"
        >
          <Select
            placeholder="Select currency"
            allowClear
            showSearch
          >
            {COMMON_CURRENCIES.map(curr => (
              <Select.Option key={curr} value={curr}>{curr}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Divider />

        <Form.Item
          label="Status"
          name="status"
          rules={[{ required: true, message: 'Please select status' }]}
        >
          <Select>
            <Select.Option value={2}>Active</Select.Option>
            <Select.Option value={1}>Inactive</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              {isEdit ? 'Update' : 'Create'}
            </Button>
            <Button onClick={() => navigate('/projects')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>

      {isEdit && (
        <>
          <Divider>Project Roles</Divider>
          <div style={{ marginBottom: 16 }}>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAddRole}>
              Add Role
            </Button>
          </div>
          <Table
            columns={roleColumns}
            dataSource={projectRoles}
            rowKey="id"
            pagination={false}
            size="small"
            style={{ maxWidth: 600 }}
          />

          <Modal
            title={editingRole ? 'Edit Role' : 'Add Role'}
            open={roleModalVisible}
            onCancel={() => setRoleModalVisible(false)}
            footer={null}
          >
            <Form
              form={roleForm}
              layout="vertical"
              onFinish={handleRoleSubmit}
            >
              <Form.Item
                label="Role Name"
                name="name"
                rules={[{ required: true, message: 'Please input role name' }]}
              >
                <Input placeholder="e.g., Backend Developer, UI Designer" />
              </Form.Item>

              <Form.Item
                label="Level"
                name="level"
                rules={[{ required: true, message: 'Please select level' }]}
              >
                <Select placeholder="Select level">
                  {ROLE_LEVELS.map(level => (
                    <Select.Option key={level.value} value={level.value}>
                      {level.label}
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>

              <Form.Item
                label="Headcount"
                name="headcount"
                rules={[
                  { required: true, message: 'Please input headcount' },
                  { type: 'number', min: 0, message: 'Headcount must be non-negative' }
                ]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  placeholder="Number of people needed"
                  min={0}
                  precision={0}
                />
              </Form.Item>

              <Form.Item>
                <Space>
                  <Button type="primary" htmlType="submit" loading={roleLoading}>
                    {editingRole ? 'Update' : 'Create'}
                  </Button>
                  <Button onClick={() => setRoleModalVisible(false)}>
                    Cancel
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Modal>
        </>
      )}
    </div>
  );
}
