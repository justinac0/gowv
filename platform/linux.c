#include "common.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stddef.h>

void CgoNativeWindowSetIcon(void* window_handle, const char* filepath) {
    gtk_window_set_icon_from_file(GTK_WINDOW(window_handle), filepath, NULL);
}

void CgoNativeWindowHide(void* window_handle) {
    gtk_widget_set_visible(GTK_WIDGET(window_handle), false);
}

void CgoNativeWindowShow(void* window_handle) {
    gtk_widget_set_visible(GTK_WIDGET(window_handle), true);
}

void CgoNativeWindowSetFullscreen(void* window_handle, bool fullscreen) {
    if (fullscreen) {
        gtk_window_fullscreen(GTK_WINDOW(window_handle));
    } else {
        gtk_window_unfullscreen(GTK_WINDOW(window_handle));
    }
}

void CgoNativeWindowSetMaximized(void* window_handle) {
    WINDOW_MAXIMIZE(GTK_WINDOW(window_handle));
}

void CgoNativeWindowSetMinimized(void* window_handle) {
    WINDOW_MINIMIZE(GTK_WINDOW(window_handle));
}
