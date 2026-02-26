import { useMemo } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

import type { QueryResult } from "../../types";

interface ResultsTableProps {
  result: QueryResult | null;
  loading: boolean;
  error: string | null;
  onPageChange?: (page: number) => void;
}

function formatValue(value: any): string {
  if (value === null) return "NULL";
  if (value === undefined) return "";
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

export function ResultsTable({
  result,
  loading,
  error,
  onPageChange,
}: ResultsTableProps) {
  const totalPages = result?.total ? Math.ceil(result.total / result.limit) : 0;
  const currentPage = result?.page || 1;

  const canGoPrevious = currentPage > 1;
  const canGoNext = currentPage < totalPages;

  const handlePrevious = () => {
    if (canGoPrevious && onPageChange) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (canGoNext && onPageChange) {
      onPageChange(currentPage + 1);
    }
  };

  const content = useMemo(() => {
    if (loading) {
      return (
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">Loading...</span>
        </div>
      );
    }

    if (error) {
      return (
        <div className="flex items-center justify-center h-32 px-4">
          <span className="text-sm text-red-500">{error}</span>
        </div>
      );
    }

    if (!result) {
      return (
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">
            Run a query to see results
          </span>
        </div>
      );
    }

    if (!result.rows || result.rows.length === 0) {
      return (
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">No rows returned</span>
        </div>
      );
    }

    return (
      <div className="h-full flex flex-col">
        <div className="flex-1 overflow-auto">
          <Table>
            <TableHeader>
              <TableRow>
                {result.columns.map((col) => (
                  <TableHead key={col} className="whitespace-nowrap">
                    {col}
                  </TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              {result.rows.map((row, i) => (
                <TableRow key={i}>
                  {row.map((cell, j) => (
                    <TableCell key={j} className="max-w-xs truncate">
                      <span
                        className={cell === null ? "text-gray-400 italic" : ""}
                        title={formatValue(cell)}
                      >
                        {formatValue(cell)}
                      </span>
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
        <div className="h-10 flex items-center justify-between px-4 border-t border-gray-200 dark:border-gray-800 text-xs text-gray-500">
          <div>
            {result.count} of {result.total} rows
          </div>
          {onPageChange && totalPages > 1 && (
            <div className="flex items-center gap-2">
              <button
                onClick={handlePrevious}
                disabled={!canGoPrevious}
                className="px-2 py-1 rounded border border-gray-200 dark:border-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                Previous
              </button>
              <span className="px-2">
                Page {currentPage} of {totalPages}
              </span>
              <button
                onClick={handleNext}
                disabled={!canGoNext}
                className="px-2 py-1 rounded border border-gray-200 dark:border-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                Next
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }, [
    result,
    loading,
    error,
    onPageChange,
    totalPages,
    currentPage,
    canGoPrevious,
    canGoNext,
  ]);

  return (
    <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
      <div className="flex-1 min-h-0 min-w-0">{content}</div>
    </div>
  );
}
