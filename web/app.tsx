import "@/app.css";
import { useEffect, useState } from "react";
import { Sidebar } from "@/components/layout/sidebar";
import { ConnectionDialog } from "@/components/connections/connection-dialog";
import { QueryEditor } from "@/components/query/query-editor";
import { ResultsTable } from "@/components/query/results-table";
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
} from "@/lib/hooks";
import type { Connection } from "@/types";

interface AlertState {
  open: boolean;
  title: string;
  description: string;
  type: "info" | "success" | "error" | "confirm";
  onConfirm?: () => void;
}

export function Head() {
  return (
    <>
      <title>Datafrost - Turso Database UI</title>
      <meta
        name="description"
        content="A simple database UI for Turso libSQL"
      />
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
    selectedTable,
    showAddDialog,
    queryResult,
    theme,
    setSelectedConnection,
    setSelectedTable,
    setShowAddDialog,
    setQueryResult,
    clearSelection,
    setTheme,
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
  const { data: tables, isLoading: tablesLoading } =
    useTablesQuery(selectedConnection);
  const { data: tableData } = useTableDataQuery(
    selectedConnection,
    selectedTable,
  );
  const executeMutation = useExecuteQueryMutation(selectedConnection);
  const createMutation = useCreateConnectionMutation();
  const deleteMutation = useDeleteConnectionMutation();
  const setLastConnectedMutation = useSetLastConnectedMutation();
  const updateMutation = useUpdateConnectionMutation();
  const testMutation = useTestConnectionMutation();
  const testExistingMutation = useTestExistingConnectionMutation();

  const connections = connectionsData?.connections || [];
  const lastId = connectionsData?.last_id || 0;

  useEffect(() => {
    if (tableData) {
      setQueryResult(tableData);
    }
  }, [tableData, setQueryResult]);

  const handleSelectConnection = async (id: number) => {
    if (selectedConnection === id) {
      clearSelection();
    } else {
      setSelectedConnection(id);
      setSelectedTable(null);
      setQueryResult(null);
      await setLastConnectedMutation.mutateAsync(id);
    }
  };

  const handleSelectTable = (name: string) => {
    setSelectedTable(name);
  };

  const handleSaveConnection = async (
    name: string,
    url: string,
    token: string,
  ) => {
    if (dialogMode === "add") {
      await createMutation.mutateAsync({ name, url, token });
    } else if (dialogMode === "edit" && editingConnection) {
      await updateMutation.mutateAsync({
        id: editingConnection.id,
        data: { name, url, token },
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

  const handleExecuteQuery = async (query: string) => {
    try {
      const result = await executeMutation.mutateAsync(query);
      setQueryResult(result);
    } catch {
      setQueryResult(null);
    }
  };

  return (
    <ErrorBoundary>
      <div className="h-screen flex flex-col bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100">
        <ResizablePanelGroup
          orientation="horizontal"
          className="flex-1 overflow-hidden"
        >
          <ResizablePanel defaultSize={20}>
            <Sidebar
              connections={connections}
              tables={tables || []}
              selectedConnection={selectedConnection}
              selectedTable={selectedTable}
              lastId={lastId}
              theme={theme}
              onSelectConnection={handleSelectConnection}
              onSelectTable={handleSelectTable}
              onAddConnection={handleOpenAddDialog}
              onEditConnection={handleEditConnection}
              onDeleteConnection={handleDeleteConnection}
              onTestConnection={handleTestConnection}
              onToggleTheme={handleToggleTheme}
            />
          </ResizablePanel>

          <ResizableHandle
            withHandle
            className="bg-zinc-200 dark:bg-zinc-800"
          />

          <ResizablePanel defaultSize={80}>
            <ResizablePanelGroup orientation="vertical">
              <ResizablePanel defaultSize={20}>
                <QueryEditor
                  onExecute={handleExecuteQuery}
                  loading={executeMutation.isPending}
                />
              </ResizablePanel>

              <ResizableHandle
                withHandle
                className="bg-zinc-200 dark:bg-zinc-800"
              />

              <ResizablePanel defaultSize={80}>
                <ResultsTable
                  result={queryResult}
                  loading={
                    connectionsLoading ||
                    tablesLoading ||
                    executeMutation.isPending
                  }
                  error={executeMutation.error?.message || null}
                />
              </ResizablePanel>
            </ResizablePanelGroup>
          </ResizablePanel>
        </ResizablePanelGroup>

        <ConnectionDialog
          open={showAddDialog}
          onOpenChange={setShowAddDialog}
          mode={dialogMode}
          connection={editingConnection}
          onSave={handleSaveConnection}
          onTest={(url, token) => testMutation.mutateAsync({ url, token })}
          testLoading={testMutation.isPending}
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
