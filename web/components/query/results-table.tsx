import { useRef } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

import type { QueryResult } from "@/types";

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

function VirtualTable({ result, onPageChange }: { result: QueryResult; onPageChange?: (page: number) => void }) {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: result.rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 36,
    overscan: 10,
  });

  const virtualItems = virtualizer.getVirtualItems();
  const totalSize = virtualizer.getTotalSize();
  const paddingTop = virtualItems.length > 0 ? virtualItems[0].start : 0;
  const paddingBottom =
    virtualItems.length > 0
      ? totalSize - virtualItems[virtualItems.length - 1].end
      : 0;

  const totalPages = result.total ? Math.ceil(result.total / result.limit) : 0;
  const currentPage = result.page || 1;

  return (
    <div className="h-full flex flex-col">
      <div ref={parentRef} className="flex-1 overflow-auto">
        <Table>
          <TableHeader className="sticky top-0 z-10 bg-white dark:bg-gray-950">
            <TableRow>
              {result.columns.map((col) => (
                <TableHead key={col} className="whitespace-nowrap">
                  {col}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {paddingTop > 0 && (
              <TableRow>
                <TableCell
                  colSpan={result.columns.length}
                  style={{ height: paddingTop, padding: 0 }}
                />
              </TableRow>
            )}
            {virtualItems.map((virtualRow) => {
              const row = result.rows[virtualRow.index];
              return (
                <TableRow key={virtualRow.index}>
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
              );
            })}
            {paddingBottom > 0 && (
              <TableRow>
                <TableCell
                  colSpan={result.columns.length}
                  style={{ height: paddingBottom, padding: 0 }}
                />
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="h-10 flex items-center justify-between px-4 border-t border-gray-200 dark:border-gray-800 text-xs text-gray-500 shrink-0">
        <div>
          {result.count} of {result.total} rows
        </div>
        {onPageChange && totalPages > 1 && (
          <div className="flex items-center gap-2">
            <button
              onClick={() => onPageChange(currentPage - 1)}
              disabled={currentPage <= 1}
              className="px-2 py-1 rounded border border-gray-200 dark:border-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              Previous
            </button>
            <span className="px-2">
              Page {currentPage} of {totalPages}
            </span>
            <button
              onClick={() => onPageChange(currentPage + 1)}
              disabled={currentPage >= totalPages}
              className="px-2 py-1 rounded border border-gray-200 dark:border-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              Next
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export function ResultsTable({
  result,
  loading,
  error,
  onPageChange,
}: ResultsTableProps) {
  if (loading && !result) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">Loading...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32 px-4">
          <span className="text-sm text-red-500">{error}</span>
        </div>
      </div>
    );
  }

  if (!result) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">Run a query to see results</span>
        </div>
      </div>
    );
  }

  if (!result.rows || result.rows.length === 0) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">No rows returned</span>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
      <VirtualTable result={result} onPageChange={onPageChange} />
    </div>
  );
}
