import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space, Select } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { GetClient, CreateClient, UpdateClient } from '../../../wailsjs/go/main/App';

const { Title } = Typography;
const { TextArea } = Input;

export default function ClientForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const isEdit = !!id;

  useEffect(() => {
    if (isEdit) {
      loadClient();
    } else {
      // Set default status to Active (2) for new clients
      form.setFieldsValue({ status: 2 });
    }
  }, [id]);

  const loadClient = async () => {
    try {
      const client = await GetClient(Number(id));
      form.setFieldsValue(client);
    } catch (error) {
      message.error('Failed to load client');
    }
  };

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      if (isEdit) {
        await UpdateClient({ ...values, id: Number(id) });
        message.success('Client updated');
      } else {
        await CreateClient(values);
        message.success('Client created');
      }
      navigate('/clients');
    } catch (error) {
      message.error(`Failed to ${isEdit ? 'update' : 'create'} client`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Title level={2}>{isEdit ? 'Edit Client' : 'Add New Client'}</Title>
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
            { required: true, message: 'Please input client name' },
            { max: 255, message: 'Name cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter client name" />
        </Form.Item>

        <Form.Item
          label="Email"
          name="email"
          rules={[
            { required: true, message: 'Please input client email' },
            { type: 'email', message: 'Please enter a valid email address' },
            { max: 255, message: 'Email cannot exceed 255 characters' }
          ]}
        >
          <Input placeholder="Enter email address" />
        </Form.Item>

        <Form.Item
          label="Phone"
          name="phone"
          rules={[
            { max: 50, message: 'Phone cannot exceed 50 characters' }
          ]}
        >
          <Input placeholder="Enter phone number" />
        </Form.Item>

        <Form.Item
          label="Address"
          name="address"
        >
          <TextArea rows={3} placeholder="Enter address" />
        </Form.Item>

        <Form.Item
          label="Contact Person"
          name="contact_person"
        >
          <Input placeholder="Enter contact person name" />
        </Form.Item>

        <Form.Item
          label="Notes"
          name="notes"
        >
          <TextArea rows={4} placeholder="Enter any additional notes" />
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
            <Button onClick={() => navigate('/clients')}>
              Cancel
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </div>
  );
}
