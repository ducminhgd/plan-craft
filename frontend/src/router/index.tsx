import { createBrowserRouter } from 'react-router-dom';
import MainLayout from '../layouts/MainLayout';
import Dashboard from '../pages/Dashboard';
import ClientList from '../pages/clients/ClientList';
import ClientForm from '../pages/clients/ClientForm';
import ProjectList from '../pages/projects/ProjectList';
import ProjectForm from '../pages/projects/ProjectForm';
import HumanResourceList from '../pages/human-resources/HumanResourceList';
import HumanResourceForm from '../pages/human-resources/HumanResourceForm';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <Dashboard /> },
      { path: 'clients', element: <ClientList /> },
      { path: 'clients/new', element: <ClientForm /> },
      { path: 'clients/:id', element: <ClientForm /> },
      { path: 'projects', element: <ProjectList /> },
      { path: 'projects/new', element: <ProjectForm /> },
      { path: 'human-resources', element: <HumanResourceList /> },
      { path: 'human-resources/new', element: <HumanResourceForm /> },
      { path: 'human-resources/:id', element: <HumanResourceForm /> },
    ],
  },
]);
