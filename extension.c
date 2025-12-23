#include "webview.h"

#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    printf("CgoNativeWindowSetIcon called!\n");
}

void CgoNativeWindowHide(void* window_handle) {
    printf("CgoNativeWindowHide called!\n");
}

void CgoNativeWindowShow(void* window_handle) {
    printf("CgoNativeWindowShow called!\n");
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    printf("CgoNativeWindowSetMaximized called!\n");
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    printf("CgoNativeWindowSetMinimized called!\n");
}
