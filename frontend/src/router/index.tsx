import { createBrowserRouter } from 'react-router-dom';
import MainLayout from '../layouts/MainLayout';
import Dashboard from '../pages/Dashboard';
import ClientList from '../pages/clients/ClientList';
import ClientForm from '../pages/clients/ClientForm';
import ProjectList from '../pages/projects/ProjectList';
import ProjectForm from '../pages/projects/ProjectForm';
import HumanResourceList from '../pages/human-resources/HumanResourceList';
import HumanResourceForm from '../pages/human-resources/HumanResourceForm';
import ResourceAllocationList from '../pages/resource-allocations/ResourceAllocationList';
import ResourceAllocationForm from '../pages/resource-allocations/ResourceAllocationForm';
import MilestoneList from '../pages/milestones/MilestoneList';
import MilestoneForm from '../pages/milestones/MilestoneForm';

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
      { path: 'projects/:id', element: <ProjectForm /> },
      { path: 'human-resources', element: <HumanResourceList /> },
      { path: 'human-resources/new', element: <HumanResourceForm /> },
      { path: 'human-resources/:id', element: <HumanResourceForm /> },
      { path: 'resource-allocations', element: <ResourceAllocationList /> },
      { path: 'resource-allocations/new', element: <ResourceAllocationForm /> },
      { path: 'resource-allocations/:id', element: <ResourceAllocationForm /> },
      { path: 'milestones', element: <MilestoneList /> },
      { path: 'milestones/new', element: <MilestoneForm /> },
      { path: 'milestones/:id', element: <MilestoneForm /> },
    ],
  },
]);
