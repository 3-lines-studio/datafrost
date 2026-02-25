import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { ConnectionsResponse, QueryResult, TableInfo } from "../types";

const API_BASE = "";

const fetchConnections = async (): Promise<ConnectionsResponse> => {
  const res = await fetch(`${API_BASE}/api/connections`);
  if (!res.ok) throw new Error("Failed to fetch connections");
  return res.json();
};

const fetchTables = async (connectionId: number): Promise<TableInfo[]> => {
  const res = await fetch(`${API_BASE}/api/connections/${connectionId}/tables`);
  if (!res.ok) throw new Error("Failed to fetch tables");
  return res.json();
};

const fetchTableData = async (
  connectionId: number,
  tableName: string,
): Promise<QueryResult> => {
  const res = await fetch(
    `${API_BASE}/api/connections/${connectionId}/tables/${encodeURIComponent(tableName)}`,
  );
  if (!res.ok) throw new Error("Failed to fetch table data");
  return res.json();
};

const createConnectionApi = async (data: {
  name: string;
  url: string;
  token: string;
}): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Failed to create connection");
};

const deleteConnectionApi = async (id: number): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Failed to delete connection");
};

const setLastConnectedApi = async (id: number): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections/${id}/select`, {
    method: "POST",
  });
  if (!res.ok) throw new Error("Failed to set last connected");
};

const updateConnectionApi = async (
  id: number,
  data: { name: string; url: string; token: string },
): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error("Failed to update connection");
};

const testConnectionApi = async (data: {
  url: string;
  token: string;
}): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections/test`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Connection test failed");
  }
};

const testExistingConnectionApi = async (id: number): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/connections/${id}/test`, {
    method: "POST",
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Connection test failed");
  }
};

const executeQueryApi = async (
  connectionId: number,
  query: string,
): Promise<QueryResult> => {
  const res = await fetch(`${API_BASE}/api/connections/${connectionId}/query`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Query failed");
  }
  return res.json();
};

export function useConnectionsQuery() {
  return useQuery({
    queryKey: ["connections"],
    queryFn: fetchConnections,
  });
}

export function useTablesQuery(connectionId: number | null) {
  return useQuery({
    queryKey: ["tables", connectionId],
    queryFn: () => fetchTables(connectionId!),
    enabled: !!connectionId,
  });
}

export function useTableDataQuery(
  connectionId: number | null,
  tableName: string | null,
) {
  return useQuery({
    queryKey: ["tableData", connectionId, tableName],
    queryFn: () => fetchTableData(connectionId!, tableName!),
    enabled: !!connectionId && !!tableName,
  });
}

export function useCreateConnectionMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createConnectionApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
    },
  });
}

export function useDeleteConnectionMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteConnectionApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
    },
  });
}

export function useSetLastConnectedMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setLastConnectedApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
    },
  });
}

export function useUpdateConnectionMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: { name: string; url: string; token: string } }) =>
      updateConnectionApi(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
    },
  });
}

export function useTestConnectionMutation() {
  return useMutation({
    mutationFn: testConnectionApi,
  });
}

export function useTestExistingConnectionMutation() {
  return useMutation({
    mutationFn: testExistingConnectionApi,
  });
}

export function useExecuteQueryMutation(connectionId: number | null) {
  return useMutation({
    mutationFn: (query: string) => {
      if (!connectionId) throw new Error("No connection selected");
      return executeQueryApi(connectionId, query);
    },
  });
}

const fetchTheme = async (): Promise<{ theme: string }> => {
  const res = await fetch(`${API_BASE}/api/theme`);
  if (!res.ok) throw new Error("Failed to fetch theme");
  return res.json();
};

const updateThemeApi = async (theme: string): Promise<void> => {
  const res = await fetch(`${API_BASE}/api/theme`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ theme }),
  });
  if (!res.ok) throw new Error("Failed to update theme");
};

export function useThemeQuery() {
  return useQuery({
    queryKey: ["theme"],
    queryFn: fetchTheme,
  });
}

export function useUpdateThemeMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateThemeApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["theme"] });
    },
  });
}
