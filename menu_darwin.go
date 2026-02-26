package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include <objc/objc-runtime.h>
#include <Cocoa/Cocoa.h>

void setupMacMenu() {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        NSMenu *mainMenu = [[NSMenu alloc] init];

        // App Menu (shows as "Datafrost" in menu bar)
        NSMenuItem *appMenuItem = [[NSMenuItem alloc] init];
        [mainMenu addItem:appMenuItem];
        NSMenu *appMenu = [[NSMenu alloc] initWithTitle:@"Datafrost"];
        [appMenuItem setSubmenu:appMenu];

        // Quit Datafrost (Cmd+Q)
        NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:@"Quit Datafrost"
                                                            action:@selector(terminate:)
                                                     keyEquivalent:@"q"];
        [appMenu addItem:quitItem];

        // Edit Menu
        NSMenuItem *editMenuItem = [[NSMenuItem alloc] init];
        [mainMenu addItem:editMenuItem];
        NSMenu *editMenu = [[NSMenu alloc] initWithTitle:@"Edit"];
        [editMenuItem setSubmenu:editMenu];

        // Undo
        NSMenuItem *undoItem = [[NSMenuItem alloc] initWithTitle:@"Undo"
                                                            action:@selector(undo:)
                                                     keyEquivalent:@"z"];
        [editMenu addItem:undoItem];

        // Redo
        NSMenuItem *redoItem = [[NSMenuItem alloc] initWithTitle:@"Redo"
                                                            action:@selector(redo:)
                                                     keyEquivalent:@"Z"];
        [redoItem setKeyEquivalentModifierMask:NSEventModifierFlagCommand | NSEventModifierFlagShift];
        [editMenu addItem:redoItem];

        [editMenu addItem:[NSMenuItem separatorItem]];

        // Cut
        NSMenuItem *cutItem = [[NSMenuItem alloc] initWithTitle:@"Cut"
                                                           action:@selector(cut:)
                                                    keyEquivalent:@"x"];
        [editMenu addItem:cutItem];

        // Copy
        NSMenuItem *copyItem = [[NSMenuItem alloc] initWithTitle:@"Copy"
                                                            action:@selector(copy:)
                                                     keyEquivalent:@"c"];
        [editMenu addItem:copyItem];

        // Paste
        NSMenuItem *pasteItem = [[NSMenuItem alloc] initWithTitle:@"Paste"
                                                             action:@selector(paste:)
                                                      keyEquivalent:@"v"];
        [editMenu addItem:pasteItem];

        // Select All
        NSMenuItem *selectAllItem = [[NSMenuItem alloc] initWithTitle:@"Select All"
                                                                 action:@selector(selectAll:)
                                                          keyEquivalent:@"a"];
        [editMenu addItem:selectAllItem];

        [app setMainMenu:mainMenu];
    }
}
*/
import "C"

func setupMacEditMenu() {
	C.setupMacMenu()
}
