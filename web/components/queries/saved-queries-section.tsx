import { FileCode, Plus, Pencil, Trash2 } from "lucide-react";
import { Button } from "../ui/button";
import type { SavedQuery } from "@/types";

interface SavedQueriesSectionProps {
  queries: SavedQuery[];
  isLoading: boolean;
  onNewQuery: () => void;
  onOpenQuery: (query: SavedQuery) => void;
  onRenameQuery: (query: SavedQuery) => void;
  onDeleteQuery: (query: SavedQuery) => void;
}

export function SavedQueriesSection({
  queries,
  isLoading,
  onNewQuery,
  onOpenQuery,
  onRenameQuery,
  onDeleteQuery,
}: SavedQueriesSectionProps) {
  return (
    <div className="mt-4 border-t border-gray-200 dark:border-gray-800 pt-4">
      <div className="flex items-center justify-between px-2 mb-2">
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
          Saved Queries
        </span>
        <Button
          variant="ghost"
          size="icon"
          onClick={onNewQuery}
          className="h-6 w-6"
          title="New Query"
        >
          <Plus className="h-3.5 w-3.5" />
        </Button>
      </div>

      {isLoading ? (
        <div className="px-2 py-2 text-sm text-gray-500">Loading...</div>
      ) : queries.length === 0 ? (
        <div className="px-2 py-2 text-sm text-gray-500 italic">
          No saved queries
        </div>
      ) : (
        <div className="space-y-1">
          {queries.map((query) => (
            <div
              key={query.id}
              className="group flex items-center justify-between px-2 py-1.5 rounded-md cursor-pointer text-sm hover:bg-gray-100 dark:hover:bg-gray-900"
            >
              <div
                className="flex items-center gap-2 flex-1 min-w-0"
                onClick={() => onOpenQuery(query)}
              >
                <FileCode className="h-3.5 w-3.5 text-gray-500 shrink-0" />
                <span className="truncate">{query.name}</span>
              </div>
              <div className="flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-5 w-5"
                  onClick={(e) => {
                    e.stopPropagation();
                    onRenameQuery(query);
                  }}
                >
                  <Pencil className="h-3 w-3 text-gray-500" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-5 w-5"
                  onClick={(e) => {
                    e.stopPropagation();
                    onDeleteQuery(query);
                  }}
                >
                  <Trash2 className="h-3 w-3 text-gray-500" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
