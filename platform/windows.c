#include "common.h"

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

// TODO: LoadImage requires cleanup
void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    HWND hwnd = window_handle;
    HICON icon = LoadImage(0, filepath, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_LOADFROMFILE);

    SendMessage(hwnd, WM_SETICON, ICON_SMALL, (LPARAM)icon);
    SendMessage(hwnd, WM_SETICON, ICON_BIG, (LPARAM)icon);

    SendMessage(GetWindow(hwnd, GW_OWNER), WM_SETICON, ICON_SMALL, (LPARAM)icon);
    SendMessage(GetWindow(hwnd, GW_OWNER), WM_SETICON, ICON_BIG, (LPARAM)icon);
}

void CgoNativeWindowHide(void* window_handle) {
    UNUSED(ShowWindowAsync(window_handle, SW_HIDE));
}

void CgoNativeWindowShow(void* window_handle) {
    UNUSED(ShowWindowAsync(window_handle, SW_SHOW));
    UNUSED(SetFocus(window_handle));
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    UNUSED(ShowWindowAsync(window_handle, SW_MAXIMIZE));
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    UNUSED(ShowWindowAsync(window_handle, SW_MINIMIZE));
}
