import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { GetCurrentDatabasePath, IsMemoryDatabase } from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

interface DatabaseContextType {
  currentDbPath: string;
  isDraft: boolean;
  refreshKey: number;
  triggerRefresh: () => void;
}

const DatabaseContext = createContext<DatabaseContextType | undefined>(undefined);

export function DatabaseProvider({ children }: { children: ReactNode }) {
  const [currentDbPath, setCurrentDbPath] = useState<string>('');
  const [isDraft, setIsDraft] = useState<boolean>(false);
  const [refreshKey, setRefreshKey] = useState<number>(0);

  const triggerRefresh = useCallback(() => {
    setRefreshKey((prev) => prev + 1);
  }, []);

  // Fetch current database path and draft status on mount
  useEffect(() => {
    GetCurrentDatabasePath().then(setCurrentDbPath).catch(() => {});
    IsMemoryDatabase().then(setIsDraft).catch(() => {});
  }, []);

  // Listen for database change events from the backend
  useEffect(() => {
    const handleDatabaseChanged = (dbPath: string, isMemory: boolean) => {
      setCurrentDbPath(dbPath);
      setIsDraft(isMemory);
      // Trigger a refresh for all components that depend on database data
      triggerRefresh();
    };

    EventsOn('database:changed', handleDatabaseChanged);

    return () => {
      EventsOff('database:changed');
    };
  }, [triggerRefresh]);

  return (
    <DatabaseContext.Provider value={{ currentDbPath, isDraft, refreshKey, triggerRefresh }}>
      {children}
    </DatabaseContext.Provider>
  );
}

export function useDatabase() {
  const context = useContext(DatabaseContext);
  if (context === undefined) {
    throw new Error('useDatabase must be used within a DatabaseProvider');
  }
  return context;
}
