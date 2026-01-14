import { useState } from 'react';
import { Menu, Modal } from 'antd';
import type { MenuProps } from 'antd';
import { Quit } from '../../wailsjs/runtime/runtime';
import './MenuBar.css';

type MenuItem = Required<MenuProps>['items'][number];

export default function MenuBar() {
  const [aboutModalOpen, setAboutModalOpen] = useState(false);

  const handleExit = () => {
    Modal.confirm({
      title: 'Exit Application',
      content: 'Are you sure you want to exit Plan Craft?',
      okText: 'Exit',
      cancelText: 'Cancel',
      onOk: () => {
        Quit();
      },
    });
  };

  const handleOpen = () => {
    // TODO: Implement directory opening functionality
    Modal.info({
      title: 'Coming Soon',
      content: 'Directory opening functionality will be implemented in a future version.',
    });
  };

  const handleAbout = () => {
    setAboutModalOpen(true);
  };

  const fileMenuItems: MenuItem[] = [
    {
      key: 'open',
      label: 'Open',
      onClick: handleOpen,
    },
    {
      type: 'divider',
    },
    {
      key: 'exit',
      label: 'Exit',
      onClick: handleExit,
    },
  ];

  const helpMenuItems: MenuItem[] = [
    {
      key: 'about',
      label: 'About',
      onClick: handleAbout,
    },
  ];

  const items: MenuItem[] = [
    {
      key: 'file',
      label: 'File',
      children: fileMenuItems,
    },
    {
      key: 'help',
      label: 'Help',
      children: helpMenuItems,
    },
  ];

  return (
    <>
      <div className="menu-bar">
        <Menu
          mode="horizontal"
          items={items}
          style={{
            border: 'none',
            borderBottom: '1px solid #f0f0f0',
          }}
        />
      </div>

      <Modal
        title="About Plan Craft"
        open={aboutModalOpen}
        onOk={() => setAboutModalOpen(false)}
        onCancel={() => setAboutModalOpen(false)}
        footer={[
          <button key="ok" onClick={() => setAboutModalOpen(false)}>
            OK
          </button>,
        ]}
      >
        <div style={{ padding: '20px 0' }}>
          <h3>Plan Craft</h3>
          <p>Version 1.0.0</p>
          <p>A desktop project management and estimation tool</p>
          <p style={{ marginTop: '20px' }}>
            Built with Go and Wails for efficient project planning,
            work breakdown structures, timeline estimation, resource planning,
            and cost estimation.
          </p>
          <p style={{ marginTop: '20px' }}>
            Contact: <a href="mailto:giaduongducminh@gmail.com">giaduongducminh@gmail.com</a>
          </p>
          <p style={{ marginTop: '20px' }}>
            © 2026 Plan Craft. All rights reserved.
          </p>
        </div>
      </Modal>
    </>
  );
}
