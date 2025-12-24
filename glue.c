#include "webview.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

// NOTE: WEBVIEW_PLATFORM_* macros don't get exported because this
// source is not compiled through a cpp compiler.
#if defined(__APPLE__)
#define GOWV_PLATFORM_DARWIN
#elif defined(__unix__)
#define GOWV_PLATFORM_LINUX

// NOTE: some function names change between gtk3/4
// namely *iconify -> *minimize
#include <gtk/gtk.h>

#if GTK_MAJOR_VERSION >= 4

#ifdef GDK_WINDOWING_X11
#include <gdk/x11/gdkx.h>
#endif

#warning "SetIcon won't work. Icons are application level since GTK4"
#define WINDOW_MINIMIZE(window) gtk_window_minimize(window)
#define WINDOW_MAXIMIZE(window) gtk_window_maximize(window)

#elif GTK_MAJOR_VERSION >= 3
#ifdef GDK_WINDOWING_X11
#include <gdk/gdkx.h>
#endif
#define WINDOW_MINIMIZE(window) gtk_window_iconify(window)
#define WINDOW_MAXIMIZE(window) gtk_window_maximize(window)
#endif

#elif defined(_WIN32)
#define GOWV_PLATFORM_WINDOWS
#else
#error "unable to detect current platform"
#endif

// NOTE: import platform specific source for handling gowv webview extension
// functionality.
#if defined(GOWV_PLATFORM_DARWIN)
#include "platform/darwin.c"
#elif defined(GOWV_PLATFORM_LINUX)
#include "platform/linux.c"
#elif defined(GOWV_PLATFORM_WINDOWS)
#include "platform/windows.c"
#else
#error "webview not supported on your platform"
#endif

// NOTE: the below code is straight from webview/webview_go with some minor changes
// mainly adding ability to check errors on webview_error_t returning functions.
struct binding_context {
    webview_t w;
    uintptr_t index;
};

void _webviewDispatchGoCallback(void *);
void _webviewBindingGoCallback(webview_t, char *, char *, uintptr_t);

static void _webview_dispatch_cb(webview_t w, void *arg) {
    _webviewDispatchGoCallback(arg);
}

static void _webview_binding_cb(const char *id, const char *req, void *arg) {
    struct binding_context *ctx = (struct binding_context *) arg;
    _webviewBindingGoCallback(ctx->w, (char *)id, (char *)req, ctx->index);
}

webview_error_t CgoWebViewDispatch(webview_t w, uintptr_t arg) {
    return webview_dispatch(w, _webview_dispatch_cb, (void *)arg);
}

webview_error_t CgoWebViewBind(webview_t w, const char *name, uintptr_t index) {
    struct binding_context *ctx = calloc(1, sizeof(struct binding_context));
    ctx->w = w;
    ctx->index = index;

    return webview_bind(w, name, _webview_binding_cb, (void *)ctx);
}

webview_error_t CgoWebViewUnbind(webview_t w, const char *name) {
    return webview_unbind(w, name);
}
