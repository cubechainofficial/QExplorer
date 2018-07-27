package main

import (
	"log"
	"os"
	"strconv"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
    outfile, _ = os.Create("./chain.log")
    logger      = log.New(outfile, "", 0)
)

var ChainVerify="0"

type ChainItem struct {
	Block		int
	Hash1		string
	Hash2		string
}

type ChainInfoModel struct {
	walk.SortedReflectTableModelBase
	dirPath string
	verify	int
	items   []*CChain
}

type ChainModel struct {
	walk.SortedReflectTableModelBase
	itemtype string
	items   []*ChainItem
}

func NewChainInfoModel() *ChainInfoModel {
	return new(ChainInfoModel)
}

func (m *ChainInfoModel) Items() interface{} {
	return m.items
}

func NewChainModel() *ChainModel {
	return new(ChainModel)
}

func (m *ChainModel) Items() interface{} {
	return m.items
}

func (m *ChainModel) ItemSetting(itype string,index int) error {
	m.itemtype = itype
	m.items = nil
	var cubeitem Cubing
	if index>0 {
		cubeitem=CubingFileRead(index)
	} else {
		cubeitem=Cubing{}
	}

	for k,v:=range cubeitem.Hash1 {
		item := &ChainItem{
			Block:  k+1,
			Hash1:	v,
			Hash2:	cubeitem.Hash2[k],
		}	
		m.items = append(m.items, item)
	}
	m.PublishRowsReset()
	return nil
}

func ChainChange(index int)  {
	ChaintableModel2.ItemSetting("Chain_Ch",index)
}

func (m *ChainInfoModel) ChainSetting(dirPath string) error {
	m.dirPath = dirPath
	m.items = nil
	cubec:=ChainFileRead(dirPath)
	for _,v:= range cubec.Chain {
		item := &CChain{
			Index:  v.Index,
			Timestamp:  v.Timestamp,
			Chash:  v.Chash,
		}		
		m.items = append(m.items, item)
	}
	m.verify=cubec.Verify
	ChainVerify=strconv.Itoa(cubec.Verify)
	m.PublishRowsReset()
	return nil
}

var ChaintableModel *ChainInfoModel
var ChaintableModel2 *ChainModel
var ChainPathtext *walk.LineEdit
var ChainSearchtext *walk.LineEdit
var ChainCopyStr,ChainCopyStr1,ChainCopyStr2 string
var ChainTableView *walk.TableView
var ChainTableView2 *walk.TableView
var ChainLineEdit1,ChainLineEdit2 *walk.LineEdit
var ChainCombo1 *walk.ComboBox
var ChainTotalModel []string
var chainNumber BlockNumber
var ChainComposite *walk.Composite

