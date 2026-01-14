import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetHumanResource, CreateHumanResource, UpdateHumanResource } from '../../../wailsjs/go/main/App';

const { Title } = Typography;

export default function HumanResourceForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const isEdit = !!id;

  useEffect(() => {
    if (isEdit) {
      loadHumanResource();
    } else {
      // Set default status to Active (2) for new human resources
      form.setFieldsValue({ status: 2 });
    }
  }, [id]);

  const loadHumanResource = async () => {
    try {
      const hr = await GetHumanResource(Number(id));
      form.setFieldsValue(hr);
    } catch (error) {
      message.error('Failed to load human resource');
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      if (isEdit) {
        await UpdateHumanResource({ ...values, id: Number(id) });
        message.success('Human resource updated');
      } else {
        await CreateHumanResource(values);
        message.success('Human resource created');
      }
      navigate('/human-resources');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} human resource`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Human Resource' : 'Add New Human Resource'}</Title>
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
            { required: true, message: 'Please input name' },
            { max: 255, message: 'Name cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter name" />
        </Form.Item>

        <Form.Item
          label="Title"
          name="title"
          rules={[
            { required: true, message: 'Please input title' },
            { max: 255, message: 'Title cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter job title (e.g., Software Engineer, Project Manager)" />
        </Form.Item>

        <Form.Item
          label="Level"
          name="level"
          rules={[
            { required: true, message: 'Please input level' },
            { max: 255, message: 'Level cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter level (e.g., Junior, Senior, Lead)" />
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
            <Button onClick={() => navigate('/human-resources')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
}
