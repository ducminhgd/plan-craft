import { RouterProvider } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import { router } from './router';
import { DatabaseProvider } from './contexts/DatabaseContext';
import './App.css';

function App() {
    return (
        <ConfigProvider theme={{ token: { colorPrimary: '#1890ff' } }}>
            <DatabaseProvider>
                <RouterProvider router={router} />
            </DatabaseProvider>
        </ConfigProvider>
    );
}

export default App
