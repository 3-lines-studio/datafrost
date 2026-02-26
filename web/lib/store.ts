import { create } from "zustand";
import type { QueryResult, Tab } from "../types";

type Theme = "light" | "dark";

interface AppState {
  selectedConnection: number | null;
  showAddDialog: boolean;
  theme: Theme;
  tabs: Tab[];
  activeTabId: string | null;
  loadedConnectionId: number | null;
  hasTabsChanged: boolean;
  setSelectedConnection: (id: number | null) => void;
  setShowAddDialog: (show: boolean) => void;
  setTheme: (theme: Theme) => void;
  setTabs: (tabs: Tab[]) => void;
  setActiveTabId: (id: string | null) => void;
  setLoadedConnectionId: (id: number | null) => void;
  resetTabsChanged: () => void;
  addTab: (tab: Tab) => void;
  closeTab: (id: string) => void;
  updateTab: (id: string, updates: Partial<Tab>) => void;
  getTab: (id: string) => Tab | undefined;
  findTabByTable: (connectionId: number, tableName: string) => Tab | undefined;
  clearSelection: () => void;
}

export const useAppStore = create<AppState>((set, get) => ({
  selectedConnection: null,
  showAddDialog: false,
  theme: "light",
  tabs: [],
  activeTabId: null,
  loadedConnectionId: null,
  hasTabsChanged: false,
  setSelectedConnection: (id) => set({ selectedConnection: id }),
  setShowAddDialog: (show) => set({ showAddDialog: show }),
  setTheme: (theme) => set({ theme }),
  setTabs: (tabs) => set({ tabs, hasTabsChanged: false }),
  setActiveTabId: (id) => set({ activeTabId: id }),
  setLoadedConnectionId: (id) => set({ loadedConnectionId: id }),
  resetTabsChanged: () => set({ hasTabsChanged: false }),
  addTab: (tab) => {
    const state = get();
    const existingTab = state.findTabByTable(tab.connectionId, tab.tableName || "");
    if (existingTab && tab.type === "table") {
      set({ activeTabId: existingTab.id });
      return;
    }
    set({
      tabs: [...state.tabs, tab],
      activeTabId: tab.id,
      hasTabsChanged: true,
    });
  },
  closeTab: (id) => {
    const state = get();
    const newTabs = state.tabs.filter((t) => t.id !== id);
    let newActiveId = state.activeTabId;
    if (state.activeTabId === id) {
      const closedIndex = state.tabs.findIndex((t) => t.id === id);
      if (newTabs.length > 0) {
        newActiveId = newTabs[Math.min(closedIndex, newTabs.length - 1)]?.id || null;
      } else {
        newActiveId = null;
      }
    }
    set({
      tabs: newTabs,
      activeTabId: newActiveId,
      hasTabsChanged: true,
    });
  },
  updateTab: (id, updates) => {
    const state = get();
    set({
      tabs: state.tabs.map((t) => (t.id === id ? { ...t, ...updates } : t)),
      hasTabsChanged: true,
    });
  },
  getTab: (id) => {
    return get().tabs.find((t) => t.id === id);
  },
  findTabByTable: (connectionId, tableName) => {
    return get().tabs.find(
      (t) => t.type === "table" && t.connectionId === connectionId && t.tableName === tableName,
    );
  },
  clearSelection: () =>
    set({ selectedConnection: null, showAddDialog: false }),
}));
