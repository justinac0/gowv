#include "common.h"

#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
#include <winuser.h>

// TODO: LoadImage requires cleanup
void CgoNativeWindowSetIcon(void *window_handle, const char *filepath) {
  HWND hwnd = window_handle;
  HICON icon = LoadImage(0, filepath, IMAGE_ICON, 0, 0,
                         LR_DEFAULTSIZE | LR_LOADFROMFILE);

  SendMessage(hwnd, WM_SETICON, ICON_SMALL, (LPARAM)icon);
  SendMessage(hwnd, WM_SETICON, ICON_BIG, (LPARAM)icon);

  SendMessage(GetWindow(hwnd, GW_OWNER), WM_SETICON, ICON_SMALL, (LPARAM)icon);
  SendMessage(GetWindow(hwnd, GW_OWNER), WM_SETICON, ICON_BIG, (LPARAM)icon);
}

void CgoNativeWindowHide(void *window_handle) {
  UNUSED(ShowWindowAsync(window_handle, SW_HIDE));
}

void CgoNativeWindowShow(void *window_handle) {
  UNUSED(ShowWindowAsync(window_handle, SW_SHOW));
  UNUSED(SetFocus(window_handle));
}

static WINDOWPLACEMENT wp_prev = {sizeof(wp_prev)};
void CgoNativeWindowSetFullscreen(void *window_handle, bool fullscreen) {
  HWND hwnd = window_handle;
  DWORD style = GetWindowLong(hwnd, GWL_STYLE);

  if (fullscreen) {
    MONITORINFO mi = {sizeof(mi)};
    if (GetWindowPlacement(hwnd, &wp_prev) &&
        GetMonitorInfo(MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY),
                       &mi)) {
      SetWindowLongPtr(hwnd, GWL_STYLE, style & ~WS_OVERLAPPEDWINDOW);
      SetWindowPos(hwnd, HWND_TOP, mi.rcMonitor.left, mi.rcMonitor.top,
                   mi.rcMonitor.right - mi.rcMonitor.left,
                   mi.rcMonitor.bottom - mi.rcMonitor.top,
                   SWP_NOOWNERZORDER | SWP_FRAMECHANGED);
    }
  } else {
    SetWindowLongPtr(hwnd, GWL_STYLE, style | WS_OVERLAPPEDWINDOW);
    SetWindowPlacement(hwnd, &wp_prev);
    SetWindowPos(hwnd, NULL, 0, 0, 0, 0,
                 SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOOWNERZORDER |
                     SWP_FRAMECHANGED);
  }
}

void CgoNativeWindowSetMaximized(void *window_handle) {
  UNUSED(ShowWindowAsync(window_handle, SW_MAXIMIZE));
}

void CgoNativeWindowSetMinimized(void *window_handle) {
  UNUSED(ShowWindowAsync(window_handle, SW_MINIMIZE));
}
