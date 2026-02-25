export interface Connection {
  id: number;
  name: string;
  url: string;
  token: string;
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
  limited: boolean;
}

export interface ConnectionsResponse {
  connections: Connection[];
  last_id: number;
}

export interface UpdateConnectionRequest {
  name: string;
  url: string;
  token: string;
}

export interface TestConnectionRequest {
  url: string;
  token: string;
}
