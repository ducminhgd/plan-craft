import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select, InputNumber } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetTask, CreateTask, UpdateTask, GetProjects, GetMilestones, GetTasks } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';

const { Title } = Typography;
const { TextArea } = Input;

export default function TaskForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [projects, setProjects] = useState<entities.Project[]>([]);
  const [milestones, setMilestones] = useState<entities.Milestone[]>([]);
  const [parentTasks, setParentTasks] = useState<entities.Task[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>();

  const isEdit = !!id;

  useEffect(() => {
    loadProjects();
  }, []);

  useEffect(() => {
    if (isEdit) {
      loadTask();
    } else {
      form.resetFields();
      form.setFieldsValue({
        status: 1, // To Do
        priority: 2, // Medium
        level: 1,
        estimated_effort: 0,
      });
    }
  }, [id]);

  useEffect(() => {
    if (selectedProjectId) {
      loadMilestones(selectedProjectId);
      loadParentTasks(selectedProjectId);
    } else {
      setMilestones([]);
      setParentTasks([]);
    }
  }, [selectedProjectId]);

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

  const loadMilestones = async (projectId: number) => {
    try {
      const params: any = {
        status: 2, // Only active milestones
        project_id: projectId,
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      const result = await GetMilestones(params);
      setMilestones(result.data || []);
    } catch (error) {
      console.error('Failed to load milestones', error);
    }
  };

  const loadParentTasks = async (projectId: number) => {
    try {
      const params: any = {
        project_id: projectId,
        pagination: {
          page: 1,
          page_size: 1000,
          total: 0,
        },
      };
      const result = await GetTasks(params);
      // Filter out current task if editing
      const tasks = (result.data || []).filter((t: entities.Task) => !isEdit || t.id !== Number(id));
      setParentTasks(tasks);
    } catch (error) {
      console.error('Failed to load parent tasks', error);
    }
  };

  const loadTask = async () => {
    try {
      const task = await GetTask(Number(id));
      form.setFieldsValue({
        ...task,
        milestone_id: task.milestone_id || undefined,
        parent_id: task.parent_id || undefined,
      });
      setSelectedProjectId(task.project_id);
    } catch (error) {
      message.error('Failed to load task');
    }
  };

  const handleProjectChange = (value: number) => {
    setSelectedProjectId(value);
    // Reset milestone and parent when project changes
    form.setFieldsValue({
      milestone_id: undefined,
      parent_id: undefined,
    });
  };

  const handleParentChange = (value: number | undefined) => {
    if (value) {
      const parent = parentTasks.find(t => t.id === value);
      if (parent) {
        // Set level to parent's level + 1
        form.setFieldsValue({ level: parent.level + 1 });
      }
    } else {
      // Reset to level 1 if no parent
      form.setFieldsValue({ level: 1 });
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const taskData: any = {
        ...values,
        milestone_id: values.milestone_id || null,
        parent_id: values.parent_id || null,
      };

      if (isEdit) {
        await UpdateTask({ ...taskData, id: Number(id) });
        message.success('Task updated');
      } else {
        await CreateTask(taskData);
        message.success('Task created');
      }
      navigate('/tasks');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} task`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Task' : 'Add New Task'}</Title>
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
            { required: true, message: 'Please input task name' },
            { max: 255, message: 'Name cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter task name" />
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
            onChange={handleProjectChange}
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
          label="Milestone"
          name="milestone_id"
        >
          <Select
            placeholder="Select a milestone (optional)"
            showSearch
            allowClear
            optionFilterProp="children"
            disabled={!selectedProjectId}
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {milestones.map(milestone => (
              <Select.Option key={milestone.id} value={milestone.id}>{milestone.name}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Parent Task"
          name="parent_id"
        >
          <Select
            placeholder="Select a parent task (optional)"
            showSearch
            allowClear
            optionFilterProp="children"
            disabled={!selectedProjectId}
            onChange={handleParentChange}
            filterOption={(input, option) =>
              (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
            }
          >
            {parentTasks.map(task => (
              <Select.Option key={task.id} value={task.id}>
                {'—'.repeat(task.level - 1)}{task.level > 1 ? ' ' : ''}{task.name}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Description"
          name="description"
        >
          <TextArea rows={4} placeholder="Enter task description" />
        </Form.Item>

        <Form.Item
          label="Level"
          name="level"
          rules={[{ required: true, message: 'Please input task level' }]}
          tooltip="1 = Epic, 2 = Task, 3 = Subtask, etc. Auto-set based on parent task."
        >
          <InputNumber min={1} style={{ width: '100%' }} disabled />
        </Form.Item>

        <Form.Item
          label="Priority"
          name="priority"
          rules={[{ required: true, message: 'Please select priority' }]}
        >
          <Select>
            <Select.Option value={1}>Low</Select.Option>
            <Select.Option value={2}>Medium</Select.Option>
            <Select.Option value={3}>High</Select.Option>
            <Select.Option value={4}>Critical</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Status"
          name="status"
          rules={[{ required: true, message: 'Please select status' }]}
        >
          <Select>
            <Select.Option value={1}>To Do</Select.Option>
            <Select.Option value={2}>In Progress</Select.Option>
            <Select.Option value={3}>Done</Select.Option>
            <Select.Option value={4}>Cancelled</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Estimated Effort (hours)"
          name="estimated_effort"
          rules={[
            { type: 'number', min: 0, message: 'Effort must be non-negative' }
          ]}
        >
          <InputNumber min={0} step={0.5} style={{ width: '100%' }} placeholder="0" />
        </Form.Item>

        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              {isEdit ? 'Update' : 'Create'}
            </Button>
            <Button onClick={() => navigate('/tasks')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
}
