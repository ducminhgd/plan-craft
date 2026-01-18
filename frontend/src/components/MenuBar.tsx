import { useState, useEffect } from 'react';
import { Menu, Modal, message } from 'antd';
import type { MenuProps } from 'antd';
import { Quit, BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { OpenDatabase, SaveDatabaseAs, GetCurrentDatabasePath, IsMemoryDatabase } from '../../wailsjs/go/services/DatabaseFileService';
import './MenuBar.css';

type MenuItem = Required<MenuProps>['items'][number];

// Detect if running on macOS
const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
const modKey = isMac ? '⌘' : 'Ctrl';

export default function MenuBar() {
  const [aboutModalOpen, setAboutModalOpen] = useState(false);
  const [currentDbPath, setCurrentDbPath] = useState<string>('');
  const [isDraft, setIsDraft] = useState<boolean>(false);

  // Fetch current database path and draft status on mount
  useEffect(() => {
    GetCurrentDatabasePath().then(setCurrentDbPath).catch(() => {});
    IsMemoryDatabase().then(setIsDraft).catch(() => {});
  }, []);

  // Register keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const modPressed = isMac ? e.metaKey : e.ctrlKey;

      if (modPressed && e.key === 'o') {
        e.preventDefault();
        handleOpenFile();
      } else if (modPressed && e.shiftKey && e.key === 's') {
        e.preventDefault();
        handleSaveAs();
      } else if (modPressed && e.key === 'q') {
        e.preventDefault();
        handleExit();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const handleExit = () => {
    if (isDraft) {
      Modal.confirm({
        title: 'Unsaved Changes',
        content: 'You have unsaved changes in draft mode. Would you like to save before exiting?',
        okText: 'Save',
        cancelText: 'Exit without saving',
        onOk: async () => {
          // Don't reload when saving before exit
          const filePath = await handleSaveAs(false);
          if (filePath) {
            Quit();
          }
        },
        onCancel: () => {
          Quit();
        },
      });
    } else {
      Modal.confirm({
        title: 'Exit Application',
        content: 'Are you sure you want to exit Plan Craft?',
        okText: 'Exit',
        cancelText: 'Cancel',
        onOk: () => {
          Quit();
        },
      });
    }
  };

  const handleOpenFile = async () => {
    try {
      const filePath = await OpenDatabase();
      if (filePath) {
        setCurrentDbPath(filePath);
        setIsDraft(false);
        message.success(`Switched to database: ${filePath.split(/[/\\]/).pop()}`);
        // Reload the page to refresh all data from the new database
        window.location.reload();
      }
    } catch (error) {
      message.error(`Failed to open database: ${error}`);
    }
  };

  const handleSaveAs = async (shouldReload: boolean = true): Promise<string | null> => {
    try {
      const filePath = await SaveDatabaseAs();
      if (filePath) {
        setCurrentDbPath(filePath);
        setIsDraft(false);
        message.success(`Database saved to: ${filePath.split(/[/\\]/).pop()}`);
        // Reload the page to refresh all data from the new database
        if (shouldReload) {
          window.location.reload();
        }
        return filePath;
      }
      return null;
    } catch (error) {
      message.error(`Failed to save database: ${error}`);
      return null;
    }
  };

  const handleGuides = () => {
    BrowserOpenURL('https://github.com/ducminhgd/plan-craft/wiki');
  };

  const handleAbout = () => {
    setAboutModalOpen(true);
  };

  const fileMenuItems: MenuItem[] = [
    {
      key: 'open',
      label: `Open file (${modKey}+O)`,
      onClick: handleOpenFile,
    },
    {
      key: 'save-as',
      label: `Save as (${modKey}+Shift+S)`,
      onClick: () => handleSaveAs(),
    },
    {
      type: 'divider',
    },
    {
      key: 'exit',
      label: `Exit (${modKey}+Q)`,
      onClick: handleExit,
    },
  ];

  const helpMenuItems: MenuItem[] = [
    {
      key: 'guides',
      label: 'Guides',
      onClick: handleGuides,
    },
    {
      type: 'divider',
    },
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
        {isDraft ? (
          <span className="current-db-path draft" title="Unsaved draft - use Save As to persist">
            draft
          </span>
        ) : currentDbPath && (
          <span className="current-db-path" title={currentDbPath}>
            {currentDbPath.split(/[/\\]/).pop()}
          </span>
        )}
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
