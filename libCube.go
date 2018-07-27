package main

import (
	"fmt"
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

var filepathSeparator=string(filepath.Separator)
var datacenterHost string
var filesetpath string
var chainpath string
var cubepath string
var blockpath string
var downloadpath string
var Datanumber=10000

var cube CBlock
var cubechain Cubechain

type Cubechain struct {
	Verify		int
	Chain		[]CChain
}

type CChain struct {
	Index		int
	Timestamp	int
	Chash       string
}

type CubeBlock struct {
	Index		int
	Timestamp	int
	Cube		[27]Block
	Chash       string
}

type Cubing struct {
	Index		int
	Timestamp	int
	Hash1		[27]string
	Hash2		[27]string
	Chash       string
}

type Block struct {
	Index			int
	Cubeno			int
	Timestamp		int
	Data			[]byte
	Hash			string
	PrevHash		string
	PrevCubeHash	string
	Nonce			int
}

type CBlock struct {
	Index		int
	Timestamp	int
	Chash       string
}

func CubePath(idx int) string {
	divn:=idx/Datanumber
	divm:=idx%Datanumber
	if divm>0 {
		divn++
	} else if divm==0 {
		divm=Datanumber
	}
	if divn==0 {
		divn++
		divm=1
	}
	nhex:=fmt.Sprintf("%x",Datanumber)
	mcnt:=len(nhex)
	nstr:=fmt.Sprintf("%0.5x",divn)
	mstr:=fmt.Sprintf("%0."+strconv.Itoa(mcnt)+"x",divm)
	path:=CubeBasePath(cubepath)
	dirname:=path+filepathSeparator+nstr+filepathSeparator+mstr
	if dirExist(dirname)==false {
		if err:=os.MkdirAll(dirname, os.FileMode(0755)); err!=nil {
			return "Directory not found\\1\\1"
		}	
	}	
	dirname=strings.Replace(dirname,filepathSeparator+filepathSeparator,filepathSeparator, -1)
	return dirname
}

func CubePathNum(path string) int {
	result:=0
	add:=0
	if strings.Index(path,".")>=0 {	
		add++				
	}
	split:=strings.Split(path, filepathSeparator)
	slen:=len(split)
	nint,_:=strconv.ParseUint(split[slen-(2+add)],16,32)
	mint,_:=strconv.ParseUint(split[slen-(1+add)],16,32)
	result=(int(nint)-1)*Datanumber+int(mint)
	return result
}

func CubeBasePath(path string) string{
	if strings.Index(cubepath,".cub")>=0 || strings.Index(cubepath,".cbi")>=0 || strings.Index(cubepath,".blc")>=0 {
		ex:=strings.Split(cubepath,filepathSeparator)
		exp:=ex[0:len(ex)-3]
		path=""
		for _,v:=range exp {
			path+=v+filepathSeparator
		}
	}
	return path
}

func CubingFileRead(index int) Cubing {
	var cubing Cubing
	path:=CubePath(index)
	filename:=fileSearch(path,".cbi")
	err:=pathRead(path+filepathSeparator+filename,&cubing)
	if err!=nil {
		logger.Print(err)
	}
	return cubing
}

func ChainFileRead(dirPath string) Cubechain {
	var cubechain Cubechain
	if strings.Index(dirPath,".chn")>=0 {
	} else {
		filename:=fileSearch(dirPath,".chn")
		dirPath=dirPath+filepathSeparator+filename
	}
	err:=pathRead(dirPath,&cubechain)
	if err!=nil {
		logger.Print(err)
	}
	return cubechain
}

func pathWrite(path string, object interface{}) error {
	file,err:=os.Create(path)
	if err==nil {
		encoder:=gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func pathRead(path string, object interface{}) error {
	file,err:=os.Open(path)
	if err==nil {
		decoder:=gob.NewDecoder(file)
		err=decoder.Decode(object)
	}
	file.Close()
	return err
}

func fileSearch(dirname string,find string) string{
    result:=""
	d,err:=os.Open(dirname)
    if err != nil {
        logger.Print(err)
    }
    defer d.Close()
    fi, err:=d.Readdir(-1)
    if err != nil {
        logger.Print(err)
    }
    for _, fi:=range fi {
        if fi.Mode().IsRegular() {
            fstr:=fi.Name()
			if strings.Index(fstr,find)>=0 {
				result=fi.Name()
				return result
			}
        }
    }
	return result
}

func MaxFind(dirpath string) string {
	find:="0"
    d, err:=os.Open(dirpath)
    if err != nil {
        fmt.Println(err)
    }
    defer d.Close()
	fi, err:=d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
    }
    for _, fi:=range fi {
        if fi.Mode().IsRegular() {
        } else {
  			if fi.Name()>find {
				find=fi.Name()
			}
		}
   }
   return find
}

func BlockPathIndex(path string) (int,int) {
	idx:=0
	cubeno:=0
	if strings.Index(path,".blc")>=0 {
		ex:=strings.Split(path,"_")
		exp:=strings.Split(ex[0],filepathSeparator)
		idx,_=strconv.Atoi(exp[len(exp)-1])
		cubeno,_=strconv.Atoi(ex[1])
		cubeno++
	}
	return idx,cubeno
}

func BlockPath(idx int,cno int) string {
	find:=strconv.Itoa(idx)+"_"+strconv.Itoa(cno-1)+"_"	
    dirname:=CubePath(idx)
	filename:=fileSearch(dirname,find)
	result:=dirname+filepathSeparator+filename
	result=strings.Replace(result,filepathSeparator+filepathSeparator,filepathSeparator, -1)
	return result
}

func CubeFileCount(rootpath string,filter string) int {
	result:=0
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
		} else if filepath.Ext(path)==filter {
			result++
		}
		return nil
	})	
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}	
	return result
}

