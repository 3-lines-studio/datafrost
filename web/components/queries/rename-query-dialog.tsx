import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import type { SavedQuery } from "../../types";

interface RenameQueryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  query: SavedQuery | null;
  onRename: (name: string) => void;
  isLoading?: boolean;
}

export function RenameQueryDialog({
  open,
  onOpenChange,
  query,
  onRename,
  isLoading,
}: RenameQueryDialogProps) {
  const [name, setName] = useState(query?.name || "");

  useEffect(() => {
    if (query) {
      setName(query.name);
    }
  }, [query]);

  const handleRename = () => {
    if (name.trim()) {
      onRename(name.trim());
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleRename();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Rename Query</DialogTitle>
          <DialogDescription>
            Enter a new name for this saved query.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Query Name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Query name"
              autoFocus
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleRename} disabled={!name.trim() || isLoading}>
            {isLoading ? "Renaming..." : "Rename"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
