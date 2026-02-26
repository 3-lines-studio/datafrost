import { X, Plus, Table, FileCode, Loader2 } from "lucide-react";
import { Button } from "../ui/button";
import type { Tab, TabType } from "../../types";

interface TabBarProps {
  tabs: Tab[];
  activeTabId: string | null;
  hasConnection: boolean;
  isLoading?: boolean;
  onTabClick: (id: string) => void;
  onTabClose: (id: string) => void;
  onNewQueryTab: () => void;
}

function getTabIcon(type: TabType) {
  if (type === "table") {
    return <Table className="h-3.5 w-3.5" />;
  }
  return <FileCode className="h-3.5 w-3.5" />;
}

export function TabBar({
  tabs,
  activeTabId,
  hasConnection,
  isLoading,
  onTabClick,
  onTabClose,
  onNewQueryTab,
}: TabBarProps) {
  return (
    <div className="flex items-center border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-950">
      <div className="flex-1 overflow-x-auto scrollbar-hide">
        <div className="flex items-center">
          {isLoading && tabs.length === 0 && (
            <div className="flex items-center gap-2 px-3 py-2 text-sm text-gray-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>Loading tabs...</span>
            </div>
          )}
          {tabs.map((tab) => {
            const isActive = activeTabId === tab.id;
            const Icon = getTabIcon(tab.type);

            return (
              <div
                key={tab.id}
                onClick={() => onTabClick(tab.id)}
                onAuxClick={(e) => {
                  if (e.button === 1) {
                    e.preventDefault();
                    onTabClose(tab.id);
                  }
                }}
                className={`
                  group flex items-center gap-2 px-3 py-2 cursor-pointer border-r border-gray-200 dark:border-gray-800
                  min-w-[120px] max-w-[200px] select-none
                  ${
                    isActive
                      ? "bg-white dark:bg-gray-900 border-b-0"
                      : "bg-gray-50 dark:bg-gray-950 hover:bg-gray-100 dark:hover:bg-gray-900"
                  }
                `}
              >
                <span className="text-gray-500">{Icon}</span>
                <span className="flex-1 truncate text-sm text-gray-700 dark:text-gray-300">
                  {tab.title}
                </span>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={(e) => {
                    e.stopPropagation();
                    onTabClose(tab.id);
                  }}
                  className="h-5 w-5 p-0 opacity-60 hover:opacity-100"
                >
                  <X className="h-3 w-3" />
                </Button>
              </div>
            );
          })}
        </div>
      </div>

      {hasConnection && (
        <div className="flex items-center px-2 border-l border-gray-200 dark:border-gray-800">
          <Button
            variant="ghost"
            size="icon"
            onClick={onNewQueryTab}
            className="h-7 w-7"
            title="New Query Tab"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>
      )}
    </div>
  );
}
