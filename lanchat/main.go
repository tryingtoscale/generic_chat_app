package main

import (
	"bufio"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"fmt"
	"log"
	"net"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// https://pkg.go.dev/github.com/gotk3/gotk3/gtk?utm_source=godoc
// https://docs.gtk.org/gtk3

const gladeTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!-- Generated with glade 3.22.2 -->
<interface>
  <requires lib="gtk+" version="3.20"/>
  <object class="GtkWindow" id="myPopup">
    <property name="name">LanChatApp</property>
    <property name="width_request">600</property>
    <property name="height_request">600</property>
    <property name="visible">True</property>
    <property name="can_focus">True</property>
    <property name="has_focus">True</property>
    <property name="is_focus">True</property>
    <property name="can_default">True</property>
    <property name="has_default">True</property>
    <property name="receives_default">True</property>
    <property name="halign">baseline</property>
    <property name="valign">baseline</property>
    <property name="border_width">2</property>
    <property name="title" translatable="yes">Lan Chat</property>
    <property name="window_position">center</property>
    <property name="destroy_with_parent">True</property>
    <property name="icon_name">network-workgroup</property>
    <property name="deletable">False</property>
    <property name="gravity">north</property>
    <property name="has_resize_grip">True</property>
    <signal name="window-state-event" handler="myMainEvent" swapped="no"/>
    <child type="titlebar">
      <placeholder/>
    </child>
    <child>
      <object class="GtkGrid" id="myGrid">
        <property name="visible">True</property>
        <property name="can_focus">False</property>
        <property name="hexpand">True</property>
        <property name="vexpand">True</property>
        <child>
          <object class="GtkEntry" id="myEntry">
            <property name="width_request">600</property>
            <property name="height_request">80</property>
            <property name="visible">True</property>
            <property name="can_focus">True</property>
            <property name="can_default">True</property>
            <property name="has_default">True</property>
            <property name="valign">baseline</property>
            <property name="hexpand">True</property>
            <property name="max_length">250</property>
            <property name="activates_default">True</property>
            <property name="width_chars">30</property>
            <property name="max_width_chars">40</property>
            <property name="secondary_icon_name">face-smile-symbolic</property>
            <property name="secondary_icon_tooltip_text" translatable="yes">Insert Emoji</property>
            <property name="secondary_icon_tooltip_markup" translatable="yes">Insert Emoji</property>
            <property name="placeholder_text" translatable="yes">Enter your Message Here:</property>
            <property name="show_emoji_icon">True</property>
            <property name="enable_emoji_completion">True</property>
            <signal name="key-press-event" handler="key_pressed" swapped="no"/>
          </object>
          <packing>
            <property name="left_attach">0</property>
            <property name="top_attach">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkScrolledWindow" id="myScroll">
            <property name="width_request">600</property>
            <property name="height_request">400</property>
            <property name="visible">True</property>
            <property name="can_focus">True</property>
            <property name="valign">baseline</property>
            <property name="hexpand">True</property>
            <property name="vexpand">True</property>
            <property name="vscrollbar_policy">always</property>
            <property name="shadow_type">in</property>
            <property name="min_content_width">600</property>
            <property name="min_content_height">400</property>
            <property name="max_content_width">600</property>
            <property name="max_content_height">400</property>
            <property name="propagate_natural_width">True</property>
            <property name="propagate_natural_height">True</property>
            <child>
              <object class="GtkTreeView" id="myText">
                <property name="width_request">600</property>
                <property name="height_request">400</property>
                <property name="visible">True</property>
                <property name="can_focus">True</property>
                <property name="valign">baseline</property>
                <property name="hexpand">True</property>
                <property name="vexpand">True</property>
                <property name="enable_grid_lines">both</property>
                <signal name="size-allocate" handler="treeview_changed" swapped="no"/>
                <child internal-child="selection">
                  <object class="GtkTreeSelection"/>
                </child>
              </object>
            </child>
          </object>
          <packing>
            <property name="left_attach">0</property>
            <property name="top_attach">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkFixed" id="myFixed">
            <property name="name">wFixed</property>
            <property name="width_request">600</property>
            <property name="height_request">32</property>
            <property name="visible">True</property>
            <property name="can_focus">False</property>
            <property name="valign">baseline</property>
            <child>
              <object class="GtkButton" id="mySend">
                <property name="label" translatable="yes">Send</property>
                <property name="width_request">100</property>
                <property name="height_request">32</property>
                <property name="visible">True</property>
                <property name="can_focus">True</property>
                <property name="receives_default">True</property>
                <property name="valign">baseline</property>
                <property name="resize_mode">immediate</property>
                <signal name="clicked" handler="btnSend" swapped="no"/>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="myMin">
                <property name="label" translatable="yes">Minimize</property>
                <property name="width_request">100</property>
                <property name="height_request">32</property>
                <property name="visible">True</property>
                <property name="can_focus">True</property>
                <property name="receives_default">True</property>
                <property name="valign">baseline</property>
                <property name="image_position">right</property>
                <signal name="clicked" handler="btnMinimize" swapped="no"/>
              </object>
              <packing>
                <property name="x">500</property>
              </packing>
            </child>
            <child>
              <object class="GtkButton" id="myConnect">
                <property name="label" translatable="yes">Connect</property>
                <property name="width_request">100</property>
                <property name="height_request">32</property>
                <property name="visible">True</property>
                <property name="can_focus">True</property>
                <property name="receives_default">True</property>
                <property name="valign">baseline</property>
                <signal name="clicked" handler="btnConnect" swapped="no"/>
              </object>
              <packing>
                <property name="x">263</property>
              </packing>
            </child>
          </object>
          <packing>
            <property name="left_attach">0</property>
            <property name="top_attach">2</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>`

func show_debug_info(info string) {
	const myDebuggingMode bool = false
	if myDebuggingMode {
		fmt.Println(info)
	}
}

type GtkUserInterface struct {
	window    *gtk.Window
	builder   *gtk.Builder
	myList    *gtk.ListStore
	mytxtBox  *gtk.Entry
	myConnect *gtk.Button
}

// you just place them in a map that names the signals, then feed the map to the builder
var signals = map[string]interface{}{
	"btnConnect":       btnConnect,
	"myMainEvent":      updatedEvent,
	"key_pressed":      txtPressed,
	"btnSend":          btnSend,
	"btnMinimize":      btnMin,
	"treeview_changed": updatedTree,
}

var userInterface *GtkUserInterface

// IDs to access the tree view columns by
const (
	COLUMN_USERNAME = iota
	COLUMN_MSG
)

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

// Creates a tree view and the list store that holds its data
func setupTreeView(treeView *gtk.TreeView) *gtk.ListStore {

	treeView.AppendColumn(createColumn("User", COLUMN_USERNAME))
	treeView.AppendColumn(createColumn("Message", COLUMN_MSG))

	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)

	return listStore
}

// Append a row to the list store for the tree view
func addRow(listStore *gtk.ListStore, username string, message string) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_USERNAME, COLUMN_MSG},
		[]interface{}{username, message})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}

func updatedTree() {
	scroll_obj, scroll_bad := userInterface.builder.GetObject("myScroll")
	if scroll_bad != nil {
		log.Fatalln("Couldn't get myScroll")
	}
	myScroll := scroll_obj.(*gtk.ScrolledWindow)

	adj := myScroll.GetVAdjustment()
	myNewSize := adj.GetUpper() - adj.GetPageSize()
	adj.SetValue(myNewSize)
}

// looks like handlers can literally be any function or method

func updatedEvent(w *gtk.Window, ev *gdk.Event) bool {
	var state_of_win gdk.WindowState

	wev := gdk.EventWindowStateNewFromEvent(ev)
	state_of_win = wev.ChangedMask()

	switch state_of_win {
	case gdk.WINDOW_STATE_FOCUSED:
		show_debug_info("focused")
	case gdk.WINDOW_STATE_ICONIFIED:
		show_debug_info("Icon")
	case gdk.WINDOW_STATE_MAXIMIZED:
		show_debug_info("maximized")
	default:
		show_debug_info("Other")
	}

	return true
}

func txtPressed(win *gtk.Entry, event *gdk.Event) {
	keyEvent := gdk.EventKeyNewFromEvent(event)
	keyVal := keyEvent.KeyVal()

	const mainEnterKey = 65293
	const keyPadEnterKey = 65421
	// Enter key will Send Text
	if keyVal == mainEnterKey || keyVal == keyPadEnterKey {
		btnSend()
	}
}

func btnConnect() {
	Do_connect()
}

var lastMsg string // Track this message sent, to make sure it does not unIcon by mistake.
func btnSend() {
	data := userInterface.mytxtBox
	s, e := data.GetText()
	if e != nil {
		log.Fatalln("Couldn't get Text:", e)
	}
	isEmptyCheck := strings.TrimSpace(s)
	if isEmptyCheck == "" {
		return
	}
	if is_connected {
		i, wErr := chat_connection.Write([]byte(s + "\n"))
		if i > 0 && wErr == nil {
			lastMsg = s
			data.SetText("") // Empty Chat Box as message was sent, successfully?
		} else {
			is_connected = false
			btn := userInterface.myConnect
			if btn != nil {
				set_btn_visiable(btn, true)
			}
			show_debug_info("Lost connection.")
		}
	}
}

func btnMin() {
	userInterface.window.Iconify()
}

func readConf(filename string) (*ChatData, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &ChatData{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

type ChatData struct {
	Conf struct {
		ConnHost string `yaml:"host"`
		ConnPort string `yaml:"port"`
		ConnType string `yaml:"type"`
	}
}

var is_connected bool = false
var chat_connection net.Conn

func set_btn_visiable(btnCnt *gtk.Button, visiable bool) {
	btnCnt.SetVisible(visiable)
}

func sayStatus(online bool) {
	var s string
	currentTime := time.Now().Format(time.Kitchen)
	if online {
		s = "I'm Online at : " + currentTime
	} else {
		s = "I'm Offline at : " + currentTime
	}

	chat_connection.Write([]byte(s + "\n"))
	lastMsg = s
}

func Do_connect() bool {
	var err error
	if !is_connected {
		chat_connection, err = net.Dial(SocketConn.Conf.ConnType, SocketConn.Conf.ConnHost+":"+SocketConn.Conf.ConnPort)
		// Note: Do not defer closing connection as it will bail out right away!
		btn := userInterface.myConnect
		if err != nil {
			show_debug_info("Error connecting:" + err.Error())
			if btn != nil {
				set_btn_visiable(btn, true)
			}
			return false
		} else {
			sayStatus(true) // Say online @ time + Fixes issue with MSGs not sending out right...
			show_debug_info("Connected.")
			is_connected = true
			if btn != nil {
				set_btn_visiable(btn, false)
			}
			return true
		}
	}
	return true
}

var stopDoubleMsgs string

func check_chat() {
	const debugging_check_chat bool = false
	if is_connected {
		go func() {
			var opps error
			var message string
			if debugging_check_chat {
				message = "Me:Testing"
			} else {
				message, opps = bufio.NewReader(chat_connection).ReadString('\n')
			}

			if opps != nil {
				is_connected = false
				btn := userInterface.myConnect
				if btn != nil {
					set_btn_visiable(btn, true)
				}
				show_debug_info("Lost Connection.")
			}

			tmsg := strings.TrimSpace(message)
			if tmsg != "" && tmsg != stopDoubleMsgs {
				if strings.Contains(tmsg, ":") { // Valid Message after Here
					realMsg := strings.SplitN(tmsg, ":", 2)
					addRow(userInterface.myList, realMsg[0], realMsg[1])
					stopDoubleMsgs = tmsg
					if lastMsg != realMsg[1] {
						userInterface.window.Deiconify() // UnMinimize -- nope, doesn't always do that...
						userInterface.window.Present()   // Popup on New MSG!
					}
				}
			}
		}()
	}
}

var SocketConn *ChatData
var nSets = 1

func main() {
	var app *gtk.Application
	var cerr error
	SocketConn, cerr = readConf("conf.yaml")
	if cerr != nil {
		log.Fatal(cerr)
	}

	const appID = "robs.lanchat"
	const doing_debugging bool = false // If true loads Chat.glade, else uses Template in here...

	var app_err error
	app, app_err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if app_err != nil {
		log.Fatalln("Couldn't create app:", app_err)
	}

	// It looks like all builder code must execute in the context of `app`.
	// If you try creating the builder inside the main function instead of
	// the `app` "activate" callback, then you will get a segfault
	app.Connect("activate", func() {
		var builder *gtk.Builder
		var build_err error
		if doing_debugging {
			builder, build_err = gtk.BuilderNewFromFile("Chat.glade")
			if build_err != nil {
				log.Fatalln("Couldn't load glade file")
			}
		} else {
			builder, build_err = gtk.BuilderNew()
			if build_err != nil {
				log.Fatalln("Couldn't make new builder:", build_err)
			}

			template_err := builder.AddFromString(gladeTemplate)
			if template_err != nil {
				log.Fatalln("Couldn't add builder from string:", template_err)
			}
		}

		builder.ConnectSignals(signals)

		obj, err := builder.GetObject("myPopup")
		if err != nil {
			log.Fatalln("Couldn't get Win")
		}
		wnd := obj.(*gtk.Window)

		mobj, bad := builder.GetObject("myText")
		if bad != nil {
			log.Fatalln("Couldn't get TextView")
		}
		tvx := mobj.(*gtk.TreeView)

		listStore := setupTreeView(tvx)

		eobj, ebad := builder.GetObject("myEntry")
		if ebad != nil {
			log.Fatalln("Couldn't get Entry text Box")
		}
		myEntry := eobj.(*gtk.Entry)

		btnObj, btnBad := builder.GetObject("myConnect")
		if btnBad != nil {
			log.Fatalln("Couldn't get Connect Button")
		}
		myBtnConnect := btnObj.(*gtk.Button)

		userInterface = &GtkUserInterface{
			window:    wnd,
			builder:   builder,
			myList:    listStore,
			mytxtBox:  myEntry,
			myConnect: myBtnConnect,
		}

		wnd.ShowAll()
		app.AddWindow(wnd)

		Do_connect()
	})

	go func() {
		for {
			time.Sleep(time.Second)
			glib.IdleAdd(check_chat)
		}
	}()

	app.Run(os.Args)

	sayStatus(false) // Say offline @ time

	chat_connection.Close()
}
