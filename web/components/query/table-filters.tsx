import { useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { ChevronDown, Plus, X } from "lucide-react";
import type { ColumnFilter, FilterOperator } from "../../types";

interface TableFiltersProps {
  columns: string[];
  filters: ColumnFilter[];
  onFiltersChange: (filters: ColumnFilter[]) => void;
}

const OPERATORS: { value: FilterOperator; label: string }[] = [
  { value: "eq", label: "=" },
  { value: "neq", label: "!=" },
  { value: "gt", label: ">" },
  { value: "lt", label: "<" },
  { value: "gte", label: ">=" },
  { value: "lte", label: "<=" },
  { value: "like", label: "LIKE" },
  { value: "not_like", label: "NOT LIKE" },
  { value: "is_null", label: "IS NULL" },
  { value: "is_not_null", label: "IS NOT NULL" },
];

function getOperatorLabel(operator: FilterOperator): string {
  return OPERATORS.find((op) => op.value === operator)?.label || operator;
}

function needsValue(operator: FilterOperator): boolean {
  return operator !== "is_null" && operator !== "is_not_null";
}

export function TableFilters({
  columns,
  filters,
  onFiltersChange,
}: TableFiltersProps) {
  const [newFilterColumn, setNewFilterColumn] = useState<string | null>(null);
  const [newFilterOperator, setNewFilterOperator] =
    useState<FilterOperator>("eq");
  const [newFilterValue, setNewFilterValue] = useState("");
  const [showAddFilter, setShowAddFilter] = useState(false);

  const addFilter = () => {
    if (!newFilterColumn) return;

    const filter: ColumnFilter = {
      id: crypto.randomUUID(),
      column: newFilterColumn,
      operator: newFilterOperator,
      value: needsValue(newFilterOperator) ? newFilterValue : "",
    };

    onFiltersChange([...filters, filter]);
    setNewFilterColumn(null);
    setNewFilterOperator("eq");
    setNewFilterValue("");
    setShowAddFilter(false);
  };

  const removeFilter = (id: string) => {
    onFiltersChange(filters.filter((f) => f.id !== id));
  };

  const clearAllFilters = () => {
    onFiltersChange([]);
  };

  const availableColumns = columns.filter(
    (col) => !filters.some((f) => f.column === col),
  );

  return (
    <div className="flex flex-col gap-2 p-3 border-b border-gray-200 dark:border-gray-800">
      <div className="flex items-center gap-2 flex-wrap">
        {filters.map((filter) => (
          <div
            key={filter.id}
            className="flex items-center gap-1.5 px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded-md text-sm"
          >
            <span className="font-medium">{filter.column}</span>
            <span className="text-gray-500">{getOperatorLabel(filter.operator)}</span>
            {needsValue(filter.operator) && (
              <span className="text-gray-700 dark:text-gray-300">{filter.value}</span>
            )}
            <button
              onClick={() => removeFilter(filter.id)}
              className="ml-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
            >
              <X className="size-3.5" />
            </button>
          </div>
        ))}

        {!showAddFilter && availableColumns.length > 0 && (
          <Button
            variant="outline"
            size="xs"
            onClick={() => setShowAddFilter(true)}
          >
            <Plus className="size-3.5 mr-1" />
            Filter
          </Button>
        )}

        {filters.length > 0 && (
          <Button
            variant="ghost"
            size="xs"
            onClick={clearAllFilters}
            className="text-gray-500"
          >
            Clear all
          </Button>
        )}
      </div>

      {showAddFilter && availableColumns.length > 0 && (
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="xs" className="min-w-[100px]">
                {newFilterColumn || "Column"}
                <ChevronDown className="size-3.5 ml-1" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              {availableColumns.map((col) => (
                <DropdownMenuItem
                  key={col}
                  onClick={() => setNewFilterColumn(col)}
                >
                  {col}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="xs" className="min-w-[80px]">
                {getOperatorLabel(newFilterOperator)}
                <ChevronDown className="size-3.5 ml-1" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              {OPERATORS.map((op) => (
                <DropdownMenuItem
                  key={op.value}
                  onClick={() => setNewFilterOperator(op.value)}
                >
                  {op.label}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          {needsValue(newFilterOperator) && (
            <Input
              type="text"
              value={newFilterValue}
              onChange={(e) => setNewFilterValue(e.target.value)}
              placeholder="Value"
              className="w-32 h-7 text-sm"
            />
          )}

          <Button
            size="xs"
            onClick={addFilter}
            disabled={!newFilterColumn || (needsValue(newFilterOperator) && !newFilterValue)}
          >
            Add
          </Button>

          <Button
            variant="ghost"
            size="xs"
            onClick={() => {
              setShowAddFilter(false);
              setNewFilterColumn(null);
              setNewFilterOperator("eq");
              setNewFilterValue("");
            }}
          >
            Cancel
          </Button>
        </div>
      )}
    </div>
  );
}
