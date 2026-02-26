Datafrost VERSION_PLACEHOLDER - Database Management GUI

REQUIREMENTS:
- macOS 11.0 (Big Sur) or later
- Apple Silicon Mac (M1, M2, M3, or later)
- Intel Macs are not supported (can build from source if needed)

INSTALLATION:
1. Extract the archive:
   tar -xzf datafrost-macos-arm64-vVERSION_PLACEHOLDER.tar.gz

2. Copy binary to a directory in your PATH:
   sudo cp datafrost-macos-arm64-vVERSION_PLACEHOLDER/datafrost /usr/local/bin/
   chmod +x /usr/local/bin/datafrost

3. On first run, macOS may show a security warning because the app is unsigned:
   - Go to System Preferences > Security & Privacy
   - Click "Open Anyway"

USAGE:
Run: datafrost

Or double-click the binary from Finder.

VERIFICATION:
Check version: datafrost --version

TROUBLESHOOTING:
- If you see "cannot be opened because the developer cannot be verified":
  Right-click the binary and select "Open", then confirm

- Intel Mac users: This build is ARM64-only. You have two options:
  1. Use Rosetta 2 (if installed) - may work but not guaranteed
  2. Build from source: git clone <repo> && go build
