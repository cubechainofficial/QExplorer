package main

import (
	"strings"
	"path/filepath"
	"strconv"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type BlockInfoModel struct {
	walk.SortedReflectTableModelBase
	dirPath string
	items   []*Block
}

type BlockNumber struct {
	BlockIndex		int
	TotalIndex		string
	BlockBlockIndex	int
}


func NewBlockInfoModel() *BlockInfoModel {
	return new(BlockInfoModel)
}

func (m *BlockInfoModel) Items() interface{} {
	return m.items
}

func (m *BlockInfoModel) BlockSetting(dirPath string,index int,cubeno int) (int,int) {
	var oblock Block
	m.dirPath = dirPath
	m.items = nil
	path:=""
	if filepath.Ext(dirPath)==".blc" {
		path=dirPath
		index,cubeno=BlockPathIndex(path)
	} else {
		dir:=CubePath(index)
		filename:=fileSearch(dir,strconv.Itoa(index)+"_"+strconv.Itoa(cubeno-1)+"_")
		path=dir+filepathSeparator+filename
	}

	err:=pathRead(path,&oblock)
	if err!=nil {
		logger.Print(err)
	}

	item := &Block{
		Index:  oblock.Index,
		Cubeno: oblock.Cubeno+1,
		Timestamp:  oblock.Timestamp,
		Data : []byte(oblock.Data),
		Hash : oblock.Hash,			
		PrevHash : oblock.PrevHash,			
		PrevCubeHash : oblock.PrevCubeHash,			
		Nonce:  oblock.Nonce,
	}		
	m.items = append(m.items, item)

	m.PublishRowsReset()
	return index,cubeno
}

var BlocktableModel *BlockInfoModel
var BlockPathtext *walk.LineEdit
var BlockSearchtext *walk.LineEdit
var BlockCopyStr,BlockCopyStr1,BlockCopyStr2,BlockIndexing string
var BlockTableView *walk.TableView
var BlockTe *walk.TextEdit
var BlockLineEdit1 *walk.LineEdit
var BlockCombo1,BlockCombo2 *walk.ComboBox
var BlockTotalModel []string
var blockNumber BlockNumber
var blockNumber1 int 
var blockNumber2 int

func newBlockPage(parent walk.Container) (Page, error) {
	var splitter *walk.Splitter
	var db *walk.DataBinder
	var BlockLabel1,BlockLabel2,BlockLabel3 *walk.Label

	BlocktableModel = NewBlockInfoModel()
	BlocktableModel.BlockSetting(blockpath,1,1)


	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "Block",
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
										AssignTo: &BlockPathtext,
										MaxSize: Size{1000, 22},
										Text: blockpath,
									},
									PushButton{
										Text: "Select",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											blockpath=BlockPathtext.Text()
											idx,cno:=BlockPathSetting()
											BlockPathtext.SetText(blockpath)
											BlockCombo1.SetText(strconv.Itoa(idx))
											BlockCombo2.SetText(strconv.Itoa(cno))
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								DataBinder: DataBinder{
									AssignTo: &db,
									DataSource:     blockNumber,
									ErrorPresenter: ToolTipErrorPresenter{},
								},
								Children: []Widget{
									Label{
										AssignTo: &BlockLabel1,
										Text: "Total Block : ",
									},
									LineEdit{
										AssignTo: &BlockLineEdit1,
										MaxSize: Size{1000, 22},
										Text: Bind("TotalIndex"),
									},
									Label{
										AssignTo: &BlockLabel2,
										Text: "Block Index : ",
									},
									ComboBox{
										AssignTo: &BlockCombo1,
										Editable: true,
										Value:    Bind("BlockIndex"),
										Model:    []string{"1","2","3","4","5","6","7","8","9","10"},
										OnCurrentIndexChanged: func() { 
											if BlockCombo1.CurrentIndex()>-1 {
												if BlockCombo2.CurrentIndex()<0 {
													BlockCombo2.SetText("1")		
												}
												BlockChangeIndex()
											}
										},
									},
									Label{
										AssignTo: &BlockLabel3,
										Text: "Block Number",
									},
									ComboBox{
										AssignTo: &BlockCombo2,
										Editable: true,
										Value:    Bind("BlockBlockIndex"),
										
										Model:    []string{"1","2","3","4","5","6","7","8","9","10","11","12","13","14","15","16","17","18","19","20","21","22","23","24","25","26","27"},
										OnCurrentIndexChanged: func() {
											if BlockCombo1.CurrentIndex()>-1 && BlockCombo2.CurrentIndex()>-1 {
												BlockChangeIndex()
											}										
										},
									},
								},
							},	
							TableView{
								AssignTo: &BlockTableView,
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
								Model: BlocktableModel,
								OnCurrentIndexChanged: func() {
									if index := BlockTableView.CurrentIndex(); index > -1 {
										blockNumber.BlockIndex=BlocktableModel.items[index].Index
										blockNumber.BlockBlockIndex=BlocktableModel.items[index].Cubeno
										BlockCombo1.SetText(strconv.Itoa(blockNumber.BlockIndex))
										BlockCombo2.SetText(strconv.Itoa(blockNumber.BlockBlockIndex))

										BlockIndexing="T:"+strconv.Itoa(BlocktableModel.items[index].Index)
										BlockCopyStr=BlocktableModel.items[index].Hash
										BlockCopyStr1=BlocktableModel.items[index].PrevHash
										BlockCopyStr2=BlocktableModel.items[index].PrevCubeHash
										BlockTe.SetText(ByteToStr(BlocktableModel.items[index].Data))
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
											walk.Clipboard().SetText(BlockCopyStr)
										},
									},
									PushButton{
										Text:  "Prev Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(BlockCopyStr1)
										},
									},
									PushButton{
										Text:  "Pattern Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(BlockCopyStr2)
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
										AssignTo: &BlockSearchtext,
										MaxSize: Size{1000, 22},
										Text: "",
									},
									PushButton{
										Text:  "Search",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											BlockSearch(BlockSearchtext.Text())
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										MinSize: Size{0, 22},
										Text: "Decoding Block Data : ",
									},
								},
							},	
							TextEdit{
								AssignTo: &BlockTe,
							},		
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Data Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(BlockTe.Text())
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
	if err := BlockTe.SetReadOnly(true); err != nil {
		return nil, err
	}	
	if err := BlockLineEdit1.SetReadOnly(true); err != nil {
		return nil, err
	}	
	
	BlockTotalModel=BlockTotalSetting()
	BlockCombo1.SetModel(BlockTotalModel)

	return p, nil
}


