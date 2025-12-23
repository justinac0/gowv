#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    printf("(linux) CgoNativeWindowSetIcon called!\n");
}

void CgoNativeWindowHide(void* window_handle) {
    printf("(linux) CgoNativeWindowHide called!\n");
}

void CgoNativeWindowShow(void* window_handle) {
    printf("(linux) CgoNativeWindowShow called!\n");
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    printf("(linux) CgoNativeWindowSetMaximized called!\n");
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    printf("(linux) CgoNativeWindowSetMinimized called!\n");
}
