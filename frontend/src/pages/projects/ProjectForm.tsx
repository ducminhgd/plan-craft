import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select, DatePicker } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetProject, CreateProject, UpdateProject, GetClients } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import { DATE_FORMAT, parseDate, toISOString } from '../../utils/date';

const { Title } = Typography;
const { TextArea } = Input;

export default function ProjectForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [clients, setClients] = useState<entities.Client[]>([]);

  const isEdit = !!id;

  useEffect(() => {
    loadClients();
  }, []);

  useEffect(() => {
    if (isEdit) {
      loadProject();
    } else {
      // Reset form and set default status to Active (2) for new projects
      form.resetFields();
      form.setFieldsValue({ status: 2 });
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
    </div>
  );
}
