#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    printf("(windows) CgoNativeWindowSetIcon called!\n");
}

void CgoNativeWindowHide(void* window_handle) {
    printf("(windows) CgoNativeWindowHide called!\n");
}

void CgoNativeWindowShow(void* window_handle) {
    printf("(windows) CgoNativeWindowShow called!\n");
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    printf("(windows) CgoNativeWindowSetMaximized called!\n");
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    printf("(windows) CgoNativeWindowSetMinimized called!\n");
}
