package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Directory struct {
	name     string
	parent   *Directory
	children []*Directory
}

func NewDirectory(name string, parent *Directory) *Directory {
	return &Directory{name: name, parent: parent}
}

var _ walk.TreeItem = new(Directory)

func (d *Directory) Text() string {
	return d.name
}

func (d *Directory) Parent() walk.TreeItem {
	if d.parent == nil {
		return nil
	}

	return d.parent
}

func (d *Directory) ChildCount() int {
	if d.children == nil {
		if err := d.ResetChildren(); err != nil {
			log.Print(err)
		}
	}
	return len(d.children)
}

func (d *Directory) ChildAt(index int) walk.TreeItem {
	return d.children[index]
}

func (d *Directory) Image() interface{} {
	return d.Path()
}

func (d *Directory) ResetChildren() error {
	d.children = nil
	dirPath := d.Path()
	if err := filepath.Walk(d.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if info == nil {
				return filepath.SkipDir
			}
		}
		name := info.Name()
		if !info.IsDir() || path == dirPath {
			return nil
		}
		d.children = append(d.children, NewDirectory(name, d))
		return filepath.SkipDir
	}); err != nil {
		return err
	}
	return nil
}

func (d *Directory) Path() string {
	elems := []string{d.name}
	dir, _ := d.Parent().(*Directory)
	for dir != nil {
		elems = append([]string{dir.name}, elems...)
		dir, _ = dir.Parent().(*Directory)
	}
	return filepath.Join(elems...)
}

type DirectoryTreeModel struct {
	walk.TreeModelBase
	roots []*Directory
}

var _ walk.TreeModel = new(DirectoryTreeModel)
func NewDirectoryTreeModel() (*DirectoryTreeModel, error) {
	model := new(DirectoryTreeModel)
	drives, err := walk.DriveNames()
	if err != nil {
		return nil, err
	}
	for _, drive := range drives {
		switch drive {
		case "A:\\", "B:\\":
			continue
		}
		model.roots = append(model.roots, NewDirectory(drive, nil))
	}
	return model, nil
}

func NewModelPath() *DirectoryTreeModel {
	model := new(DirectoryTreeModel)
	dirname:=filesetpath
	model.roots = append(model.roots, NewDirectory(dirname, nil))
	return model
}

func (*DirectoryTreeModel) SavePath() bool {
	//ex=strings.Split(chainpath,filepathSeparator)
	return true
}

func (*DirectoryTreeModel) LazyPopulation() bool {
	return true
}

func (m *DirectoryTreeModel) RootCount() int {
	return len(m.roots)
}

func (m *DirectoryTreeModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}

type FileInfo struct {
	Name     string
	Size     int64
	Modified time.Time
}

type FileInfoModel struct {
	walk.SortedReflectTableModelBase
	dirPath string
	items   []*FileInfo
}

var _ walk.ReflectTableModel = new(FileInfoModel)

func NewFileInfoModel() *FileInfoModel {
	return new(FileInfoModel)
}

func (m *FileInfoModel) Items() interface{} {
	return m.items
}

func (m *FileInfoModel) SetFilePath(dirPath string) error {
	m.dirPath = dirPath
	m.items = nil

	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if info == nil {
				return filepath.SkipDir
			}
		}
		name := info.Name()
		if path == dirPath || (!info.IsDir() && shouldExclude(name)) {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		item := &FileInfo{
			Name:     name,
			Size:     info.Size(),
			Modified: info.ModTime(),
		}
		m.items = append(m.items, item)
		return nil
	}); err != nil {
		return err
	}
	m.PublishRowsReset()
	return nil
}

func (m *FileInfoModel) SetDirPath(dirPath string) error {
	m.dirPath = dirPath
	m.items = nil

	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if info == nil {
				return filepath.SkipDir
			}
		}
		name := info.Name()
		if path == dirPath || (!info.IsDir() && shouldExclude(name)) {
			return nil
		}
		item := &FileInfo{
			Name:     name,
			Size:     info.Size(),
			Modified: info.ModTime(),
		}

		m.items = append(m.items, item)

		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return err
	}
	m.PublishRowsReset()
	return nil
}

func (m *FileInfoModel) Image(row int) interface{} {
	return filepath.Join(m.dirPath, m.items[row].Name)
}

func shouldExclude(name string) bool {
	ln:=len(name)
	if ln>30 {
		switch name[len(name)-4:] {
		case ".blc",".chn",".cub":
			return false
		}
	}
	return true
}


