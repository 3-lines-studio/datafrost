Datafrost VERSION_PLACEHOLDER - Database Management GUI

PREREQUISITES:
This application requires WebKit2GTK to be installed on your system.

Install on Ubuntu/Debian:
  sudo apt-get install libwebkit2gtk-4.0-37 libgtk-3-0

Install on Fedora/RHEL:
  sudo dnf install webkit2gtk3 gtk3

Install on Arch:
  sudo pacman -S webkit2gtk gtk3

INSTALLATION:
1. Copy binary to /usr/local/bin:
   sudo cp datafrost /usr/local/bin/
   sudo chmod +x /usr/local/bin/datafrost

2. Copy icon to pixmaps:
   sudo cp icon.png /usr/local/share/pixmaps/datafrost.png

3. Copy desktop entry:
   sudo cp datafrost.desktop /usr/local/share/applications/

4. Update desktop database:
   sudo update-desktop-database /usr/local/share/applications/

USAGE:
Run: datafrost

Or find "Datafrost" in your application menu.

VERIFICATION:
Check version: datafrost --version

TROUBLESHOOTING:
- If the application doesn't start, check that WebKit2GTK is installed
- Look for error messages in the terminal
- Check your system's package manager for webkit2gtk/gtk3 packages
