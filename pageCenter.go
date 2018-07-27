package main

import (
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DataCenter struct {
	No		int
	Type	string
	Index	int
	File	string
	FilePath string
}

type CenterInfoModel struct {
	walk.SortedReflectTableModelBase
	items   []*DataCenter
}

type CenterNumberT struct {
	CenterIndex		int
	TotalIndex		string
	CenterCenterIndex	int
}

func NewCenterInfoModel() *CenterInfoModel {
	return new(CenterInfoModel)
}

func (m *CenterInfoModel) Items() interface{} {
	return m.items
}

func (m *CenterInfoModel) CenterSetting() () {
	m.items = nil
	var file1,file2,file3 int
	if CenterChk1.CheckState() == walk.CheckChecked	{
		file1=1
	}	
	if CenterChk2.CheckState() == walk.CheckChecked	{
		file2=1
	}
	if CenterChk3.CheckState() == walk.CheckChecked	{
		file3=1
	}	
	reader :=strings.NewReader("file1="+strconv.Itoa(file1)+"&file2="+strconv.Itoa(file2)+"&file3="+strconv.Itoa(file3)+"&mode="+strconv.Itoa(CenterRbVal))
	request,_ := http.NewRequest("POST", "http://"+ datacenterHost +"/chain/datacenter.html", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res,_ := client.Do(request)
	body,_ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	CenterTe.SetText(ByteToStr(body))
	
	rdate:=strings.Split(ByteToStr(body),",")

	for k,v:=range rdate {
		filename:=PathToFile(v,"/")
		fileext:=FileToExt(filename)
		cidx:=0
		if fileext!="chn" {
			cidx=CubePathNum(strings.Replace(v,"/",filepathSeparator,-1))
		}
		item := &DataCenter{
			No:k+1,
			Type:fileext,
			Index:cidx,
			File:filename,
			FilePath:v,
		}
		m.items = append(m.items, item)		
	}

	m.PublishRowsReset()
	return 
}

var CentertableModel *CenterInfoModel
var CenterSearchtext *walk.LineEdit
var CenterCopyStr,CenterCopyStr1,CenterCopyStr2,CenterIndexing string
var CenterTableView *walk.TableView
var CenterTe *walk.TextEdit
var CenterLineEdit1 *walk.LineEdit
var CenterCombo1,CenterCombo2 *walk.ComboBox
var CenterTotalModel []string
var CenterNumber CenterNumberT
var CenterNumber1 int 
var CenterNumber2 int

var	CenterChk1   *walk.CheckBox
var	CenterChk2   *walk.CheckBox
var	CenterChk3   *walk.CheckBox

var	CenterRb1   *walk.RadioButton
var	CenterRb2   *walk.RadioButton
var CenterRbVal int

func newCenterPage(parent walk.Container) (Page, error) {
	var splitter *walk.Splitter
	CentertableModel = NewCenterInfoModel()

	CenterRbVal=1
	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "Center",
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
									GroupBox {
										Title: "",
										Layout: HBox{ MarginsZero: true},
										Children: []Widget {
											HSpacer { Size:8 },
											Label{
												Text: "Mode:",
											},
											RadioButtonGroup{
												Buttons: []RadioButton{
													RadioButton {
														AssignTo: &CenterRb1,
														Text: "ALL",
														Value: 1,
														OnClicked: func() {
															CenterRbVal=1
														},
													},
													RadioButton {
														AssignTo: &CenterRb2,
														Text: "Verify After Cube",
														Value: 2,
														OnClicked: func() {
															CenterRbVal=2
														},
													},
												},
											},
											HSpacer { Size:8 },
											VSpacer { Size:3 },
										},
									},
									GroupBox {
										Title: "",
										Layout: HBox{MarginsZero: true},
										Children: []Widget {
											HSpacer { Size:6 },
											Label{
												Text: "File Type:",
											},
											CheckBox {
												AssignTo: &CenterChk1,
												Text: "Chain File",
												Checked: true,
												OnCheckedChanged: func() {},
											},
											HSpacer { Size:2 },
											CheckBox {
												AssignTo: &CenterChk2,
												Text: "Cube File",
												Checked: true,
												OnCheckedChanged: func() {},
											},
											HSpacer { Size:2 },
											CheckBox {
												AssignTo: &CenterChk3,
												Text: "Block File",
												Checked: true,
												OnCheckedChanged: func() {},
											},
											HSpacer { Size:6 },
											VSpacer { Size:32 },
										},
									},
									PushButton{
										Text: "List Import",
										MinSize: Size{100, 40},
										OnClicked: func(){ 
											CentertableModel.CenterSetting()
										},
									},
								},
							},	

							TableView{
								AssignTo: &CenterTableView,
								StretchFactor: 2,
								AlternatingRowBGColor: walk.RGB(239, 239, 239),
								Columns: []TableViewColumn{
									{DataMember: "No",Format:"%d",Width:60,},
									{DataMember: "Type",Width:60,},
									{DataMember: "Index",Width:60,},
									{DataMember: "File"},
								},
								Model: CentertableModel,
								OnCurrentIndexChanged: func() {
									if index := CenterTableView.CurrentIndex(); index > -1 {
										
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
											walk.Clipboard().SetText(CenterCopyStr)
										},
									},
									PushButton{
										Text:  "Prev Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CenterCopyStr1)
										},
									},
									PushButton{
										Text:  "Pattern Hash Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CenterCopyStr2)
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
										AssignTo: &CenterSearchtext,
										MaxSize: Size{1000, 22},
										Text: "",
									},
									PushButton{
										Text:  "Search",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
										},
									},
								},
							},	
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										MinSize: Size{0, 22},
										Text: "Decoding Center Data : ",
									},
								},
							},	
							TextEdit{
								AssignTo: &CenterTe,
							},		
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Data Copy",
										MaxSize: Size{100, 22},
										OnClicked: func(){ 
											walk.Clipboard().SetText(CenterTe.Text())
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
	if err := CenterTe.SetReadOnly(true); err != nil {
		return nil, err
	}	
	
	return p, nil
}

