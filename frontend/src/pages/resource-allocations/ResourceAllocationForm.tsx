import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select, DatePicker, InputNumber } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetProjectResource, CreateProjectResource, UpdateProjectResource, GetProjects, GetHumanResources } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';
import { DATE_FORMAT, parseDate, toISOString } from '../../utils/date';

const { Title } = Typography;
const { TextArea } = Input;

export default function ResourceAllocationForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [projects, setProjects] = useState<entities.Project[]>([]);
  const [humanResources, setHumanResources] = useState<entities.HumanResource[]>([]);
  const [selectedProject, setSelectedProject] = useState<entities.Project | null>(null);

  const isEdit = !!id;

  useEffect(() => {
    loadProjects();
    loadHumanResources();
  }, []);

  useEffect(() => {
    if (isEdit) {
      loadAllocation();
    } else {
      // Reset form and set defaults for new allocations
      form.resetFields();
      form.setFieldsValue({ status: 2, allocation: 100, cost: 0 });
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

  const loadAllocation = async () => {
    try {
      const allocation = await GetProjectResource(Number(id));
      // Convert date strings to dayjs objects for DatePicker
      const formValues: any = {
        ...allocation,
        start_date: parseDate(allocation.start_date),
        end_date: parseDate(allocation.end_date),
      };
      form.setFieldsValue(formValues);
      // Set the selected project based on the loaded allocation
      const project = projects.find(p => p.id === allocation.project_id);
      if (project) {
        setSelectedProject(project);
      }
    } catch (error) {
      message.error('Failed to load allocation');
    }
  };

  // Update selected project when project_id changes or projects load
  useEffect(() => {
    const projectId = form.getFieldValue('project_id');
    if (projectId && projects.length > 0) {
      const project = projects.find(p => p.id === projectId);
      if (project) {
        setSelectedProject(project);
      }
    }
  }, [projects]);

  const handleProjectChange = (projectId: number) => {
    const project = projects.find(p => p.id === projectId);
    setSelectedProject(project || null);
    // Set start and end dates from project if available
    const projectStartDate = project?.start_date ? parseDate(project.start_date) : null;
    const projectEndDate = project?.end_date ? parseDate(project.end_date) : null;
    form.setFieldsValue({ start_date: projectStartDate, end_date: projectEndDate });
  };

  // Helper to get project date constraints for validation messages
  const getProjectDateRangeText = () => {
    if (!selectedProject) return '';
    const parts = [];
    if (selectedProject.start_date) {
      parts.push(`from ${parseDate(selectedProject.start_date)?.format(DATE_FORMAT)}`);
    }
    if (selectedProject.end_date) {
      parts.push(`to ${parseDate(selectedProject.end_date)?.format(DATE_FORMAT)}`);
    }
    return parts.length > 0 ? `Project runs ${parts.join(' ')}` : '';
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      // Convert dayjs objects back to ISO strings for the backend
      const allocationData: any = {
        ...values,
        start_date: toISOString(values.start_date),
        end_date: toISOString(values.end_date),
      };

      if (isEdit) {
        await UpdateProjectResource({ ...allocationData, id: Number(id) });
        message.success('Allocation updated');
      } else {
        await CreateProjectResource(allocationData);
        message.success('Allocation created');
      }
      navigate('/resource-allocations');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} allocation`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Project Allocation' : 'Add New Project Allocation'}</Title>
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        style={{ maxWidth: 600 }}
      >
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
            onChange={handleProjectChange}
          >
            {projects.map(project => (
              <Select.Option key={project.id} value={project.id}>{project.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Human Resource"
          name="human_resource_id"
          rules={[{ required: true, message: 'Please select a human resource' }]}
        >
          <Select
            placeholder="Select a human resource"
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {humanResources.map(hr => (
              <Select.Option key={hr.id} value={hr.id}>{hr.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Role"
          name="role"
          rules={[
            { max: 255, message: 'Role cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter role in project (e.g., Developer, Tech Lead, QA)" />
        </Form.Item>

        <Form.Item
          label="Allocation (%)"
          name="allocation"
          rules={[
            { required: true, message: 'Please input allocation percentage' },
            { type: 'number', min: 0, max: 100, message: 'Allocation must be between 0 and 100' }
          ]}
        >
          <InputNumber
            style={{ width: '100%' }}
            placeholder="Enter allocation percentage"
            min={0}
            max={100}
            precision={0}
          />
        </Form.Item>

        <Form.Item
          label={selectedProject?.currency ? `Cost (${selectedProject.currency})` : 'Cost'}
          name="cost"
          style={{ width: '100%' }}
          rules={[
            { type: 'number', min: 0, message: 'Cost must be a positive number' }
          ]}
        >
          <InputNumber
            style={{ width: '100%' }}
            placeholder="Enter cost for this resource"
            min={0}
            precision={2}
            addonAfter={selectedProject?.currency || undefined}
          />
        </Form.Item>

        <Form.Item
          label="Start Date"
          name="start_date"
          dependencies={['project_id']}
          extra={getProjectDateRangeText()}
          rules={[
            () => ({
              validator(_, value) {
                if (!value || !selectedProject) {
                  return Promise.resolve();
                }
                const projectStart = parseDate(selectedProject.start_date);
                const projectEnd = parseDate(selectedProject.end_date);

                if (projectStart && value.isBefore(projectStart, 'day')) {
                  return Promise.reject(new Error(`Start date cannot be before project start date (${projectStart.format(DATE_FORMAT)})`));
                }
                if (projectEnd && value.isAfter(projectEnd, 'day')) {
                  return Promise.reject(new Error(`Start date cannot be after project end date (${projectEnd.format(DATE_FORMAT)})`));
                }
                return Promise.resolve();
              },
            }),
          ]}
        >
          <DatePicker style={{ width: '100%' }} placeholder="Select start date" format={DATE_FORMAT} />
        </Form.Item>

        <Form.Item
          label="End Date"
          name="end_date"
          dependencies={['start_date', 'project_id']}
          rules={[
            ({ getFieldValue }) => ({
              validator(_, value) {
                const startDate = getFieldValue('start_date');

                // Validate against resource start date
                if (value && startDate) {
                  if (!value.isSame(startDate, 'day') && !value.isAfter(startDate, 'day')) {
                    return Promise.reject(new Error('End date must be greater than or equal to start date'));
                  }
                }

                // Validate against project date range
                if (value && selectedProject) {
                  const projectStart = parseDate(selectedProject.start_date);
                  const projectEnd = parseDate(selectedProject.end_date);

                  if (projectStart && value.isBefore(projectStart, 'day')) {
                    return Promise.reject(new Error(`End date cannot be before project start date (${projectStart.format(DATE_FORMAT)})`));
                  }
                  if (projectEnd && value.isAfter(projectEnd, 'day')) {
                    return Promise.reject(new Error(`End date cannot be after project end date (${projectEnd.format(DATE_FORMAT)})`));
                  }
                }

                return Promise.resolve();
              },
            }),
          ]}
        >
          <DatePicker style={{ width: '100%' }} placeholder="Select end date" format={DATE_FORMAT} />
        </Form.Item>

        <Form.Item
          label="Notes"
          name="notes"
        >
          <TextArea rows={4} placeholder="Enter additional notes" />
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
            <Button onClick={() => navigate('/resource-allocations')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
}
