import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";
import { Loader2, Key, AlertCircle } from "lucide-react";
import type { TableSchema } from "@/types";

interface TableSchemaViewProps {
  schema: TableSchema | null;
  loading: boolean;
  error: string | null;
}

export function TableSchemaView({
  schema,
  loading,
  error,
}: TableSchemaViewProps) {
  if (loading) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32">
          <Loader2 className="h-5 w-5 animate-spin text-gray-500" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex flex-col items-center justify-center h-32 px-4 gap-2">
          <AlertCircle className="h-5 w-5 text-red-500" />
          <span className="text-sm text-red-500">{error}</span>
        </div>
      </div>
    );
  }

  if (!schema) {
    return (
      <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800">
        <div className="flex items-center justify-center h-32">
          <span className="text-sm text-gray-500">No schema data</span>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full border-t border-gray-200 dark:border-gray-800 overflow-auto">
      <div className="p-4 space-y-6">
        <section>
          <h3 className="text-sm font-medium mb-3 text-gray-900 dark:text-gray-100">
            Columns
          </h3>
          <div className="border rounded-md overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-10"></TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Nullable</TableHead>
                  <TableHead>Default</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {schema.columns.map((col) => (
                  <TableRow key={col.name}>
                    <TableCell className="w-10">
                      {col.is_primary_key && (
                        <Key className="h-4 w-4 text-yellow-500" />
                      )}
                    </TableCell>
                    <TableCell className="font-medium">{col.name}</TableCell>
                    <TableCell className="text-gray-600 dark:text-gray-400">
                      {col.type}
                    </TableCell>
                    <TableCell>
                      {col.nullable ? (
                        <span className="text-green-600">Yes</span>
                      ) : (
                        <span className="text-red-600">No</span>
                      )}
                    </TableCell>
                    <TableCell className="text-gray-500 truncate max-w-xs">
                      {col.default_value || "-"}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </section>

        {schema.indexes && schema.indexes.length > 0 && (
          <section>
            <h3 className="text-sm font-medium mb-3 text-gray-900 dark:text-gray-100">
              Indexes
            </h3>
            <div className="border rounded-md overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Columns</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {schema.indexes.map((idx) => (
                    <TableRow key={idx.name}>
                      <TableCell className="font-medium">{idx.name}</TableCell>
                      <TableCell>
                        {idx.unique ? (
                          <span className="text-blue-600">Unique</span>
                        ) : (
                          <span>Index</span>
                        )}
                      </TableCell>
                      <TableCell>{idx.columns.join(", ")}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </section>
        )}

        {schema.constraints && schema.constraints.length > 0 && (
          <section>
            <h3 className="text-sm font-medium mb-3 text-gray-900 dark:text-gray-100">
              Constraints
            </h3>
            <div className="border rounded-md overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Column</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {schema.constraints.map((constraint) => (
                    <TableRow key={constraint.name}>
                      <TableCell className="font-medium">
                        {constraint.name}
                      </TableCell>
                      <TableCell>{constraint.type}</TableCell>
                      <TableCell>{constraint.column}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </section>
        )}
      </div>
    </div>
  );
}
