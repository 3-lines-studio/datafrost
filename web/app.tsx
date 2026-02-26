import "@/app.css";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Sidebar } from "@/components/layout/sidebar";
import { ConnectionDialog } from "@/components/connections/connection-dialog";
import { TabBar } from "@/components/tabs/tab-bar";
import { TableTab } from "@/components/tabs/table-tab";
import { QueryTab } from "@/components/tabs/query-tab";
import { SaveQueryDialog } from "@/components/queries/save-query-dialog";
import { RenameQueryDialog } from "@/components/queries/rename-query-dialog";
import { ErrorBoundary } from "@/components/error-boundary";
import { QueryProvider } from "@/lib/query-provider";
import { useAppStore } from "@/lib/store";

import { AlertDialog } from "@/components/ui/alert-dialog";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import {
  useConnectionsQuery,
  useAdaptersQuery,
  useTablesQuery,
  useTableDataQuery,
  useExecuteQueryMutation,
  useCreateConnectionMutation,
  useDeleteConnectionMutation,
  useSetLastConnectedMutation,
  useThemeQuery,
  useUpdateThemeMutation,
  useUpdateConnectionMutation,
  useTestConnectionMutation,
  useTestExistingConnectionMutation,
  useLayoutQuery,
  useSaveLayoutMutation,
  useTabsQuery,
  useSaveTabsMutation,
  useSavedQueriesQuery,
  useCreateSavedQueryMutation,
  useUpdateSavedQueryMutation,
  useDeleteSavedQueryMutation,
} from "@/lib/hooks";
import type { Connection, QueryResult, Tab, SavedQuery } from "@/types";

interface AlertState {
  open: boolean;
  title: string;
  description: string;
  type: "info" | "success" | "error" | "confirm";
  onConfirm?: () => void;
}

interface TabResult {
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
}

export function Head() {
  return (
    <>
      <title>Datafrost - Turso Database UI</title>
      <meta name="description" content="A simple database UI" />
      <link
        href="https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500&display=swap"
        rel="stylesheet"
      />
    </>
  );
}

