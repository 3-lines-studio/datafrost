import { ResultsTable } from "../query/results-table";
import type { QueryResult } from "../../types";

interface TableTabProps {
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
  page?: number;
  onPageChange?: (page: number) => void;
}

export function TableTab({
  result,
  loading,
  error,
  page,
  onPageChange,
}: TableTabProps) {
  return (
    <div className="h-full">
      <ResultsTable
        result={result}
        loading={loading}
        error={error}
        onPageChange={onPageChange}
      />
    </div>
  );
}
