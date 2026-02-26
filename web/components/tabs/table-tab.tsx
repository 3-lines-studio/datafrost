import { ResultsTable } from "../query/results-table";
import { TableFilters } from "../query/table-filters";
import type { QueryResult, ColumnFilter } from "@/types";

interface TableTabProps {
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
  page?: number;
  filters?: ColumnFilter[];
  onPageChange?: (page: number) => void;
  onFiltersChange?: (filters: ColumnFilter[]) => void;
}

export function TableTab({
  result,
  loading,
  error,
  page,
  filters = [],
  onPageChange,
  onFiltersChange,
}: TableTabProps) {
  const columns = result?.columns || [];

  return (
    <div className="h-full flex flex-col">
      {onFiltersChange && (
        <TableFilters
          columns={columns}
          filters={filters}
          onFiltersChange={onFiltersChange}
        />
      )}
      <div className="flex-1 min-h-0">
        <ResultsTable
          result={result}
          loading={loading}
          error={error}
          onPageChange={onPageChange}
        />
      </div>
    </div>
  );
}