func newChainPage(parent walk.Container) (Page, error) {
	var splitter *walk.Splitter
	var db *walk.DataBinder

	ChaintableModel = NewChainInfoModel()
	ChaintableModel.ChainSetting(chainpath)

	ChaintableModel2 = NewChainModel()
	ChaintableModel2.ItemSetting("Chain",0)

	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "Chain",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				AssignTo: &splitter,
				Children: []Widget{
					Composite{
						AssignTo: &ChainComposite,

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
										AssignTo: &ChainPathtext,
										MaxSize: Size{1000, 22},
										Text: chainpath,
									},
									PushButton{
										Text:  "Select",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											chainpath=ChainPathtext.Text()
											ChainPathSetting()
											ChainPathtext.SetText(chainpath)
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								DataBinder: DataBinder{
									AssignTo: &db,
									DataSource:     chainNumber,
									ErrorPresenter: ToolTipErrorPresenter{},
								},
								Children: []Widget{
									Label{
										Text: "Total Cube : ",
									},
									LineEdit{
										AssignTo: &ChainLineEdit1,
										MaxSize: Size{1000, 22},
										Text: Bind("TotalIndex"),
										ReadOnly:true,
									},
									Label{
										Text: "Cube Index : ",
									},
									ComboBox{
										AssignTo: &ChainCombo1,
										Editable: true,
										Value:    Bind("BlockIndex"),
										Model:    []string{"1","2","3","4","5","6","7","8","9","10"},
										OnCurrentIndexChanged: func() { 
											if ChainCombo1.CurrentIndex()>-1 {
												ChainChangeIndex()
											}
										},
									},
									Label{
										Text: "Verify Chain : ",
									},
									LineEdit{
										AssignTo: &ChainLineEdit2,
										MaxSize: Size{1000, 22},
										Text: ChainVerify,
										ReadOnly:true,
									},
								},
							},	
							TableView{
								AssignTo: &ChainTableView,
								StretchFactor: 2,
								AlternatingRowBGColor: walk.RGB(239, 239, 239),
								Columns: []TableViewColumn{
									{DataMember: "Index",Format:"%d",Width:60,},
									{DataMember: "Timestamp",Format:"%d",Width:90,},
									{DataMember: "Chash",Width:200,},
								},
								Model: ChaintableModel,
								OnCurrentIndexChanged: func() {
									if index := ChainTableView.CurrentIndex(); index > -1 {
										indexset := ChaintableModel.items[index].Index
										ChainChange(indexset)
										ChainCombo1.SetText(strconv.Itoa(indexset))
										ChainCopyStr=ChaintableModel.items[index].Chash
									}							
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Chash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(ChainCopyStr)
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
										AssignTo: &ChainSearchtext,
										MaxSize: Size{1000, 22},
										Text: "",
									},
									PushButton{
										Text:  "Search",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											ChainSearch(ChainSearchtext.Text())
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										MinSize: Size{0, 22},
										Text: "Block Hash Data : ",
									},
								},
							},	
							TableView{
								AssignTo:      &ChainTableView2,
								StretchFactor: 2,
								Columns: []TableViewColumn{
									TableViewColumn{
										DataMember: "Block",
										Format:     "%d",
										Width:      60,
									},
									TableViewColumn{
										DataMember: "Hash1",
										Width:      160,
									},
									TableViewColumn{
										DataMember: "Hash2",
										Width:      160,
									},
								},
								Model: ChaintableModel2,
								OnCurrentIndexChanged: func() {
									if index := ChainTableView2.CurrentIndex(); index > -1 {
										ChainCopyStr1=ChaintableModel2.items[index].Hash1
										ChainCopyStr2=ChaintableModel2.items[index].Hash2
									}
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Hash1 Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(ChainCopyStr1)
										},
									},
									PushButton{
										Text:  "Hash2 Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(ChainCopyStr2)
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
        logger.Print("Create Chain:")
        logger.Print(err)
		return nil, err
	}
	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	if err := ChainLineEdit1.SetReadOnly(true); err != nil {
		return nil, err
	}	

	ChainTotalModel=ChainTotalSetting()
	ChainCombo1.SetModel(ChainTotalModel)
	
	return p, nil
}

func ChainTotalSetting() []string{
	var result []string
	cb:=len(ChaintableModel.items)
	ChainLineEdit1.SetText(strconv.Itoa(cb))
	for i:=0;i<cb;i++ {
		result=append(result,strconv.Itoa(i+1))
	}
	return result
}

func ChainChangeIndex() {
	index,_:=strconv.Atoi(ChainCombo1.Text())
	for k,v:=range ChaintableModel.items {
		if v.Index==index {
			ChainTableView.SetCurrentIndex(k)
		}
	}
	return 
}


func ChainPathSetting() error {
	var dlgForm walk.Form
	dlg := new(walk.FileDialog)
	dlg.FilePath = chainpath
	dlg.Filter = "Chain Files (*.chn)|*.chn"
	dlg.Title = "Select an chain files"
	if ok, err := dlg.ShowOpen(dlgForm); err != nil {
		return err
	} else if !ok {
		return nil
	}
	chainpath=dlg.FilePath
	ChaintableModel.ChainSetting(chainpath)
	return nil
}

func ChainSearch(search string) {
	for k,v:=range ChaintableModel.items {
		if v.Chash==search {
			ChainTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Cube Hash ["+strconv.Itoa(v.Index)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		}
	}
	for k,v:=range ChaintableModel2.items {
		if v.Hash1==search {
			ChainTableView2.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block Hash1 ["+strconv.Itoa(v.Block)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		} else if v.Hash2==search {
			ChainTableView2.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block Hash2 ["+strconv.Itoa(v.Block)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		}
	}
	walk.MsgBox(mw, "Search", "Not found hash!", walk.MsgBoxOK|walk.MsgBoxIconWarning)
	return
}
