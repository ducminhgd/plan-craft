import { useState, useEffect } from 'react';
import { Table, Button, Space, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { GetClients, DeleteClient } from '../../../wailsjs/go/main/App';
import { entities } from '../../../wailsjs/go/models';

const { Title } = Typography;

export default function ClientList() {
  const navigate = useNavigate();
  const [clients, setClients] = useState<entities.Client[]>([]);
  const [loading, setLoading] = useState(false);

  const loadClients = async () => {
    setLoading(true);
    try {
      const data = await GetClients({} as entities.ClientQueryParams);
      setClients(data || []);
    } catch (error) {
      message.error('Failed to load clients');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClients();
  }, []);

  const handleDelete = async (id: number) => {
    try {
      await DeleteClient(id);
      message.success('Client deleted');
      loadClients();
    } catch (error) {
      message.error('Failed to delete client');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'ID', key: 'ID' },
    { title: 'Name', dataIndex: 'Name', key: 'Name' },
    { title: 'Description', dataIndex: 'Description', key: 'Description' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: any) => (
        <Space>
          <Button
            icon={<EditOutlined />}
            onClick={() => navigate(`/clients/${record.ID}`)}
          />
          <Button
            icon={<DeleteOutlined />}
            danger
            onClick={() => handleDelete(record.ID)}
          />
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={2}>Clients</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/clients/new')}
        >
          Add New
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={clients}
        loading={loading}
        rowKey="ID"
      />
    </div>
  );
}
