import { ResultsTable } from "../query/results-table";
import type { QueryResult } from "../../types";

interface TableTabProps {
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
}

export function TableTab({ result, loading, error }: TableTabProps) {
  return (
    <div className="h-full">
      <ResultsTable result={result} loading={loading} error={error} />
    </div>
  );
}
