import { RouterProvider } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import { router } from './router';
import './App.css';

function App() {
    return (
        <ConfigProvider theme={{ token: { colorPrimary: '#1890ff' } }}>
            <RouterProvider router={router} />
        </ConfigProvider>
    );
}

export default App
