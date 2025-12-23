#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    printf("(darwin) CgoNativeWindowSetIcon called!\n");
}

void CgoNativeWindowHide(void* window_handle) {
    printf("(darwin) CgoNativeWindowHide called!\n");
}

void CgoNativeWindowShow(void* window_handle) {
    printf("(darwin) CgoNativeWindowShow called!\n");
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    printf("(darwin) CgoNativeWindowSetMaximized called!\n");
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    printf("(darwin) CgoNativeWindowSetMinimized called!\n");
}
