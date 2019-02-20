package dataStruct

import (
	"sync"
)

/*
function:定义rpm包的基本信息
author：wuwei
Data：Thu Aur 16 2018
*/
type FpmBasicInfo struct {
	Name            string //rpm包名
	Version         string //rpm包版本号
	Iteration       string //rpm包release号
	Dependence      string //依赖的库，有多个库时用空格分开
	Packager        string //打包的人员
	Description     string //包的描述
	PreInstallSh    string //安装rpm包前执行的脚本
	PostInstallSh   string //安装rpm包后执行的脚本
	PreUninstallSh  string //卸载rpm包前执行的脚本
	PostUninstallSh string //卸载rpm包后执行的脚本
}

/////////////////////////////InfoOutputStart///////////////////////////////
/*
function:制作一个线程安全的map，该结构体用于并发执行rpm打包时，将不同的rpm包的制作信息和rpm包
绑定到一起（key:包名，value:包信息），不会出现同时打包多个包时的程序运行出错问题。
*/
type InfoOutput struct {
	data map[string][]byte
	Lock sync.RWMutex
}

var RpmInfoOutput InfoOutput //并发运行时，用于查看每个包的打包状态
//var UploadFile InfoOutput    //并发运行时，查看是否已经通过网页上传文件（true），没有通过网页上传的话为（false)

/*
function:初始化map
*/
func (d *InfoOutput) InfoInit() {
	d.data = make(map[string][]byte)
}

/*
function:读取map对应的key的value值
*/
func (d *InfoOutput) InfoGet(k string) []byte {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.data[k]
}

/*
function:设置map中key对应的map值
*/
func (d *InfoOutput) InfoSet(k string, v []byte) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.data[k] = v
}

/////////////////////////////InfoOutputEnd///////////////////////////////

/////////////////////////////ErrorOutputStart///////////////////////////////
/*
function:存放错误包信息的线程安全的map，存放制作包的流程中可能发生的错误
*/
type ErrorOutput struct {
	data map[string]error
	Lock sync.RWMutex
}

var Err ErrorOutput

func (d *ErrorOutput) ErrInit() {
	d.data = make(map[string]error)
}

func (d *ErrorOutput) ErrGet(k string) error {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.data[k]
}

func (d *ErrorOutput) ErrSet(k string, v error) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.data[k] = v
}

/////////////////////////////ErrorOutputEnd///////////////////////////////

/////////////////////////////ApiDataStart///////////////////////////////
/*
function:返回json字符串，与前端交互的接口
*/
type ResponseData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	State   string `json:"state"`
}

type ResponseInfo struct {
	ErrorMsg string       `json:"errmsg"`
	ErrorNo  int          `json:"errno"`
	Data     ResponseData `json:"data"`
}

/////////////////////////////ApiDataEnd///////////////////////////////

/////////////////////////////ApiDataStart///////////////////////////////
/*
function:json字符串，返回数据库中查询出的数据
*/
type ResponseMysqlData struct {
	Iteration string `json:"iteration"`
	Version   string `json:"version"`
}

type ResponseMysqlInfo struct {
	ErrorMsg string            `json:"errmsg"`
	ErrorNo  int               `json:"errno"`
	Data     ResponseMysqlData `json:"data"`
}

/////////////////////////////ApiDataEnd///////////////////////////////
