import {
  Database,
  Plus,
  Trash2,
  ChevronRight,
  ChevronDown,
  Table,
} from "lucide-react";
import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";
import type { Connection, TableInfo } from "../../types";

interface SidebarProps {
  connections: Connection[];
  tables: TableInfo[];
  selectedConnection: number | null;
  selectedTable: string | null;
  lastId: number;
  onSelectConnection: (id: number) => void;
  onSelectTable: (name: string) => void;
  onAddConnection: () => void;
  onDeleteConnection: (id: number) => void;
}

export function Sidebar({
  connections,
  tables,
  selectedConnection,
  selectedTable,
  lastId,
  onSelectConnection,
  onSelectTable,
  onAddConnection,
  onDeleteConnection,
}: SidebarProps) {
  return (
    <div className="w-64 h-full border-r border-zinc-200 dark:border-zinc-800 flex flex-col bg-zinc-50 dark:bg-zinc-950">
      <div className="p-4 border-b border-zinc-200 dark:border-zinc-800">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Database className="h-5 w-5 text-zinc-600 dark:text-zinc-400" />
            <span className="font-semibold">Databases</span>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={onAddConnection}
            className="h-8 w-8"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <ScrollArea className="flex-1">
        <div className="p-2 space-y-1">
          {connections?.length === 0 ? (
            <p className="text-sm text-zinc-500 text-center py-4">
              No connections
            </p>
          ) : (
            connections?.map((conn) => (
              <div key={conn.id}>
                <div
                  className={`
                    flex items-center justify-between px-2 py-2 rounded-md cursor-pointer text-sm
                    ${
                      selectedConnection === conn.id
                        ? "bg-zinc-200 dark:bg-zinc-800"
                        : "hover:bg-zinc-100 dark:hover:bg-zinc-900"
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
                      <ChevronDown className="h-4 w-4 shrink-0 text-zinc-500" />
                    ) : (
                      <ChevronRight className="h-4 w-4 shrink-0 text-zinc-500" />
                    )}
                    <span className="truncate">{conn.name}</span>
                    {lastId === conn.id && (
                      <span className="text-xs text-blue-500">&#8226;</span>
                    )}
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-6 w-6 opacity-0 group-hover:opacity-100 hover:opacity-100"
                    onClick={(e) => {
                      e.stopPropagation();
                      onDeleteConnection(conn.id);
                    }}
                  >
                    <Trash2 className="h-3 w-3 text-zinc-500" />
                  </Button>
                </div>

                {selectedConnection === conn.id && tables?.length > 0 && (
                  <div className="ml-4 mt-1 space-y-1">
                    {tables.map((table) => (
                      <div
                        key={table.name}
                        onClick={() => onSelectTable(table.name)}
                        className={`
                          flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer text-sm
                          ${
                            selectedTable === table.name
                              ? "bg-zinc-200 dark:bg-zinc-800"
                              : "hover:bg-zinc-100 dark:hover:bg-zinc-900"
                          }
                        `}
                      >
                        <Table className="h-3.5 w-3.5 text-zinc-500" />
                        <span className="truncate">{table.name}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </ScrollArea>
    </div>
  );
}
