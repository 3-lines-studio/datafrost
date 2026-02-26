import {
  Database,
  Plus,
  Trash2,
  Pencil,
  Activity,
  ChevronRight,
  ChevronDown,
  Table,
  Sun,
  Moon,
  MoreVertical,
  Loader2,
} from "lucide-react";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import type { Connection, TableInfo, SavedQuery } from "@/types";
import { SavedQueriesSection } from "../queries/saved-queries-section";

interface SidebarProps {
  connections: Connection[];
  tables: TableInfo[];
  tablesLoading: boolean;
  savedQueries: SavedQuery[];
  savedQueriesLoading: boolean;
  selectedConnection: number | null;
  lastId: number;
  theme: "light" | "dark";
  onSelectConnection: (id: number) => void;
  onSelectTable: (name: string) => void;
  onAddConnection: () => void;
  onEditConnection: (id: number) => void;
  onDeleteConnection: (id: number) => void;
  onTestConnection: (id: number) => void;
  onDisconnectConnection: () => void;
  onNewQuery: () => void;
  onOpenSavedQuery: (query: SavedQuery) => void;
  onRenameSavedQuery: (query: SavedQuery) => void;
  onDeleteSavedQuery: (query: SavedQuery) => void;
  onToggleTheme: () => void;
}

export function Sidebar({
  connections,
  tables,
  tablesLoading,
  savedQueries,
  savedQueriesLoading,
  selectedConnection,
  lastId,
  theme,
  onSelectConnection,
  onSelectTable,
  onAddConnection,
  onEditConnection,
  onDeleteConnection,
  onTestConnection,
  onDisconnectConnection,
  onNewQuery,
  onOpenSavedQuery,
  onRenameSavedQuery,
  onDeleteSavedQuery,
  onToggleTheme,
}: SidebarProps) {
  return (
    <div className="h-full border-r border-gray-200 dark:border-gray-800 flex flex-col bg-gray-50 dark:bg-gray-950">
      <div className="px-2 border-b border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-between py-1">
          <div className="flex items-center gap-2">
            <Database className="size-4 text-gray-600 dark:text-gray-400" />
            <span className="font-medium">Datafrost</span>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={onAddConnection}
            className="h-7 w-7"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        <div className="p-2 space-y-1">
          {connections?.length === 0 ? (
            <p className="text-sm text-gray-500 text-center py-4">
              No connections
            </p>
          ) : (
            connections?.map((conn) => (
              <div key={conn.id}>
                <div
                  className={`
                    group flex items-center justify-between px-2 py-2 rounded-md cursor-pointer text-sm
                    ${
                      selectedConnection === conn.id
                        ? "bg-gray-200 dark:bg-gray-800"
                        : "hover:bg-gray-100 dark:hover:bg-gray-900"
                    }
                    ${
                      lastId === conn.id && selectedConnection !== conn.id
                        ? "border-l-2 border-blue-500"
                        : ""
                    }
                  `}
                >
                  <div
                    className="flex items-center gap-2 flex-1 min-w-0"
                    onClick={() => onSelectConnection(conn.id)}
                  >
                    {selectedConnection === conn.id ? (
                      <ChevronDown className="h-4 w-4 shrink-0 text-gray-500" />
                    ) : (
                      <ChevronRight className="h-4 w-4 shrink-0 text-gray-500" />
                    )}
                    <span className="truncate">{conn.name}</span>
                    {lastId === conn.id && (
                      <span className="text-xs text-blue-500">&#8226;</span>
                    )}
                  </div>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <MoreVertical className="h-3 w-3 text-gray-500" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-40">
                      {selectedConnection === conn.id && (
                        <DropdownMenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            onDisconnectConnection();
                          }}
                        >
                          <Activity className="mr-2 h-4 w-4" />
                          Disconnect
                        </DropdownMenuItem>
                      )}
                      <DropdownMenuItem
                        onClick={(e) => {
                          e.stopPropagation();
                          onTestConnection(conn.id);
                        }}
                      >
                        <Activity className="mr-2 h-4 w-4" />
                        Test Connection
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={(e) => {
                          e.stopPropagation();
                          onEditConnection(conn.id);
                        }}
                      >
                        <Pencil className="mr-2 h-4 w-4" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={(e) => {
                          e.stopPropagation();
                          onDeleteConnection(conn.id);
                        }}
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                {selectedConnection === conn.id && (
                  <>
                    {tablesLoading ? (
                      <div className="mt-2 flex items-center justify-center py-4">
                        <Loader2 className="h-4 w-4 animate-spin text-gray-500" />
                      </div>
                    ) : tables?.length > 0 ? (
                      <div className="mt-1 space-y-1">
                        {tables.map((table) => (
                          <div
                            key={table.name}
                            onClick={() => onSelectTable(table.name)}
                            className="flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer text-sm hover:bg-gray-100 dark:hover:bg-gray-900"
                          >
                            <Table className="h-3.5 w-3.5 text-gray-500" />
                            <span className="truncate">{table.name}</span>
                          </div>
                        ))}
                      </div>
                    ) : null}
                    <SavedQueriesSection
                      queries={savedQueries}
                      isLoading={savedQueriesLoading}
                      onNewQuery={onNewQuery}
                      onOpenQuery={onOpenSavedQuery}
                      onRenameQuery={onRenameSavedQuery}
                      onDeleteQuery={onDeleteSavedQuery}
                    />
                  </>
                )}
              </div>
            ))
          )}
        </div>
      </div>

      <div className="p-2 border-t border-gray-200 dark:border-gray-800">
        <Button
          variant="ghost"
          size="icon"
          onClick={onToggleTheme}
          className="h-8 w-8"
        >
          {theme === "light" ? (
            <Moon className="h-4 w-4" />
          ) : (
            <Sun className="h-4 w-4" />
          )}
        </Button>
      </div>
    </div>
  );
}
