#include "common.h"

#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    UNUSED(window_handle);
    UNUSED(filepath);
}

void CgoNativeWindowHide(void* controller_handle) {
    UNUSED(controller_handle);
}

void CgoNativeWindowShow(void* controller_handle) {
    UNUSED(controller_handle);
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    UNUSED(window_handle);
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    UNUSED(window_handle);
}
