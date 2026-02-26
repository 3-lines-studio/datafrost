import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Eye, EyeOff, Loader2, Check, X } from "lucide-react";
import type { Connection } from "../../types";

interface ConnectionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  mode: "add" | "edit";
  connection?: Connection;
  onSave: (name: string, url: string, token: string) => Promise<void>;
  onTest: (url: string, token: string) => Promise<void>;
  testLoading: boolean;
}

export function ConnectionDialog({
  open,
  onOpenChange,
  mode,
  connection,
  onSave,
  onTest,
  testLoading,
}: ConnectionDialogProps) {
  const [name, setName] = useState("");
  const [url, setUrl] = useState("");
  const [token, setToken] = useState("");
  const [showToken, setShowToken] = useState(false);
  const [saveLoading, setSaveLoading] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);

  useEffect(() => {
    if (open) {
      if (mode === "edit" && connection) {
        setName(connection.name);
        setUrl(connection.url);
        setToken(connection.token);
      } else {
        setName("");
        setUrl("");
        setToken("");
      }
      setTestResult(null);
      setShowToken(false);
    }
  }, [open, mode, connection]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !url) return;

    setSaveLoading(true);
    try {
      await onSave(name, url, token);
      onOpenChange(false);
    } finally {
      setSaveLoading(false);
    }
  };

  const handleTest = async () => {
    if (!url) return;

    setTestResult(null);
    try {
      await onTest(url, token);
      setTestResult({ success: true, message: "Connection successful!" });
    } catch (err: any) {
      setTestResult({
        success: false,
        message: err.message || "Connection failed",
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {mode === "add"
              ? "Add Database Connection"
              : "Edit Database Connection"}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSave} className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="My Database"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="url">URL</Label>
            <Input
              id="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="libsql://..."
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="token">Token</Label>
            <div className="relative">
              <Input
                id="token"
                type={showToken ? "text" : "password"}
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder="your-auth-token"
                className="pr-10"
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                onClick={() => setShowToken(!showToken)}
              >
                {showToken ? (
                  <EyeOff className="h-4 w-4 text-gray-500" />
                ) : (
                  <Eye className="h-4 w-4 text-gray-500" />
                )}
              </Button>
            </div>
          </div>

          {testResult && (
            <div
              className={`flex items-center gap-2 text-sm ${
                testResult.success ? "text-green-600" : "text-red-600"
              }`}
            >
              {testResult.success ? (
                <Check className="h-4 w-4" />
              ) : (
                <X className="h-4 w-4" />
              )}
              <span>{testResult.message}</span>
            </div>
          )}

          <div className="flex justify-between gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={handleTest}
              disabled={!url || testLoading}
            >
              {testLoading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Testing...
                </>
              ) : (
                "Test Connection"
              )}
            </Button>
            <div className="flex gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={saveLoading || !name || !url}>
                {saveLoading ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Saving...
                  </>
                ) : mode === "add" ? (
                  "Add Connection"
                ) : (
                  "Save Changes"
                )}
              </Button>
            </div>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
