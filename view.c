#include "_cgo_export.h"
#include "view_my.h"
#include <gdk/gdkkeysyms.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

char *view_ui =
    "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
    "<!-- Generated with glade 3.20.0 -->\n"
    "<interface>\n"
    "  <requires lib=\"gtk+\" version=\"3.18\"/>\n"
    "  <object class=\"GtkWindow\" id=\"main-window\">\n"
    "    <property name=\"can_focus\">False</property>\n"
    "    <property name=\"border_width\">5</property>\n"
    "    <property name=\"title\" translatable=\"yes\">Alpino "
    "Viewer</property>\n"
    "    <property name=\"default_width\">1800</property>\n"
    "    <property name=\"default_height\">1100</property>\n"
    "    <property name=\"icon_name\">face-monkey</property>\n"
    "    <signal name=\"delete-event\" handler=\"delete_event\" "
    "swapped=\"no\"/>\n"
    "    <signal name=\"destroy\" handler=\"destroy\" swapped=\"no\"/>\n"
    "    <child>\n"
    "      <object class=\"GtkBox\" id=\"my-box\">\n"
    "        <property name=\"visible\">True</property>\n"
    "        <property name=\"can_focus\">False</property>\n"
    "        <property name=\"vexpand\">True</property>\n"
    "        <property name=\"orientation\">vertical</property>\n"
    "        <child>\n"
    "          <placeholder/>\n"
    "        </child>\n"
    "      </object>\n"
    "    </child>\n"
    "  </object>\n"
    "</interface>\n";

WebKitWebView *webview = NULL;

G_MODULE_EXPORT gboolean web_view_key_pressed(WebKitWebView *web_view,
                                              GdkEventKey *event,
                                              gpointer user_data) {

  if (event->keyval == GDK_KEY_q && (event->state & GDK_CONTROL_MASK)) {
    gtk_main_quit();
    return TRUE;
  }
  if (event->keyval == GDK_KEY_minus && (event->state & GDK_CONTROL_MASK)) {
    gdouble lvl;
    lvl = webkit_web_view_get_zoom_level (webview) - .05;
    if (lvl < .2) {
      lvl = .2;
    }
    webkit_web_view_set_zoom_level (webview, lvl);
  }
  if (event->keyval == GDK_KEY_equal && (event->state & GDK_CONTROL_MASK)) {
    gdouble lvl;
    lvl = webkit_web_view_get_zoom_level (webview) + .05;
    if (lvl > 3) {
      lvl = 3;
    }
    webkit_web_view_set_zoom_level (webview, lvl);
  }
  if (event->keyval == GDK_KEY_0 && (event->state & GDK_CONTROL_MASK)) {
    webkit_web_view_set_zoom_level (webview, 1);
  }

  return FALSE;
}

static void web_view_load_changed(WebKitWebView *web_view,
                                  WebKitLoadEvent load_event,
                                  gpointer user_data) {
  static char buf[1000];
  const gchar *provisional_uri, *redirected_uri, *uri;

  return;

  switch (load_event) {
  case WEBKIT_LOAD_STARTED:
    /* New load, we have now a provisional URI */
    provisional_uri = webkit_web_view_get_uri(web_view);
    printf("WEBKIT_LOAD_STARTED %s\n", provisional_uri);
    break;
  case WEBKIT_LOAD_REDIRECTED:
    redirected_uri = webkit_web_view_get_uri(web_view);
    printf("WEBKIT_LOAD_REDIRECTED %s\n", redirected_uri);
    break;
  case WEBKIT_LOAD_COMMITTED:
    /* The load is being performed. Current URI is
     * the final one and it won't change unless a new
     * load is requested or a navigation within the
     * same page is performed */
    uri = webkit_web_view_get_uri(web_view);
    printf("WEBKIT_LOAD_COMMITTED %s\n", uri);
    break;
  case WEBKIT_LOAD_FINISHED:
    /* Load finished */
    uri = webkit_web_view_get_uri(web_view);
    printf("WEBKIT_LOAD_FINISHED %s\n", uri);
    g_snprintf(buf, 999, "%s", uri);
    go_message(idLOADED, buf);
    break;
  }
}

void run(char const *url, char const *title) {
  static char buf[1000];
  GtkBuilder *builder;
  GError *error = NULL;
  GtkWidget *window, *box;
  WebKitSettings *settings = NULL;

  gtk_init(NULL, NULL);

  builder = gtk_builder_new();
  if (!gtk_builder_add_from_string(builder, view_ui, -1, &error)) {
    g_snprintf(buf, 999, "%s", error->message);
    go_message(idERROR, buf);
    return;
  }
  gtk_builder_connect_signals(builder, NULL);

  window = GTK_WIDGET(gtk_builder_get_object(builder, "main-window"));
  if (strlen(title) > 0) {
    gtk_window_set_title(GTK_WINDOW(window), title);
  }

  box = GTK_WIDGET(gtk_builder_get_object(builder, "my-box"));
  settings = webkit_settings_new();
  webkit_settings_set_default_font_size(settings, 18);
  webkit_settings_set_default_monospace_font_size(settings, 14);
  webkit_settings_set_default_charset(settings, "utf-8");
  webkit_settings_set_default_font_family(settings, "serif");
  webview = WEBKIT_WEB_VIEW(webkit_web_view_new_with_settings(settings));
  gtk_box_pack_start(GTK_BOX(box), GTK_WIDGET(webview), TRUE, TRUE, 0);
  g_signal_connect(webview, "load-changed", G_CALLBACK(web_view_load_changed),
                   NULL);
  g_signal_connect(webview, "key-press-event", G_CALLBACK(web_view_key_pressed),
                   NULL);
  webkit_web_view_load_uri(webview, url);

  go_message(idREADY, "Let's begin!");
  gtk_widget_show_all(window);
  gtk_main();
}

G_MODULE_EXPORT gboolean delete_event(GtkWidget *widget, GdkEvent *event,
                                      gpointer data) {
  /* If you return FALSE in the "delete-event" signal handler,
   * GTK will emit the "destroy" signal. Returning TRUE means
   * you don't want the window to be destroyed.
   * This is useful for popping up 'are you sure you want to quit?'
   * type dialogs. */

  go_message(idDELETE, "Delete!");

  return FALSE;
}

G_MODULE_EXPORT void destroy(GtkWidget *widget, gpointer data) {
  gtk_main_quit();

  go_message(idDESTROY, "Destroy!");
}
