import { create } from "zustand";
import type { QueryResult } from "../types";

type Theme = "light" | "dark";

interface AppState {
  selectedConnection: number | null;
  selectedTable: string | null;
  showAddDialog: boolean;
  queryResult: QueryResult | null;
  theme: Theme;
  setSelectedConnection: (id: number | null) => void;
  setSelectedTable: (name: string | null) => void;
  setShowAddDialog: (show: boolean) => void;
  setQueryResult: (result: QueryResult | null) => void;
  clearSelection: () => void;
  setTheme: (theme: Theme) => void;
}

export const useAppStore = create<AppState>((set) => ({
  selectedConnection: null,
  selectedTable: null,
  showAddDialog: false,
  queryResult: null,
  theme: "light",
  setSelectedConnection: (id) => set({ selectedConnection: id }),
  setSelectedTable: (name) => set({ selectedTable: name }),
  setShowAddDialog: (show) => set({ showAddDialog: show }),
  setQueryResult: (result) => set({ queryResult: result }),
  clearSelection: () =>
    set({ selectedConnection: null, selectedTable: null, queryResult: null }),
  setTheme: (theme) => set({ theme }),
}));
