import { create } from "zustand";
import type { QueryResult } from "../types";

interface AppState {
  selectedConnection: number | null;
  selectedTable: string | null;
  showAddDialog: boolean;
  queryResult: QueryResult | null;
  setSelectedConnection: (id: number | null) => void;
  setSelectedTable: (name: string | null) => void;
  setShowAddDialog: (show: boolean) => void;
  setQueryResult: (result: QueryResult | null) => void;
  clearSelection: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  selectedConnection: null,
  selectedTable: null,
  showAddDialog: false,
  queryResult: null,
  setSelectedConnection: (id) => set({ selectedConnection: id }),
  setSelectedTable: (name) => set({ selectedTable: name }),
  setShowAddDialog: (show) => set({ showAddDialog: show }),
  setQueryResult: (result) => set({ queryResult: result }),
  clearSelection: () =>
    set({ selectedConnection: null, selectedTable: null, queryResult: null }),
}));