function PageContent() {
  const {
    selectedConnection,
    showAddDialog,
    theme,
    tabs,
    activeTabId,
    loadedConnectionId,
    hasTabsChanged,
    setSelectedConnection,
    setShowAddDialog,
    setTheme,
    setTabs,
    setActiveTabId,
    setLoadedConnectionId,
    resetTabsChanged,
    addTab,
    closeTab,
    updateTab,
    setTabPage,
    setTabFilters,
    clearSelection,
  } = useAppStore();

  const [dialogMode, setDialogMode] = useState<"add" | "edit">("add");
  const [editingConnection, setEditingConnection] = useState<
    Connection | undefined
  >(undefined);
  const [alertState, setAlertState] = useState<AlertState>({
    open: false,
    title: "",
    description: "",
    type: "info",
  });

  const [tabResults, setTabResults] = useState<Record<string, TabResult>>({});
  const [saveQueryDialogOpen, setSaveQueryDialogOpen] = useState(false);
  const [renameQueryDialogOpen, setRenameQueryDialogOpen] = useState(false);
  const [queryToRename, setQueryToRename] = useState<SavedQuery | null>(null);
  const [activeSavedQueryId, setActiveSavedQueryId] = useState<number | null>(
    null,
  );

  const { data: themeData, isLoading: themeLoading } = useThemeQuery();
  const updateThemeMutation = useUpdateThemeMutation();

  useEffect(() => {
    if (themeData?.theme) {
      setTheme(themeData.theme as "light" | "dark");
    }
  }, [themeData, setTheme]);

  useEffect(() => {
    if (theme === "dark") {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, [theme]);

  const handleToggleTheme = async () => {
    const newTheme = theme === "light" ? "dark" : "light";
    await updateThemeMutation.mutateAsync(newTheme);
    setTheme(newTheme);
  };

  const { data: connectionsData, isLoading: connectionsLoading } =
    useConnectionsQuery();
  const { data: adapters, isLoading: adaptersLoading } = useAdaptersQuery();
  const { data: tables, isLoading: tablesLoading } =
    useTablesQuery(selectedConnection);
  const executeMutation = useExecuteQueryMutation(selectedConnection);
  const createMutation = useCreateConnectionMutation();
  const deleteMutation = useDeleteConnectionMutation();
  const setLastConnectedMutation = useSetLastConnectedMutation();
  const updateMutation = useUpdateConnectionMutation();
  const testMutation = useTestConnectionMutation();
  const testExistingMutation = useTestExistingConnectionMutation();

  const connections = connectionsData?.connections || [];
  const lastId = connectionsData?.last_id || 0;

  const activeTab = useMemo(() => {
    return tabs.find((t) => t.id === activeTabId) || null;
  }, [tabs, activeTabId]);

  const { data: tableData } = useTableDataQuery(
    selectedConnection,
    activeTab?.type === "table" ? activeTab.tableName || null : null,
    activeTab?.type === "table" ? activeTab.page || 1 : 1,
    activeTab?.type === "table" ? activeTab.filters || [] : [],
  );

  useEffect(() => {
    if (activeTab?.type === "table" && tableData && activeTabId) {
      setTabResults((prev) => ({
        ...prev,
        [activeTabId]: {
          result: tableData,
          loading: false,
          error: null,
        },
      }));
    }
  }, [tableData, activeTab, activeTabId]);

  const { data: tabsData, isLoading: tabsLoading } =
    useTabsQuery(selectedConnection);
  const saveTabsMutation = useSaveTabsMutation();

  const { data: savedQueriesData, isLoading: savedQueriesLoading } =
    useSavedQueriesQuery(selectedConnection);
  const createSavedQueryMutation = useCreateSavedQueryMutation();
  const updateSavedQueryMutation = useUpdateSavedQueryMutation();
  const deleteSavedQueryMutation = useDeleteSavedQueryMutation();

  useEffect(() => {
    if (!selectedConnection) {
      setTabs([]);
      setActiveTabId(null);
      setLoadedConnectionId(null);
      return;
    }

    if (selectedConnection !== loadedConnectionId && tabsData?.tabs) {
      setTabs(tabsData.tabs);
      setLoadedConnectionId(selectedConnection);
      if (tabsData.tabs.length > 0) {
        setActiveTabId(tabsData.tabs[0].id);
      } else {
        setActiveTabId(null);
      }
    }
  }, [
    selectedConnection,
    tabsData,
    setTabs,
    setActiveTabId,
    setLoadedConnectionId,
    loadedConnectionId,
  ]);

  useEffect(() => {
    if (!selectedConnection || !hasTabsChanged) return;

    const timeout = setTimeout(() => {
      saveTabsMutation.mutate({ connectionId: selectedConnection, tabs });
      resetTabsChanged();
    }, 500);
    return () => clearTimeout(timeout);
  }, [tabs, selectedConnection, hasTabsChanged, resetTabsChanged]);

  const handleSelectConnection = async (id: number) => {
    if (selectedConnection === id) {
      clearSelection();
      return;
    }
    setSelectedConnection(id);
    await setLastConnectedMutation.mutateAsync(id);
  };

  const handleDisconnectConnection = () => {
    clearSelection();
  };

  const handleSaveQuery = async (name: string) => {
    if (!selectedConnection || !activeTab || activeTab.type !== "query") return;

    const queryText = activeTab.query || "";

    if (activeSavedQueryId) {
      await updateSavedQueryMutation.mutateAsync({
        connectionId: selectedConnection,
        queryId: activeSavedQueryId,
        name,
        query: queryText,
      });
    } else {
      const savedQuery = await createSavedQueryMutation.mutateAsync({
        connectionId: selectedConnection,
        name,
        query: queryText,
      });
      setActiveSavedQueryId(savedQuery.id);
      updateTab(activeTab.id, { title: name });
    }

    setSaveQueryDialogOpen(false);
  };

  const handleOpenSavedQuery = (savedQuery: SavedQuery) => {
    if (!selectedConnection) return;

    const existingTab = tabs.find(
      (t) =>
        t.type === "query" &&
        t.connectionId === selectedConnection &&
        t.query === savedQuery.query,
    );

    if (existingTab) {
      setActiveTabId(existingTab.id);
      setActiveSavedQueryId(savedQuery.id);
      return;
    }

    const isEmptyTab =
      activeTab?.type === "query" &&
      (!activeTab.query || activeTab.query === "");

    if (isEmptyTab && activeTab) {
      updateTab(activeTab.id, {
        query: savedQuery.query,
        title: savedQuery.name,
      });
      setActiveSavedQueryId(savedQuery.id);
    } else {
      const newTab: Tab = {
        id: crypto.randomUUID(),
        type: "query",
        title: savedQuery.name,
        connectionId: selectedConnection,
        query: savedQuery.query,
      };
      addTab(newTab);
      setActiveSavedQueryId(savedQuery.id);
      setTabResults((prev) => ({
        ...prev,
        [newTab.id]: {
          result: null,
          loading: false,
          error: null,
        },
      }));
    }
  };

  const handleRenameSavedQuery = (savedQuery: SavedQuery) => {
    setQueryToRename(savedQuery);
    setRenameQueryDialogOpen(true);
  };

  const handleConfirmRename = async (name: string) => {
    if (!selectedConnection || !queryToRename) return;

    await updateSavedQueryMutation.mutateAsync({
      connectionId: selectedConnection,
      queryId: queryToRename.id,
      name,
      query: queryToRename.query,
    });

    const tabToUpdate = tabs.find(
      (t) =>
        t.type === "query" &&
        t.connectionId === selectedConnection &&
        t.query === queryToRename.query,
    );

    if (tabToUpdate) {
      updateTab(tabToUpdate.id, { title: name });
    }

    setRenameQueryDialogOpen(false);
    setQueryToRename(null);
  };

  const handleDeleteSavedQuery = async (savedQuery: SavedQuery) => {
    if (!selectedConnection) return;

    setAlertState({
      open: true,
      title: "Delete Saved Query",
      description: `Are you sure you want to delete "${savedQuery.name}"? This action cannot be undone.`,
      type: "confirm",
      onConfirm: async () => {
        await deleteSavedQueryMutation.mutateAsync({
          connectionId: selectedConnection,
          queryId: savedQuery.id,
        });
      },
    });
  };

  const handleSelectTable = (name: string) => {
    if (!selectedConnection) return;

    const existingTab = tabs.find(
      (t) =>
        t.type === "table" &&
        t.connectionId === selectedConnection &&
        t.tableName === name,
    );

    if (existingTab) {
      setActiveTabId(existingTab.id);
      return;
    }

    const newTab: Tab = {
      id: crypto.randomUUID(),
      type: "table",
      title: name,
      connectionId: selectedConnection,
      tableName: name,
    };

    addTab(newTab);
    setTabResults((prev) => ({
      ...prev,
      [newTab.id]: {
        result: null,
        loading: true,
        error: null,
      },
    }));
  };

  const handleSaveConnection = async (
    name: string,
    type: string,
    credentials: Record<string, any>,
  ) => {
    if (dialogMode === "add") {
      await createMutation.mutateAsync({ name, type, credentials });
    } else if (dialogMode === "edit" && editingConnection) {
      await updateMutation.mutateAsync({
        id: editingConnection.id,
        data: { name, type, credentials },
      });
    }
  };

  const handleOpenAddDialog = () => {
    setDialogMode("add");
    setEditingConnection(undefined);
    setShowAddDialog(true);
  };

  const handleEditConnection = (id: number) => {
    const connection = connections.find((c) => c.id === id);
    if (connection) {
      setEditingConnection(connection);
      setDialogMode("edit");
      setShowAddDialog(true);
    }
  };

  const handleTestConnection = async (id: number) => {
    try {
      await testExistingMutation.mutateAsync(id);
      setAlertState({
        open: true,
        title: "Connection Test",
        description: "Connection test successful!",
        type: "success",
      });
    } catch (err: any) {
      setAlertState({
        open: true,
        title: "Connection Test Failed",
        description: err.message || "Connection test failed",
        type: "error",
      });
    }
  };

  const handleDeleteConnection = (id: number) => {
    setAlertState({
      open: true,
      title: "Delete Connection",
      description:
        "Are you sure you want to delete this connection? This action cannot be undone.",
      type: "confirm",
      onConfirm: async () => {
        await deleteMutation.mutateAsync(id);
        if (selectedConnection === id) {
          clearSelection();
        }
      },
    });
  };

  const handleExecuteQuery = useCallback(
    async (tabId: string, query: string) => {
      if (!selectedConnection) return;

      setTabResults((prev) => ({
        ...prev,
        [tabId]: {
          ...prev[tabId],
          loading: true,
          error: null,
        },
      }));

      try {
        const result = await executeMutation.mutateAsync(query);
        setTabResults((prev) => ({
          ...prev,
          [tabId]: {
            result,
            loading: false,
            error: null,
          },
        }));
      } catch (err: any) {
        setTabResults((prev) => ({
          ...prev,
          [tabId]: {
            result: null,
            loading: false,
            error: err.message || "Query failed",
          },
        }));
      }
    },
    [selectedConnection, executeMutation],
  );

  const handleNewQueryTab = () => {
    if (!selectedConnection) return;

    const newTab: Tab = {
      id: crypto.randomUUID(),
      type: "query",
      title: `Query ${tabs.filter((t) => t.type === "query").length + 1}`,
      connectionId: selectedConnection,
      query: "",
    };

    addTab(newTab);
    setTabResults((prev) => ({
      ...prev,
      [newTab.id]: {
        result: null,
        loading: false,
        error: null,
      },
    }));
  };

  const handleTabClick = (id: string) => {
    setActiveTabId(id);
  };

  const handleTabClose = (id: string) => {
    closeTab(id);
    setTabResults((prev) => {
      const newResults = { ...prev };
      delete newResults[id];
      return newResults;
    });
  };

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey) {
        if (e.key === "w") {
          e.preventDefault();
          if (activeTabId) {
            handleTabClose(activeTabId);
          }
        }
        if (e.key === "t") {
          e.preventDefault();
          if (selectedConnection) {
            handleNewQueryTab();
          }
        }
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [activeTabId, handleTabClose, handleNewQueryTab, selectedConnection]);

  const handleQueryChange = (tabId: string, query: string) => {
    updateTab(tabId, { query });
  };

  const { data: horizontalData, isLoading: horizontalLoading } =
    useLayoutQuery("horizontal");
  const saveLayoutMutation = useSaveLayoutMutation();

  if (horizontalLoading || themeLoading) {
    return null;
  }

  const horizontalLayout = horizontalData?.layout
    ? parseInt(horizontalData.layout)
    : 20;

  const handleHorizontalLayoutChange = (layout: { [id: string]: number }) => {
    const firstPanelSize = Object.values(layout)[0];
    if (firstPanelSize) {
      saveLayoutMutation.mutate({
        key: "horizontal",
        layout: firstPanelSize.toString(),
      });
    }
  };

  const currentTabResult = activeTabId ? tabResults[activeTabId] : null;

  const renderTabContent = () => {
    if (!activeTab || !selectedConnection) {
      return (
        <div className="flex items-center justify-center h-full text-gray-500">
          <p>Select a database connection to view tabs</p>
        </div>
      );
    }

    if (activeTab.type === "table") {
      return (
        <TableTab
          result={currentTabResult?.result || null}
          loading={currentTabResult?.loading || tablesLoading}
          error={currentTabResult?.error || null}
          filters={activeTab.filters}
          onFiltersChange={(filters) => setTabFilters(activeTab.id, filters)}
          onPageChange={(page) => setTabPage(activeTab.id, page)}
        />
      );
    }

    return (
      <QueryTab
        query={activeTab.query || ""}
        onQueryChange={(q) => {
          handleQueryChange(activeTab.id, q);
          if (q !== activeTab.query) {
            setActiveSavedQueryId(null);
          }
        }}
        onExecute={(q) => handleExecuteQuery(activeTab.id, q)}
        onSave={() => setSaveQueryDialogOpen(true)}
        result={currentTabResult?.result || null}
        loading={currentTabResult?.loading || false}
        error={currentTabResult?.error || null}
        executeLoading={executeMutation.isPending}
      />
    );
  };

  return (
    <ErrorBoundary>
      <div className="h-screen flex flex-col bg-white dark:bg-gray-950 text-gray-900 dark:text-gray-100">
        <ResizablePanelGroup
          orientation="horizontal"
          className="flex-1 overflow-hidden"
          onLayoutChanged={handleHorizontalLayoutChange}
        >
          <ResizablePanel
            defaultSize={horizontalLayout}
            minSize={200}
            maxSize={400}
          >
            <Sidebar
              connections={connections}
              tables={tables || []}
              tablesLoading={tablesLoading}
              savedQueries={savedQueriesData?.queries || []}
              savedQueriesLoading={savedQueriesLoading}
              selectedConnection={selectedConnection}
              lastId={lastId}
              theme={theme}
              onSelectConnection={handleSelectConnection}
              onSelectTable={handleSelectTable}
              onAddConnection={handleOpenAddDialog}
              onEditConnection={handleEditConnection}
              onDeleteConnection={handleDeleteConnection}
              onTestConnection={handleTestConnection}
              onDisconnectConnection={handleDisconnectConnection}
              onNewQuery={handleNewQueryTab}
              onOpenSavedQuery={handleOpenSavedQuery}
              onRenameSavedQuery={handleRenameSavedQuery}
              onDeleteSavedQuery={handleDeleteSavedQuery}
              onToggleTheme={handleToggleTheme}
            />
          </ResizablePanel>

          <ResizableHandle
            withHandle
            className="bg-gray-200 dark:bg-gray-800"
          />

          <ResizablePanel defaultSize={100 - horizontalLayout}>
            <div className="flex flex-col h-full">
              <TabBar
                tabs={tabs}
                activeTabId={activeTabId}
                hasConnection={!!selectedConnection}
                isLoading={tabsLoading}
                onTabClick={handleTabClick}
                onTabClose={handleTabClose}
                onNewQueryTab={handleNewQueryTab}
              />
              <div className="flex-1 min-h-0">{renderTabContent()}</div>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>

        <ConnectionDialog
          open={showAddDialog}
          onOpenChange={setShowAddDialog}
          mode={dialogMode}
          connection={editingConnection}
          adapters={adapters || []}
          adaptersLoading={adaptersLoading}
          onSave={handleSaveConnection}
          onTest={(type, credentials) =>
            testMutation.mutateAsync({ type, credentials })
          }
          testLoading={testMutation.isPending}
        />

        <SaveQueryDialog
          open={saveQueryDialogOpen}
          onOpenChange={setSaveQueryDialogOpen}
          defaultName={activeTab?.type === "query" ? activeTab.title : ""}
          onSave={handleSaveQuery}
          isLoading={
            createSavedQueryMutation.isPending ||
            updateSavedQueryMutation.isPending
          }
        />

        <RenameQueryDialog
          open={renameQueryDialogOpen}
          onOpenChange={setRenameQueryDialogOpen}
          query={queryToRename}
          onRename={handleConfirmRename}
          isLoading={updateSavedQueryMutation.isPending}
        />

        <AlertDialog
          open={alertState.open}
          onOpenChange={(open) => setAlertState({ ...alertState, open })}
          title={alertState.title}
          description={alertState.description}
          type={alertState.type}
          onConfirm={alertState.onConfirm}
        />
      </div>
    </ErrorBoundary>
  );
}

export function Page() {
  return (
    <QueryProvider>
      <PageContent />
    </QueryProvider>
  );
}
