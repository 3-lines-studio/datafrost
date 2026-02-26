import { AlignLeft, Loader2, Play, Save } from "lucide-react";
import { useCallback, useState } from "react";
import { format } from "sql-formatter";
import * as EditorModule from "react-simple-code-editor";
import { Button } from "../ui/button";

const Editor =
  (EditorModule as any).default?.default ||
  (EditorModule as any).default ||
  EditorModule;
import { highlight, languages } from "prismjs/components/prism-core";

import "prismjs/components/prism-sql";

interface QueryEditorProps {
  onExecute: (query: string) => Promise<void>;
  loading: boolean;
  query?: string;
  onQueryChange?: (query: string) => void;
  onSave?: () => void;
}

const sqlHighlight = (code: string) => {
  return highlight(code, languages.sql, "sql");
};

export function QueryEditor({
  onExecute,
  loading,
  query: controlledQuery,
  onQueryChange,
  onSave,
}: QueryEditorProps) {
  const [internalQuery, setInternalQuery] = useState("SELECT * FROM ");
  const query = controlledQuery !== undefined ? controlledQuery : internalQuery;
  const setQuery =
    onQueryChange !== undefined ? onQueryChange : setInternalQuery;

  const handleExecute = useCallback(async () => {
    if (!query.trim() || loading) return;
    await onExecute(query.trim());
  }, [query, loading, onExecute]);

  const handleFormat = useCallback(() => {
    if (!query.trim()) return;
    const formattedQuery = format(query);
    setQuery(formattedQuery);
  }, [query, setQuery]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
        e.preventDefault();
        handleExecute();
      } else if (e.shiftKey && e.altKey && e.key === "f") {
        e.preventDefault();
        handleFormat();
      }
    },
    [handleExecute, handleFormat],
  );

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-950">
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
          Query Editor
        </span>
        <div className="flex items-center gap-2">
          {onSave && (
            <Button
              size="sm"
              variant="outline"
              onClick={onSave}
              disabled={!query.trim()}
            >
              <Save className="h-4 w-4 mr-2" />
              Save
            </Button>
          )}
          <Button
            size="sm"
            variant="outline"
            onClick={handleFormat}
            disabled={!query.trim()}
          >
            <AlignLeft className="h-4 w-4 mr-2" />
            Format
          </Button>
          <Button
            size="sm"
            onClick={handleExecute}
            disabled={loading || !query.trim()}
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Running...
              </>
            ) : (
              <>
                <Play className="h-4 w-4 mr-2" />
                Run Query
              </>
            )}
          </Button>
        </div>
      </div>

      <div
        className="flex-1 relative font-mono text-sm overflow-auto"
        onKeyDown={handleKeyDown}
      >
        <Editor
          value={query}
          onValueChange={setQuery}
          highlight={sqlHighlight}
          padding={16}
          className="font-mono text-sm"
          textareaClassName="focus:outline-none"
          style={{
            fontFamily: '"IBM Plex Mono", monospace',
            fontSize: 14,
            backgroundColor: "transparent",
            overflow: "auto",
            height: "100%",
          }}
        />
      </div>

      <div className="px-4 py-1 text-xs text-gray-500 border-t border-gray-200 dark:border-gray-800">
        Ctrl+Enter to run · Shift+Alt+F to format
      </div>
    </div>
  );
}
