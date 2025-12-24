#include "common.h"

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    UNUSED(window_handle);
    UNUSED(filepath);
}

void CgoNativeWindowHide(void* controller_handle) {
    BOOL visible = ShowWindowAsync(controller_handle, SW_HIDE);
    UNUSED(visible);
}

void CgoNativeWindowShow(void* controller_handle) {
    BOOL visible = ShowWindowAsync(controller_handle, SW_SHOW);
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    BOOL visible = ShowWindowAsync(window_handle, SW_MAXIMIZE);
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    BOOL visible = ShowWindowAsync(window_handle, SW_MINIMIZE);
}
