import { Loader2, Play } from "lucide-react";
import { useCallback, useState } from "react";
import * as EditorModule from "react-simple-code-editor";
import { Button } from "../ui/button";

const Editor =
  (EditorModule as any).default?.default ||
  (EditorModule as any).default ||
  EditorModule;
import { highlight, languages } from "prismjs/components/prism-core";

import "prismjs/components/prism-sql";
import "prismjs/themes/prism.css";

interface QueryEditorProps {
  onExecute: (query: string) => Promise<void>;
  loading: boolean;
}

const sqlHighlight = (code: string) => {
  return highlight(code, languages.sql, "sql");
};

export function QueryEditor({ onExecute, loading }: QueryEditorProps) {
  const [query, setQuery] = useState("SELECT * FROM ");

  const handleExecute = useCallback(async () => {
    if (!query.trim() || loading) return;
    await onExecute(query.trim());
  }, [query, loading, onExecute]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
        e.preventDefault();
        handleExecute();
      }
    },
    [handleExecute],
  );

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-4 py-2 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-950">
        <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
          Query Editor
        </span>
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

      <div
        className="flex-1 relative font-mono text-sm"
        onKeyDown={handleKeyDown}
      >
        <Editor
          value={query}
          onValueChange={setQuery}
          highlight={sqlHighlight}
          padding={16}
          className="font-mono text-sm min-h-full"
          textareaClassName="focus:outline-none"
          style={{
            fontFamily: '"Fira Code", "Monaco", "Consolas", monospace',
            fontSize: 14,
            backgroundColor: "transparent",
          }}
        />
      </div>

      <div className="px-4 py-1 text-xs text-zinc-500 border-t border-zinc-200 dark:border-zinc-800">
        Press Ctrl+Enter to run
      </div>
    </div>
  );
}
