import "@/app.css";
import { useEffect } from "react";
import { Sidebar } from "@/components/layout/sidebar";
import { AddConnectionDialog } from "@/components/connections/add-connection-dialog";
import { QueryEditor } from "@/components/query/query-editor";
import { ResultsTable } from "@/components/query/results-table";
import { ErrorBoundary } from "@/components/error-boundary";
import { QueryProvider } from "@/lib/query-provider";
import { useAppStore } from "@/lib/store";
import {
  useConnectionsQuery,
  useTablesQuery,
  useTableDataQuery,
  useExecuteQueryMutation,
  useCreateConnectionMutation,
  useDeleteConnectionMutation,
  useSetLastConnectedMutation,
} from "@/lib/hooks";

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
    setSelectedConnection,
    setSelectedTable,
    setShowAddDialog,
    setQueryResult,
    clearSelection,
  } = useAppStore();

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

  const handleAddConnection = async (
    name: string,
    url: string,
    token: string,
  ) => {
    await createMutation.mutateAsync({ name, url, token });
  };

  const handleDeleteConnection = async (id: number) => {
    if (confirm("Delete this connection?")) {
      await deleteMutation.mutateAsync(id);
      if (selectedConnection === id) {
        clearSelection();
      }
    }
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
        <div className="flex flex-1 overflow-hidden">
          <Sidebar
            connections={connections}
            tables={tables || []}
            selectedConnection={selectedConnection}
            selectedTable={selectedTable}
            lastId={lastId}
            onSelectConnection={handleSelectConnection}
            onSelectTable={handleSelectTable}
            onAddConnection={() => setShowAddDialog(true)}
            onDeleteConnection={handleDeleteConnection}
          />

          <div className="flex-1 min-w-0 flex flex-col">
            <div className="flex-1 min-h-0">
              <QueryEditor
                onExecute={handleExecuteQuery}
                loading={executeMutation.isPending}
              />
            </div>

            <div className="h-80">
              <ResultsTable
                result={queryResult}
                loading={
                  connectionsLoading ||
                  tablesLoading ||
                  executeMutation.isPending
                }
                error={executeMutation.error?.message || null}
              />
            </div>
          </div>
        </div>

        <AddConnectionDialog
          open={showAddDialog}
          onOpenChange={setShowAddDialog}
          onAdd={handleAddConnection}
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
