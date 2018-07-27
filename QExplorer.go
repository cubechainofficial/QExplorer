package main

import (
	"bytes"
	"time"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mw AppMainWindow
func main() {
	walk.Resources.SetRootDirPath("./img")

	var openAction *walk.Action
	var recentMenu *walk.Menu

	app := walk.App()
	app.SetOrganizationName("Cubesystem")
	app.SetProductName("Cubechain Explorer")
	settings := walk.NewIniFileSettings("settings.ini")
	settings.SetExpireDuration(time.Hour * 24 * 30 * 1000)
	if err := settings.Load(); err != nil {
		logger.Fatal(err)
	}
	
	app.SetSettings(settings)
	
	cfg := &MultiPageMainWindowConfig{
		Name:    "Cubechain Explorer",
		MinSize: Size{800, 600},
		Size:     Size{1248, 800},
		
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Open",
						Image:       "open.png",
						Enabled:     Bind("enabledCB.Checked"),
						Visible:     Bind("!openHiddenCB.Checked"),
						Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: mw.openAction_Triggered,
					},
					Menu{
						AssignTo: &recentMenu,
						Text:     "Recent",
					},
					Separator{},
					Action{
						Text:        "About",
						OnTriggered: func() { mw.aboutAction_Triggered() },
					},
					Action{
						Text:        "E&xit",
						Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Cube",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "&Cube open",
						Image:       "open.png",
						Enabled:     Bind("enabledCB.Checked"),
						Visible:     Bind("!openHiddenCB.Checked"),
						Shortcut:    Shortcut{walk.ModControl, walk.KeyC},
						OnTriggered: mw.openAction_Triggered,
					},
					Menu{
						AssignTo: &recentMenu,
						Text:     "Recent",
					},
					Separator{},
					Action{
						Text:        "About",
						OnTriggered: func() { mw.aboutAction_Triggered() },
					},
					Action{
						Text:        "E&xit",
						Shortcut:    Shortcut{walk.ModControl, walk.KeyX},
						OnTriggered: func() { mw.Close() },
					},
				},
			},
		},
		OnCurrentPageChanged: func() {
			mw.updateTitle(mw.CurrentPageTitle())
		},
		PageCfgs: []PageConfig{
			{"File", "document-new.png", newFilePage},
			{"Chain", "document-new.png", newChainPage},
			{"Cube", "document-new.png", newCubePage},
			{"Block", "document-properties.png", newBlockPage},
			{"DataCenter", "system-shutdown.png", newCenterPage},
			{"Setting", "document-properties.png", newSettingPage},
		},

	}
	mpmw, err := NewMultiPageMainWindow(cfg)
	if err != nil {
		panic(err)
	}
	mw.MultiPageMainWindow = mpmw
	mw.updateTitle(mw.CurrentPageTitle())

	mw.SetFullscreen(true)
	mw.SetFullscreen(false)
	mw.Run()

	if err := settings.Save(); err != nil {
		logger.Print("Setting file not save")
		logger.Fatal(err)
	}
}

type AppMainWindow struct {
	*MultiPageMainWindow
}

func (mw *AppMainWindow) updateTitle(prefix string) {
	var buf bytes.Buffer
	buf.WriteString("Cubechain Explorer")
	if prefix != "" {
		buf.WriteString(" - ")
		buf.WriteString(prefix)
	}
	mw.SetTitle(buf.String())
}

func (mw *AppMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
}

func (mw *AppMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw,
		"About Walk Multiple Pages Example",
		"An example that demonstrates a main window that supports multiple pages.",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

type basePage struct {
	*walk.Composite
}
