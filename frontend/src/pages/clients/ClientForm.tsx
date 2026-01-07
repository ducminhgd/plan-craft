import { useState, useEffect } from 'react';
import { Form, Input, Button, Typography, message, Space } from 'antd';
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
        await UpdateClient({ ...values, ID: Number(id) });
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
          name="Name"
          rules={[{ required: true, message: 'Please input client name' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Description"
          name="Description"
        >
          <TextArea rows={4} />
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
