import { QueryEditor } from "../query/query-editor";
import { ResultsTable } from "../query/results-table";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "../ui/resizable";
import type { QueryResult } from "../../types";

interface QueryTabProps {
  query: string;
  onQueryChange: (query: string) => void;
  onExecute: (query: string) => Promise<void>;
  onSave?: () => void;
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
  executeLoading: boolean;
}

export function QueryTab({
  query,
  onQueryChange,
  onExecute,
  onSave,
  result,
  loading,
  error,
  executeLoading,
}: QueryTabProps) {
  return (
    <div className="h-full">
      <ResizablePanelGroup orientation="vertical">
        <ResizablePanel defaultSize={40} minSize={150}>
          <QueryEditor
            query={query}
            onQueryChange={onQueryChange}
            onExecute={onExecute}
            onSave={onSave}
            loading={executeLoading}
          />
        </ResizablePanel>

        <ResizableHandle withHandle className="bg-gray-200 dark:bg-gray-800" />

        <ResizablePanel defaultSize={60} minSize={150}>
          <ResultsTable result={result} loading={loading} error={error} />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
}
