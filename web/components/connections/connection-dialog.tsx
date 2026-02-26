import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Eye, EyeOff, Loader2, Check, X, FileJson } from "lucide-react";
import type { Connection, AdapterInfo, FieldConfig, UIMode } from "../../types";

interface ConnectionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  mode: "add" | "edit";
  connection?: Connection;
  adapters: AdapterInfo[];
  adaptersLoading: boolean;
  onSave: (name: string, type: string, credentials: Record<string, any>) => Promise<void>;
  onTest: (type: string, credentials: Record<string, any>) => Promise<void>;
  testLoading: boolean;
}

export function ConnectionDialog({
  open,
  onOpenChange,
  mode,
  connection,
  adapters,
  adaptersLoading,
  onSave,
  onTest,
  testLoading,
}: ConnectionDialogProps) {
  const [name, setName] = useState("");
  const [selectedType, setSelectedType] = useState("");
  const [credentials, setCredentials] = useState<Record<string, any>>({});
  const [showPassword, setShowPassword] = useState<Record<string, boolean>>({});
  const [saveLoading, setSaveLoading] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);

  const selectedAdapter = adapters.find((a) => a.type === selectedType);
  const hasMultipleModes = selectedAdapter?.ui_config.modes && selectedAdapter.ui_config.modes.length > 0;
  const selectedMode = credentials.mode || (hasMultipleModes ? selectedAdapter?.ui_config.modes?.[0]?.key : undefined);
  const currentMode = hasMultipleModes
    ? selectedAdapter?.ui_config.modes?.find((m: UIMode) => m.key === selectedMode)
    : undefined;

  useEffect(() => {
    if (open) {
      if (mode === "edit" && connection) {
        setName(connection.name);
        setSelectedType(connection.type);
        setCredentials(connection.credentials || {});
      } else {
        setName("");
        setSelectedType("");
        setCredentials({});
      }
      setTestResult(null);
      setShowPassword({});
      setFileError(null);
    }
  }, [open, mode, connection]);

  useEffect(() => {
    if (selectedAdapter && hasMultipleModes) {
      const defaultMode = selectedAdapter.ui_config.modes?.[0]?.key;
      if (!credentials.mode && defaultMode) {
        setCredentials((prev) => ({ ...prev, mode: defaultMode }));
      }
    }
  }, [selectedAdapter, hasMultipleModes]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !selectedType) return;

    setSaveLoading(true);
    try {
      await onSave(name, selectedType, credentials);
      onOpenChange(false);
    } finally {
      setSaveLoading(false);
    }
  };

  const handleTest = async () => {
    if (!selectedType) return;

    setTestResult(null);
    try {
      await onTest(selectedType, credentials);
      setTestResult({ success: true, message: "Connection successful!" });
    } catch (err: any) {
      setTestResult({
        success: false,
        message: err.message || "Connection failed",
      });
    }
  };

  const handleFieldChange = (key: string, value: any) => {
    setCredentials((prev) => ({ ...prev, [key]: value }));
  };

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>, fieldKey: string) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith(".json")) {
      setFileError("Please upload a JSON file");
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      try {
        JSON.parse(content);
        handleFieldChange(fieldKey, content);
        setFileError(null);
      } catch {
        setFileError("Invalid JSON file");
      }
    };
    reader.readAsText(file);
  };

  const renderField = (field: FieldConfig) => {
    const value = credentials[field.key] || "";

    if (field.type === "textarea") {
      return (
        <div key={field.key} className="space-y-2">
          <Label htmlFor={field.key}>{field.label}</Label>
          <Textarea
            id={field.key}
            value={value}
            onChange={(e) => handleFieldChange(field.key, e.target.value)}
            placeholder={field.placeholder}
            required={field.required}
            className="min-h-[100px] font-mono text-sm"
          />
        </div>
      );
    }

    if (field.type === "password") {
      return (
        <div key={field.key} className="space-y-2">
          <Label htmlFor={field.key}>{field.label}</Label>
          <div className="relative">
            <Input
              id={field.key}
              type={showPassword[field.key] ? "text" : "password"}
              value={value}
              onChange={(e) => handleFieldChange(field.key, e.target.value)}
              placeholder={field.placeholder}
              required={field.required}
              className="pr-10"
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
              onClick={() =>
                setShowPassword((prev) => ({
                  ...prev,
                  [field.key]: !prev[field.key],
                }))
              }
            >
              {showPassword[field.key] ? (
                <EyeOff className="h-4 w-4 text-gray-500" />
              ) : (
                <Eye className="h-4 w-4 text-gray-500" />
              )}
            </Button>
          </div>
        </div>
      );
    }

    return (
      <div key={field.key} className="space-y-2">
        <Label htmlFor={field.key}>{field.label}</Label>
        <Input
          id={field.key}
          type={field.type === "number" ? "number" : "text"}
          value={value}
          onChange={(e) => handleFieldChange(field.key, e.target.value)}
          placeholder={field.placeholder}
          required={field.required}
        />
      </div>
    );
  };

  const getFieldsToValidate = () => {
    if (!selectedAdapter) return [];
    if (hasMultipleModes && currentMode) {
      return currentMode.fields;
    }
    return selectedAdapter.ui_config.fields || [];
  };

  const fieldsToValidate = getFieldsToValidate();
  const canTest = selectedType && fieldsToValidate.every((f: FieldConfig) => !f.required || credentials[f.key]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {mode === "add" ? "Add Database Connection" : "Edit Database Connection"}
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
            <Label>Database Type</Label>
            <Select
              value={selectedType}
              onValueChange={(value) => {
                setSelectedType(value);
                setCredentials({});
                setTestResult(null);
              }}
              disabled={mode === "edit" || adaptersLoading}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select database type..." />
              </SelectTrigger>
              <SelectContent>
                {adapters.map((adapter) => (
                  <SelectItem key={adapter.type} value={adapter.type}>
                    {adapter.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {selectedAdapter && (
              <p className="text-xs text-gray-500">{selectedAdapter.description}</p>
            )}
          </div>

          {selectedAdapter && (
            <>
              <div className="border-t pt-4">
                {hasMultipleModes && (
                  <div className="space-y-2 mb-4">
                    <Label>Connection Mode</Label>
                    <Select
                      value={selectedMode}
                      onValueChange={(value) => {
                        setCredentials({ mode: value });
                        setTestResult(null);
                      }}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {selectedAdapter.ui_config.modes?.map((m: UIMode) => (
                          <SelectItem key={m.key} value={m.key}>
                            {m.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                )}

                <h4 className="text-sm font-medium mb-4">Connection Details</h4>
                <div className="space-y-4">
                  {fieldsToValidate.map(renderField)}
                </div>

                {selectedAdapter.ui_config.supports_file && (
                  <div className="mt-4">
                    <Label htmlFor="file-upload" className="flex items-center gap-2 cursor-pointer">
                      <FileJson className="h-4 w-4" />
                      Upload JSON File
                    </Label>
                    <Input
                      id="file-upload"
                      type="file"
                      accept={selectedAdapter.ui_config.file_types?.join(",")}
                      onChange={(e) => {
                        const allFields = selectedAdapter.ui_config.modes
                          ?.flatMap((m: UIMode) => m.fields)
                          || selectedAdapter.ui_config.fields || [];
                        const credField = allFields.find((f: FieldConfig) => f.key === "credentials");
                        if (credField) {
                          handleFileUpload(e, credField.key);
                        }
                      }}
                      className="hidden"
                    />
                    {fileError && (
                      <p className="text-xs text-red-500 mt-1">{fileError}</p>
                    )}
                  </div>
                )}
              </div>
            </>
          )}

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
              disabled={!canTest || testLoading}
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
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={saveLoading || !name || !selectedType}>
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
