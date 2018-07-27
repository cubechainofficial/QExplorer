package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DataSetting struct {
	BasePath	string
	ChainPath	string
	CubePath	string
	BlockPath	string
	DownPath	string
}

var SettingLineEdit1,SettingLineEdit2,SettingLineEdit3,SettingLineEdit4,SettingLineEdit5 *walk.LineEdit

func newSettingPage(parent walk.Container) (Page, error) {
	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "Setting",
		Layout:   VBox{SpacingZero:true,MarginsZero: true},
		Children: []Widget{
			GroupBox {
				Title: "Setting file in cubechain data path",
				Layout: Grid{Columns: 3},
				MaxSize: Size{700, 180},
				Children: []Widget{
					Label{ Text: "Base Path:",MaxSize: Size{100, 22},},
					LineEdit{ AssignTo: &SettingLineEdit1,Text: filesetpath,},
					PushButton{	Text:  "Select", OnClicked: func() { SettingPath(1) },},
					Label{ Text: "Chain Path:",MaxSize: Size{100, 22},},
					LineEdit{ AssignTo: &SettingLineEdit2,Text: chainpath,},
					PushButton{	Text:  "Select", OnClicked: func() { SettingPath(2) },},
					Label{ Text: "Cube Path:",MaxSize: Size{100, 22},},
					LineEdit{ AssignTo: &SettingLineEdit3,Text: cubepath,},
					PushButton{	Text:  "Select", OnClicked: func() { SettingPath(3) },},
					Label{ Text: "Block Path:",MaxSize: Size{100, 22},},
					LineEdit{ AssignTo: &SettingLineEdit4,Text: blockpath,},
					PushButton{	Text:  "Select", OnClicked: func() { SettingPath(4) },},
					Label{ Text: "Download Path:",MaxSize: Size{100, 22},},
					LineEdit{ AssignTo: &SettingLineEdit5,Text: downloadpath,},
					PushButton{	Text:  "Select", OnClicked: func() { SettingPath(5) },},
				},
			},
			Composite{
				Layout: Grid{Columns: 5},
				Children: []Widget{
					Label{ Text: "",MaxSize: Size{100, 22},},
					PushButton{	Text:  "Current File Load", OnClicked: func() { SettingLoad() },},
					Label{ Text: "",MaxSize: Size{100, 22},},
					PushButton{	Text:  "Edit Setting Save", OnClicked: func() { SettingSave() },},
					Label{ Text: "",MaxSize: Size{100, 22},},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
        logger.Print(err)
		return nil, err
	}
	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	return p, nil
}

func SettingPath(no int) error {
	path:=""
	title:=""
	filter:=""
	switch no {
		case 1: 
			path=filesetpath
			title="Base"
			filter=""
		case 2:
			path=chainpath
			title="Chain"
			filter = "Chain Files (*.chn)|*.chn"
		case 3:
			path=cubepath
			title="Cube"
			filter = "Cube Files (*.cub)|*.cub"
		case 4:
			path=blockpath
			title="Block"
			filter = "Block Files (*.blc)|*.blc"
		case 5:
			path=downloadpath
			title="Download"
			filter=""
		default: return nil
	}

	var dlgForm walk.Form
	dlg := new(walk.FileDialog)

	filesetpath=ArrangeSeparator(filesetpath)
	dlg.FilePath=path
	dlg.Title="Select an "+title+" path"
	if filter=="" {
		if ok, err := dlg.ShowBrowseFolder(dlgForm); err != nil {
			return err
		} else if !ok {
			return nil
		}
	} else {
		dlg.Filter = filter
		if ok, err := dlg.ShowOpen(dlgForm); err != nil {
			return err
		} else if !ok {
			return nil
		}	
	}

	switch no {
		case 1: SettingLineEdit1.SetText(dlg.FilePath)
		case 2: SettingLineEdit2.SetText(dlg.FilePath)
		case 3: SettingLineEdit3.SetText(dlg.FilePath)
		case 4: SettingLineEdit4.SetText(dlg.FilePath)
		case 5: SettingLineEdit5.SetText(dlg.FilePath)
	}
	return nil
}

func SettingLoad() {
	var datasetting DataSetting
	err:=pathRead(setfilePath(),&datasetting)
	if err!=nil {
		logger.Print(err)
	}
	filesetpath=datasetting.BasePath
	chainpath=datasetting.ChainPath
	cubepath=datasetting.CubePath
	blockpath=datasetting.BlockPath
	downloadpath=datasetting.DownPath
	SettingLineEdit1.SetText(filesetpath)
	SettingLineEdit2.SetText(chainpath)
	SettingLineEdit3.SetText(cubepath)
	SettingLineEdit4.SetText(blockpath)
	SettingLineEdit5.SetText(downloadpath)
	return 
}

func SettingSave() {
	datasetting:=&DataSetting{SettingLineEdit1.Text(),SettingLineEdit2.Text(),SettingLineEdit3.Text(),SettingLineEdit4.Text(),SettingLineEdit5.Text()}
	err:=pathWrite(setfilePath(),datasetting)
	if err!=nil {
		logger.Print(err)
	} else {
		SettingLoad()
	}
	return 
}
