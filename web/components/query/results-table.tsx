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
}

function formatValue(value: any): string {
  if (value === null) return "NULL";
  if (value === undefined) return "";
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

export function ResultsTable({ result, loading, error }: ResultsTableProps) {
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
        <div className="h-8 flex items-center px-4 border-t border-gray-200 dark:border-gray-800 text-xs text-gray-500">
          {result.count} rows
          {result.limited && (
            <span className="ml-2 text-amber-600">(showing max 100)</span>
          )}
        </div>
      </div>
    );
  }, [result, loading, error]);

  return (
    <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
      <div className="px-4 py-2 border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-950">
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
          Results
        </span>
      </div>
      <div className="flex-1 min-h-0 min-w-0">{content}</div>
    </div>
  );
}
