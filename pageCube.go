package main

import (
	"path/filepath"
	"strconv"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type CubeInfoModel struct {
	walk.SortedReflectTableModelBase
	dirPath string
	items   []*Block
}

type CubeNumber struct {
	CubeIndex		int
	TotalIndex		string
	CubeBlockIndex	int
}

func NewCubeInfoModel() *CubeInfoModel {
	return new(CubeInfoModel)
}

func (m *CubeInfoModel) Items() interface{} {
	return m.items
}

func (m *CubeInfoModel) CubeSetting(dirPath string) error {
	var cubeblock CubeBlock
	m.dirPath = dirPath
	m.items = nil
	path:=""

	if filepath.Ext(dirPath)==".cub" {
		path=dirPath
	} else {
		dir:=dirPath
		filename:=fileSearch(dir,".cub")
		path=dir+filepathSeparator+filename
	}

	err:=pathRead(path,&cubeblock)
	if err!=nil {
		logger.Print(err)
	}
	for _,v:= range cubeblock.Cube {
		item := &Block{
			Index:  v.Index,
			Cubeno: v.Cubeno+1,
			Timestamp:  v.Timestamp,
			Data : []byte(v.Data),
			Hash : v.Hash,			
			PrevHash : v.PrevHash,			
			PrevCubeHash : v.PrevCubeHash,			
			Nonce:  v.Nonce,
		}		
		m.items = append(m.items, item)
	}
	m.PublishRowsReset()
	return nil
}

var CubetableModel *CubeInfoModel
var CubePathtext *walk.LineEdit
var CubeSearchtext *walk.LineEdit
var CubeCopyStr,CubeCopyStr1,CubeCopyStr2,CubeIndexing string
var CubeTableView *walk.TableView
var CubeTe *walk.TextEdit
var CubeLineEdit1 *walk.LineEdit
var CubeCombo1,CubeCombo2 *walk.ComboBox
var CubeTotalModel []string

func newCubePage(parent walk.Container) (Page, error) {
	var splitter *walk.Splitter
	var db *walk.DataBinder
	var CubeLabel1,CubeLabel2,CubeLabel3 *walk.Label
	var cubeNumber=CubeNumber{0,"2000",0}

	CubetableModel = NewCubeInfoModel()
	CubetableModel.CubeSetting(cubepath)


	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "Cube",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				AssignTo: &splitter,
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "BasePath:",
										MaxSize: Size{80, 22},
									},
									LineEdit{
										AssignTo: &CubePathtext,
										MaxSize: Size{1000, 22},
										Text: cubepath,
									},
									PushButton{
										Text: "Select",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											cubepath=CubePathtext.Text()
											CubePathSetting()
											CubePathtext.SetText(cubepath)
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								DataBinder: DataBinder{
									AssignTo: &db,
									DataSource:     cubeNumber,
									ErrorPresenter: ToolTipErrorPresenter{},
								},
								Children: []Widget{
									Label{
										AssignTo: &CubeLabel1,
										Text: "Total Cube : ",
									},
									LineEdit{
										AssignTo: &CubeLineEdit1,
										MaxSize: Size{1000, 22},
										Text: Bind("TotalIndex"),
									},
									Label{
										AssignTo: &CubeLabel2,
										Text: "Cube Index : ",
									},
									ComboBox{
										AssignTo: &CubeCombo1,
										Editable: true,
										Value:    Bind("CubeIndex"),
										Model:    []string{"1","2","3","4","5","6","7","8","9","10"},
										OnCurrentIndexChanged: func() { 
											if CubeCombo1.CurrentIndex() > -1  {
												if CubeCombo2.CurrentIndex() < 0  {
													CubeCombo2.SetText("1")
												}
												c,_:=strconv.Atoi(CubeCombo1.Text())
												CubeChangeIndex(c)
												c2,_:=strconv.Atoi(CubeCombo2.Text())
												CubeChangeBlock(c2)
											}
										},
									},
									Label{
										AssignTo: &CubeLabel3,
										Text: "Cube Number",
									},
									ComboBox{
										AssignTo: &CubeCombo2,
										Editable: true,
										Value:    Bind("CubeBlockIndex"),
										Model:    []string{"1","2","3","4","5","6","7","8","9","10","11","12","13","14","15","16","17","18","19","20","21","22","23","24","25","26","27"},
										OnCurrentIndexChanged: func() {
											if CubeCombo1.CurrentIndex() > -1 && CubeCombo2.CurrentIndex() > -1 {
												c,_:=strconv.Atoi(CubeCombo2.Text())
												CubeChangeBlock(c)
											}										
										},
									},
								},
							},	
							TableView{
								AssignTo: &CubeTableView,
								StretchFactor: 2,
								AlternatingRowBGColor: walk.RGB(239, 239, 239),
								Columns: []TableViewColumn{
									{DataMember: "Index",Format:"%d",Width:60,},
									{DataMember: "Cubeno",Format:"%d",Width:60,},
									{DataMember: "Timestamp",Format:"%d",Width:90,},
									{DataMember: "Hash"}, 
									{DataMember: "PrevHash"},
									{DataMember: "PrevCubeHash"},
									{DataMember: "Nonce",Format:"%d"},
								},
								Model: CubetableModel,
								OnCurrentIndexChanged: func() {
									if index := CubeTableView.CurrentIndex(); index > -1 {
										cubeNumber.CubeIndex=CubetableModel.items[index].Index
										cubeNumber.CubeBlockIndex=CubetableModel.items[index].Cubeno
										CubeCombo1.SetText(strconv.Itoa(cubeNumber.CubeIndex))
										CubeCombo2.SetText(strconv.Itoa(cubeNumber.CubeBlockIndex))

										CubeIndexing="T:"+strconv.Itoa(CubetableModel.items[index].Index)
										CubeCopyStr=CubetableModel.items[index].Hash
										CubeCopyStr1=CubetableModel.items[index].PrevHash
										CubeCopyStr2=CubetableModel.items[index].PrevCubeHash
										CubeTe.SetText(ByteToStr(CubetableModel.items[index].Data))
									}							
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CubeCopyStr)
										},
									},
									PushButton{
										Text:  "Prev Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CubeCopyStr1)
										},
									},
									PushButton{
										Text:  "Pattern Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CubeCopyStr2)
										},
									},
								},
							},	

						},
					},	
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "Search:",
										MaxSize: Size{80, 22},
									},
									LineEdit{
										AssignTo: &CubeSearchtext,
										MaxSize: Size{1000, 22},
										Text: "",
									},
									PushButton{
										Text:  "Search",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											CubeSearch(CubeSearchtext.Text())
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								DataBinder: DataBinder{
									AssignTo: &db,
									DataSource:     cubeNumber,
									ErrorPresenter: ToolTipErrorPresenter{},
								},
								Children: []Widget{
									Label{
										MinSize: Size{0, 22},
										Text: "Decoding Block Data : ",
									},
								},
							},	
							TextEdit{
								AssignTo: &CubeTe,
							},		
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Data Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CubeTe.Text())
										},
									},
								},
							},								
						},
					},	
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
	if err := CubeTe.SetReadOnly(true); err != nil {
		return nil, err
	}	
	if err := CubeLineEdit1.SetReadOnly(true); err != nil {
		return nil, err
	}	
	
	CubeTotalModel=CubeTotalSetting()
	CubeCombo1.SetModel(CubeTotalModel)
	return p, nil
}