func CurrentHeight(path string) int {
	result:=0
	f:=MaxFind(path+filepathSeparator)
	if f=="0" {
		return 1
	}
	f2:=MaxFind(path+filepathSeparator+f)
	if f2=="0" {
		return 1
	}
	nint,_:=strconv.ParseUint(f,16,32)
	mint,_:=strconv.ParseUint(f2,16,32)
	result=(int(nint)-1)*Datanumber+int(mint)
	if fileSearch(CubePath(result),".cub")>"" {
		result++
	}
	return result	
}

func dirExist(dirName string) bool{
	result:=true
	_,err:=os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err ) {
			result=false
		}
	}
	return result
}

func StrToByte(str string) []byte {
	sb := make([]byte, len(str))
	for k, v := range str {
		sb[k] = byte(v)
	}
	return sb[:]
}

func ByteToStr(bytes []byte) string {
	var str []byte
	for _, v := range bytes {
		if v != 0x0 {
			str = append(str, v)
		}
	}
	return fmt.Sprintf("%s", str)
}

func ArrangeSeparator(path string) string {
	result:=strings.Replace(path,filepathSeparator+filepathSeparator,filepathSeparator, -1)
	return result
}

func PathToDir(path string,separator string) string{
	if separator=="" {
		separator=filepathSeparator
	}
	ex:=strings.Split(path,separator)
	ex1:=ex[:len(ex)-1]
	dir:=strings.Join(ex1,separator)
	return dir
}

func PathToFile(path string,separator string) string{
	if separator=="" {
		separator=filepathSeparator
	}
	ex:=strings.Split(path,separator)
	file:=ex[len(ex)-1]
	return file
}

func FileToExt(file string) string{
	ex:=strings.Split(file,".")
	fileext:=ex[len(ex)-1]
	return fileext
}

func setfilePath() string {
	path,_:=os.Executable()
	dir:=PathToDir(path,"")
	return dir+filepathSeparator+"SettingPath.dat"
}