var FileLineEdit1 *walk.LineEdit
var FileTreeView *walk.TreeView
var actName string


func newFilePage(parent walk.Container) (Page, error) {
	var splitter *walk.Splitter
	var tableView *walk.TableView
	var msgForm walk.Form
	FileTreeViewFlag:=false	

	rootModel, err := NewDirectoryTreeModel()
	if err != nil {
		log.Fatal(err)
	}

	treeModel:= NewModelPath()
	tableModel := NewFileInfoModel()

	p := new(basePage)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "File",
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			Composite{
				StretchFactor: 1,
				MaxSize: Size{0, 22},
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "BasePath:",
					},
					LineEdit{
						AssignTo: &FileLineEdit1,
						Text: filesetpath,
					},
					PushButton{
						Text:  "Select",
						OnClicked: func() { FilePathSetting()  },
					},
					PushButton{
						Text:  "Save Path",
						OnClicked: func() { 
							filesetpath=FileLineEdit1.Text()
						},
					},
					PushButton{
						Text:  "Move Page",
						OnClicked: func() {  },
					},
				},
			},			
			HSplitter{
				AssignTo: &splitter,
				StretchFactor: 50,
				MinSize: Size{0, 400},
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text:  "Root Drive",
										OnClicked: func() {  
											FileTreeView.SetModel(rootModel)
											FileTreeViewFlag=false	
										},
									},
									PushButton{
										Text:  "Save Path",
										OnClicked: func() { 
											if FileTreeViewFlag {
												filesetpath=FileTreeView.CurrentItem().(*Directory).Path()
												FileLineEdit1.SetText(filesetpath)
											} else {
												walk.MsgBox(msgForm,"Error","At first select file path.",walk.MsgBoxOK|walk.MsgBoxIconError)							
											}
										},
									},
									PushButton{
										Text:  "Load Path",
										OnClicked: func() { 
											FileTreeView.SetModel(NewModelPath())
											FileTreeView.SetExpanded(FileTreeView.CurrentItem(),true)
											FileTreeViewFlag=false	
										},
									},
								},
							},	
							TreeView{
								AssignTo: &FileTreeView,
								Model:    treeModel,
								OnCurrentItemChanged: func() {
									FileTreeViewFlag=true
									dir := FileTreeView.CurrentItem().(*Directory)
									if err := tableModel.SetFilePath(dir.Path()); err != nil {
										walk.MsgBox(msgForm,"Error",err.Error(),walk.MsgBoxOK|walk.MsgBoxIconError)
									}
								},
							},

						},
					},	
					TableView{
						AssignTo:      &tableView,
						StretchFactor: 2,
						Columns: []TableViewColumn{
							TableViewColumn{
								DataMember: "Name",
								Width:      192,
							},
							TableViewColumn{
								DataMember: "Size",
								Format:     "%d",
								Alignment:  AlignFar,
								Width:      64,
							},
							TableViewColumn{
								DataMember: "Modified",
								Format:     "2006-01-02 15:04:05",
								Width:      120,
							},
						},
						Model: tableModel,
						OnCurrentIndexChanged: func() {
							if index := tableView.CurrentIndex(); index > -1 {
								name := tableModel.items[index].Name
								dir := FileTreeView.CurrentItem().(*Directory)
								fpath := filepath.Join(dir.Path(), name)
								actName=FileAction(fpath)
								for _, action := range mpmw.pageActions {
									action.SetChecked(false)
									logger.Print(actName)
									if action.Text()==actName {
										mpmw.setCurrentAction(action)
									}
								}
							}
						},
					},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}
	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	return p, nil
}


func FilePathSetting() error {
	var dlgForm walk.Form
	dlg := new(walk.FileDialog)

	filesetpath=ArrangeSeparator(filesetpath)
	dlg.FilePath = filesetpath
	dlg.Title = "Select an cubechain path"
    
	if ok, err := dlg.ShowBrowseFolder(dlgForm); err != nil {
        logger.Print("[Dialog Error:]")
        logger.Print(err)
		return err
	} else if !ok {
		return nil
	}
	filesetpath=dlg.FilePath
	FileLineEdit1.SetText(filesetpath)
	FileTreeView.SetModel(NewModelPath())
	return nil
}

func FileAction(path string) string {
	result:=""
	ext:=path[len(path)-4:]
	switch ext {
		case ".chn": 
			chainpath=path
			result="Chain"
		case ".cub":
			cubepath=path
			result="Cube"
		case ".blc":
			blockpath=path
			result="Block"
		default: return ""
	}
	return result
}
