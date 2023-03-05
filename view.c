#include "_cgo_export.h"
#include "view_my.h"
#include <gdk/gdkkeysyms.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

char *view_ui_single =
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

char *view_ui_multi = "\
<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n\
<!-- Generated with glade 3.40.0 -->\n\
<interface>\n\
  <requires lib=\"gtk+\" version=\"3.18\"/>\n\
  <object class=\"GtkWindow\" id=\"main-window\">\n\
    <property name=\"can-focus\">False</property>\n\
    <property name=\"border-width\">5</property>\n\
    <property name=\"title\" translatable=\"yes\">Alpino Viewer</property>\n\
    <property name=\"default-width\">1600</property>\n\
    <property name=\"default-height\">1000</property>\n\
    <property name=\"icon-name\">face-monkey</property>\n\
    <signal name=\"delete-event\" handler=\"delete_event\" swapped=\"no\"/>\n\
    <signal name=\"destroy\" handler=\"destroy\" swapped=\"no\"/>\n\
    <child>\n\
      <object class=\"GtkPaned\" id=\"panels\">\n\
        <property name=\"visible\">True</property>\n\
        <property name=\"can-focus\">True</property>\n\
        <child>\n\
          <object class=\"GtkTreeView\" id=\"files\">\n\
            <property name=\"visible\">True</property>\n\
            <property name=\"can-focus\">True</property>\n\
            <child internal-child=\"selection\">\n\
              <object class=\"GtkTreeSelection\" id=\"file\"/>\n\
            </child>\n\
          </object>\n\
          <packing>\n\
            <property name=\"resize\">False</property>\n\
            <property name=\"shrink\">True</property>\n\
          </packing>\n\
        </child>\n\
        <child>\n\
          <object class=\"GtkBox\" id=\"my-box\">\n\
            <property name=\"visible\">True</property>\n\
            <property name=\"can-focus\">False</property>\n\
            <property name=\"orientation\">vertical</property>\n\
            <child>\n\
              <placeholder/>\n\
            </child>\n\
          </object>\n\
          <packing>\n\
            <property name=\"resize\">True</property>\n\
            <property name=\"shrink\">True</property>\n\
          </packing>\n\
        </child>\n\
      </object>\n\
    </child>\n\
  </object>\n\
</interface>\n";

int nfiles = 0;
char const **filenames = NULL;

WebKitWebView *webview = NULL;
GtkWidget *window = NULL;

G_MODULE_EXPORT void tree_row_activated_cb(GtkTreeView *treeview,
                                           GtkTreePath *path,
                                           GtkTreeViewColumn *column,
                                           gpointer userdata) {
  GtkTreeIter iter;
  GtkTreeModel *model;
  gchar *item;

  model = gtk_tree_view_get_model(treeview);
  if (gtk_tree_model_get_iter(model, &iter, path)) {
    gtk_tree_model_get(model, &iter, 0, &item, -1);

    go_message(idSELECT, item);

    g_free(item);
  }
}

G_MODULE_EXPORT gboolean web_view_key_pressed(WebKitWebView *web_view,
                                              GdkEventKey *event,
                                              gpointer user_data) {

  if (event->keyval == GDK_KEY_q && (event->state & GDK_CONTROL_MASK)) {
    gtk_main_quit();
    return TRUE;
  }
  if (event->keyval == GDK_KEY_minus && (event->state & GDK_CONTROL_MASK)) {
    gdouble lvl;
    lvl = webkit_web_view_get_zoom_level(webview) - .05;
    if (lvl < .2) {
      lvl = .2;
    }
    webkit_web_view_set_zoom_level(webview, lvl);
  }
  if (event->keyval == GDK_KEY_equal && (event->state & GDK_CONTROL_MASK)) {
    gdouble lvl;
    lvl = webkit_web_view_get_zoom_level(webview) + .05;
    if (lvl > 3) {
      lvl = 3;
    }
    webkit_web_view_set_zoom_level(webview, lvl);
  }
  if (event->keyval == GDK_KEY_0 && (event->state & GDK_CONTROL_MASK)) {
    webkit_web_view_set_zoom_level(webview, 1);
  }

  return FALSE;
}

static GtkTreeModel *create_and_fill_model(void) {
  GtkListStore *store = gtk_list_store_new(1, G_TYPE_STRING);

  /* Append a row and fill in some data */
  GtkTreeIter iter;
  for (int i = 0; i < nfiles; i++) {
    gtk_list_store_append(store, &iter);
    gtk_list_store_set(store, &iter, 0, filenames[i], -1);
  }

  return GTK_TREE_MODEL(store);
}

void setnfiles(int n) {
  filenames = (char const **)malloc(n * sizeof(char const *));
}

void addfile(char const *filename) { filenames[nfiles++] = filename; }

void reload(char const *title) {
  gtk_window_set_title(GTK_WINDOW(window), title);
  webkit_web_view_reload(webview);
}

void run(char const *url, char const *title) {
  static char buf[1000];
  GtkBuilder *builder;
  GError *error = NULL;
  GtkWidget *box, *files;
  WebKitSettings *settings = NULL;

  gtk_init(NULL, NULL);

  builder = gtk_builder_new();
  if (!gtk_builder_add_from_string(
          builder, nfiles > 1 ? view_ui_multi : view_ui_single, -1, &error)) {
    g_snprintf(buf, 999, "%s", error->message);
    go_message(idERROR, buf);
    return;
  }
  gtk_builder_connect_signals(builder, NULL);

  window = GTK_WIDGET(gtk_builder_get_object(builder, "main-window"));
  if (strlen(title) > 0) {
    gtk_window_set_title(GTK_WINDOW(window), title);
  }

  if (nfiles > 1) {
    files = GTK_WIDGET(gtk_builder_get_object(builder, "files"));
    GtkCellRenderer *renderer;
    renderer = gtk_cell_renderer_text_new();
    gtk_tree_view_insert_column_with_attributes(
        GTK_TREE_VIEW(files), -1, "Bestand", renderer, "text", 0, NULL);
    GtkTreeModel *model = create_and_fill_model();
    gtk_tree_view_set_model(GTK_TREE_VIEW(files), model);
    g_object_unref(model);
    g_signal_connect(G_OBJECT(files), "row-activated",
                     G_CALLBACK(tree_row_activated_cb), NULL);
  }

  box = GTK_WIDGET(gtk_builder_get_object(builder, "my-box"));
  settings = webkit_settings_new();
  webkit_settings_set_default_font_size(settings, 18);
  webkit_settings_set_default_monospace_font_size(settings, 14);
  webkit_settings_set_default_charset(settings, "utf-8");
  webkit_settings_set_default_font_family(settings, "serif");
  webview = WEBKIT_WEB_VIEW(webkit_web_view_new_with_settings(settings));
  gtk_box_pack_start(GTK_BOX(box), GTK_WIDGET(webview), TRUE, TRUE, 0);
  webkit_web_view_load_uri(webview, url);

  g_signal_connect(window, "key-press-event", G_CALLBACK(web_view_key_pressed),
                   NULL);

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
