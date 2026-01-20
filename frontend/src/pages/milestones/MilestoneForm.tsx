import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select, DatePicker } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetMilestone, CreateMilestone, UpdateMilestone, GetProjects } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import { DATE_FORMAT, parseDate, toISOString } from '../../utils/date';

const { Title } = Typography;
const { TextArea } = Input;

export default function MilestoneForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [projects, setProjects] = useState<entities.Project[]>([]);

  const isEdit = !!id;

  useEffect(() => {
    loadProjects();
  }, []);

  useEffect(() => {
    if (isEdit) {
      loadMilestone();
    } else {
      // Reset form and set defaults for new milestones
      form.resetFields();
      form.setFieldsValue({
        status: 2,
      });
    }
  }, [id]);

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

  const loadMilestone = async () => {
    try {
      const milestone = await GetMilestone(Number(id));
      // Convert date strings to dayjs objects for DatePicker
      const formValues: any = {
        ...milestone,
        start_date: parseDate(milestone.start_date),
        end_date: parseDate(milestone.end_date),
      };
      form.setFieldsValue(formValues);
    } catch (error) {
      message.error('Failed to load milestone');
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      // Convert dayjs objects back to ISO strings for the backend
      const milestoneData: any = {
        ...values,
        start_date: toISOString(values.start_date),
        end_date: toISOString(values.end_date),
      };

      if (isEdit) {
        await UpdateMilestone({ ...milestoneData, id: Number(id) });
        message.success('Milestone updated');
      } else {
        await CreateMilestone(milestoneData);
        message.success('Milestone created');
      }
      navigate('/milestones');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} milestone`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Milestone' : 'Add New Milestone'}</Title>
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
            { required: true, message: 'Please input milestone name' },
            { max: 255, message: 'Name cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter milestone name" />
        </Form.Item>

        <Form.Item
          label="Project"
          name="project_id"
          rules={[{ required: true, message: 'Please select a project' }]}
        >
          <Select
            placeholder="Select a project"
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {projects.map(project => (
              <Select.Option key={project.id} value={project.id}>{project.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Description"
          name="description"
        >
          <TextArea rows={4} placeholder="Enter milestone description" />
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
            <Button onClick={() => navigate('/milestones')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
}
