import * as React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "./dialog";
import { Button } from "./button";
import { AlertTriangle, CheckCircle, XCircle, Info } from "lucide-react";

interface AlertDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  type?: "info" | "success" | "error" | "confirm";
  onConfirm?: () => void;
  confirmText?: string;
  cancelText?: string;
}

export function AlertDialog({
  open,
  onOpenChange,
  title,
  description,
  type = "info",
  onConfirm,
  confirmText = "OK",
  cancelText = "Cancel",
}: AlertDialogProps) {
  const isConfirm = type === "confirm";

  const iconMap = {
    info: <Info className="h-5 w-5 text-blue-500" />,
    success: <CheckCircle className="h-5 w-5 text-green-500" />,
    error: <XCircle className="h-5 w-5 text-red-500" />,
    confirm: <AlertTriangle className="h-5 w-5 text-amber-500" />,
  };

  const handleConfirm = () => {
    onConfirm?.();
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <div className="flex items-center gap-3">
            {iconMap[type]}
            <DialogTitle>{title}</DialogTitle>
          </div>
        </DialogHeader>
        <div className="py-2">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            {description}
          </p>
        </div>
        <DialogFooter>
          <div className="flex gap-2 justify-end">
            {isConfirm && (
              <Button variant="outline" onClick={() => onOpenChange(false)}>
                {cancelText}
              </Button>
            )}
            <Button
              variant={isConfirm ? "destructive" : "default"}
              onClick={handleConfirm}
            >
              {confirmText}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