func BlockTotalSetting() []string{
	var result []string
	path:=CubeBasePath(BlockPathtext.Text())
	cb:=CubeFileCount(path,".blc")
	BlockLineEdit1.SetText(strconv.Itoa(cb))
	c:=CurrentHeight(path)-1
	for i:=0;i<c;i++ {
		result=append(result,strconv.Itoa(i+1))
	}
	return result
}

func BlockPathSetting() (int,int) {
	var dlgForm walk.Form
	dlg := new(walk.FileDialog)

	blockpath=strings.Replace(blockpath,filepathSeparator+filepathSeparator,filepathSeparator, -1)
	dlg.FilePath = blockpath
	dlg.Filter = "Block Files (*.blc)|*.blc"
	dlg.Title = "Select an block files"
    

	if ok, err := dlg.ShowOpen(dlgForm); err != nil {
        logger.Print("[Dialog Error:]")
        logger.Print(err)
		return 0,0
	} else if !ok {
		return 0,0
	}
	blockpath=dlg.FilePath


	i,c:=BlocktableModel.BlockSetting(blockpath,1,1)

	return i,c
}

func BlockSearch(search string) {
	for k,v:=range BlocktableModel.items {
		if v.Hash==search {
			BlockTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block Hash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		} else if v.PrevHash==search {
			BlockTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block PevHash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		} else if v.PrevCubeHash==search {
			BlockTableView.SetCurrentIndex(k)
			walk.MsgBox(mw, "Search", "Search Block PevCubeHash ["+strconv.Itoa(v.Index)+"_"+strconv.Itoa(v.Cubeno)+"]", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			return
		}
	}
	walk.MsgBox(mw, "Search", "Not found hash!", walk.MsgBoxOK|walk.MsgBoxIconWarning)
	return
}

func BlockChangeIndex() (int,int) {
	dir:=""
	index,_:=strconv.Atoi(BlockCombo1.Text())
	cubeno,_:=strconv.Atoi(BlockCombo2.Text())
	BlocktableModel.BlockSetting(dir,index,cubeno)
	BlockTe.SetText("No found data")
	BlockTe.SetText(ByteToStr(BlocktableModel.items[0].Data))
	blockpath=BlockPath(index,cubeno)
	BlockPathtext.SetText(blockpath)
	return index,cubeno
}

func BlockChangeCombo(index int,cubeno int) {
	blockNumber.BlockIndex=index
	blockNumber.BlockBlockIndex=cubeno
	return
}

