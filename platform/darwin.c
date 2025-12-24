#include "common.h"

#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    UNUSED(window_handle);
    UNUSED(filepath);
    printf("not supported on darwin yet");
}

void CgoNativeWindowHide(void* window_handle) {
    UNUSED(window_handle);
    printf("not supported on darwin yet");
}

void CgoNativeWindowShow(void* window_handle) {
    UNUSED(window_handle);
    printf("not supported on darwin yet");
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    UNUSED(window_handle);
    printf("not supported on darwin yet");
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    UNUSED(window_handle);
    printf("not supported on darwin yet");
}