func CubeTotalSetting() []string{
	var result []string
	path:=CubeBasePath(CubeSearchtext.Text())
	c:=CurrentHeight(path)-1
	CubeLineEdit1.SetText(strconv.Itoa(c))
	for i:=0;i<c;i++ {
		result=append(result,strconv.Itoa(i+1))
	}
	return result
}

func CubePathSetting() error {
	var dlgForm walk.Form
	dlg := new(walk.FileDialog)
	dlg.FilePath = cubepath
	dlg.Filter = "Cube Files (*.cub)|*.cub"
	dlg.Title = "Select an cube files"
	if ok, err := dlg.ShowOpen(dlgForm); err != nil {
		return err
	} else if !ok {
		return nil
	}
	cubepath=dlg.FilePath
	CubetableModel.CubeSetting(cubepath)
	i:=CubePathNum(cubepath)
	CubeCombo1.SetText(strconv.Itoa(i))
	CubeCombo2.SetText("1")
	CubeChangeBlock(1)
	return nil
}

func CubeSearch(search string) {
	for k,v:=range CubetableModel.items {
		if v.Hash==search {
			CubeTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block Hash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		} else if v.PrevHash==search {
			CubeTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block PevHash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		} else if v.PrevCubeHash==search {
			CubeTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block PevCubeHash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		}
	}
	walk.MsgBox(mw, "Search", "Not found hash!", walk.MsgBoxOK|walk.MsgBoxIconWarning)
	return
}

func CubeChangeIndex(index int) {
	dir:=CubePath(index)
	filename:=fileSearch(dir,".cub")
	path:=dir+filepathSeparator+filename
	CubetableModel.CubeSetting(path)
	cubepath=path
	CubePathtext.SetText(cubepath)
	return
}

func CubeChangeBlock(index int) {
	for _,v:=range CubetableModel.items {
		if v.Cubeno==index {
			CubeTe.SetText(ByteToStr(v.Data))
			return
		}
	}
	CubeTe.SetText("No found data")
	return
}
