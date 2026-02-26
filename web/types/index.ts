export interface Connection {
  id: number;
  name: string;
  type: string;
  credentials: Record<string, any>;
  created_at: string;
}

export interface TableInfo {
  name: string;
  type: string;
}

export interface QueryResult {
  columns: string[];
  rows: any[][];
  count: number;
  total: number;
  page: number;
  limit: number;
}

export interface ConnectionsResponse {
  connections: Connection[];
  last_id: number;
}

export interface CreateConnectionRequest {
  name: string;
  type: string;
  credentials: Record<string, any>;
}

export interface UpdateConnectionRequest {
  name: string;
  type: string;
  credentials: Record<string, any>;
}

export interface TestConnectionRequest {
  type: string;
  credentials: Record<string, any>;
}

export interface AdapterInfo {
  type: string;
  name: string;
  description: string;
  ui_config: UIConfig;
}

export interface UIConfig {
  modes?: UIMode[];
  fields?: FieldConfig[];
  supports_file: boolean;
  file_types?: string[];
}

export interface UIMode {
  key: string;
  label: string;
  fields: FieldConfig[];
}

export interface FieldConfig {
  key: string;
  label: string;
  type: string;
  required: boolean;
  placeholder?: string;
}

export type TabType = "table" | "query";

export interface Tab {
  id: string;
  type: TabType;
  title: string;
  connectionId: number;
  tableName?: string;
  query?: string;
  page?: number;
  filters?: ColumnFilter[];
}

export interface SavedQuery {
  id: number;
  connectionId: number;
  name: string;
  query: string;
  createdAt: string;
  updatedAt: string;
}

export type FilterOperator =
  | "eq"
  | "neq"
  | "gt"
  | "lt"
  | "gte"
  | "lte"
  | "like"
  | "not_like"
  | "is_null"
  | "is_not_null";

export interface ColumnFilter {
  id: string;
  column: string;
  operator: FilterOperator;
  value: string;
}
