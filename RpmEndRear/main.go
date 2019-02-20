// FPM project main.go
package main

import (
	"FPMtestUpload/S3Download"
	"FPMtestUpload/command"
	"FPMtestUpload/dataStruct"
	"FPMtestUpload/dbInterface"
	"FPMtestUpload/progressQuery"
	"FPMtestUpload/releaseRpm"
	"FPMtestUpload/rpmtest"
	"FPMtestUpload/tools"
	"FPMtestUpload/uploadFile"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

/*
function:Detecte Post value
author:wuwei
input:postValue string
output:bool
Data Fri Aur 17 2018
*/
func detectePostValue(postValue string, response http.ResponseWriter) bool {
	if postValue == "" {
		//response.Write([]byte("packager|rpmName|version|iteration is nil"))
		return false
	}
	return true
}

/*
function:HTTP callback function
author:wuwei
input:response http.ResponseWriter
input:request *http.Request
output:void
Data Thu Aur 16 2018
*/
func processFpmInfo(response http.ResponseWriter, request *http.Request) {
	//http.Redirect(response, request, "http://www.baidu.com", http.StatusNotFound)

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var result dataStruct.ResponseInfo
	rpmBasicInfo := dataStruct.FpmBasicInfo{"", "",
		"", "", "", "", "", "", "", ""}

	rpmBasicInfo.Packager = request.PostFormValue("pkg")
	rpmBasicInfo.Name = request.PostFormValue("n")
	rpmBasicInfo.Version = request.PostFormValue("vNum")
	rpmBasicInfo.Iteration = request.PostFormValue("iNum")
	rpmBasicInfo.Dependence = request.PostFormValue("dpd")
	rpmBasicInfo.Description = request.PostFormValue("dp")
	rpmBasicInfo.PreInstallSh = request.PostFormValue("pris")
	rpmBasicInfo.PostInstallSh = request.PostFormValue("pois")
	rpmBasicInfo.PreUninstallSh = request.PostFormValue("prus")
	rpmBasicInfo.PostUninstallSh = request.PostFormValue("pous")

	result.Data.Name = rpmBasicInfo.Name
	result.Data.Version = rpmBasicInfo.Version
	result.Data.State = "making a bag,please wait patiently..."
	jsonResult, er := json.Marshal(result)
	if er == nil {
		response.Write(jsonResult)
		fmt.Println(jsonResult)
		//renderHTML(response, "FPMClientCheck.htm", jsonResult)
	}
	fmt.Println(rpmBasicInfo)
	//Legality check
	if !detectePostValue(rpmBasicInfo.Packager, response) || !detectePostValue(rpmBasicInfo.Name, response) || !detectePostValue(rpmBasicInfo.Version, response) || !detectePostValue(rpmBasicInfo.Iteration, response) {
		return
	}

	dataStruct.RpmInfoOutput.InfoSet(rpmBasicInfo.Name+rpmBasicInfo.Version, nil)
	dataStruct.Err.ErrSet(rpmBasicInfo.Name+rpmBasicInfo.Version, nil)
	go command.ProcessCommand(rpmBasicInfo)

}

/*
function:HTTP callback function
author:wuwei
input:response http.ResponseWriter
input:request *http.Request
output:void
Data Thu Aur 16 2018
*/
func wshellProcessFpmInfo(response http.ResponseWriter, request *http.Request) {
	//http.Redirect(response, request, "http://www.baidu.com", http.StatusNotFound)

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var result dataStruct.ResponseInfo
	rpmBasicInfo := dataStruct.FpmBasicInfo{"", "",
		"", "", "", "", "", "", "", ""}

	rpmBasicInfo.Packager = request.PostFormValue("pkg")
	rpmBasicInfo.Name = request.PostFormValue("n")
	rpmBasicInfo.Version = request.PostFormValue("vNum")
	rpmBasicInfo.Iteration = request.PostFormValue("iNum")
	rpmBasicInfo.Dependence = request.PostFormValue("dpd")
	rpmBasicInfo.Description = request.PostFormValue("dp")
	rpmBasicInfo.PreInstallSh = request.PostFormValue("pris")
	rpmBasicInfo.PostInstallSh = request.PostFormValue("pois")
	rpmBasicInfo.PreUninstallSh = request.PostFormValue("prus")
	rpmBasicInfo.PostUninstallSh = request.PostFormValue("pous")

	result.Data.Name = rpmBasicInfo.Name
	result.Data.Version = rpmBasicInfo.Version
	result.Data.State = "making a bag,please wait patiently..."
	jsonResult, er := json.Marshal(result)
	if er == nil {
		response.Write(jsonResult)
		fmt.Println(jsonResult)
		//renderHTML(response, "FPMClientCheck.htm", jsonResult)
	}
	fmt.Println(rpmBasicInfo)
	//Legality check
	if !detectePostValue(rpmBasicInfo.Packager, response) || !detectePostValue(rpmBasicInfo.Name, response) || !detectePostValue(rpmBasicInfo.Version, response) || !detectePostValue(rpmBasicInfo.Iteration, response) {
		return
	}

	dataStruct.RpmInfoOutput.InfoSet(rpmBasicInfo.Name+rpmBasicInfo.Version, nil)
	dataStruct.Err.ErrSet(rpmBasicInfo.Name+rpmBasicInfo.Version, nil)
	command.ProcessCommand(rpmBasicInfo)

}

// 渲染页面并输出
func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	// 获取页面内容
	t, err := template.New(file).ParseFiles(file)
	tools.CheckErr(err)
	// 将页面渲染后反馈给客户端
	t.Execute(w, data)
}

func main() {

	dataStruct.RpmInfoOutput.InfoInit()
	dataStruct.Err.ErrInit()
	//dataStruct.UploadFile.InfoInit()
	err := dbInterface.ConnectMysql()
	if err != nil {
		fmt.Println(err)
		return
	}
	S3Download.Client = S3Download.NewS3Client()

	http.HandleFunc("/fpmInfo", processFpmInfo)                //执行rpm打包的一系列操作
	go http.HandleFunc("/wshellfpmInfo", wshellProcessFpmInfo) //执行rpm打包的一系列操作
	//http.HandleFunc("/upload", uploadFile.Upload)  //通过http上传文件夹
	http.HandleFunc("/UploadFileInfo", uploadFile.UploadFileInfo)  //通过http上传文件夹信息
	http.HandleFunc("/progressQuery", progressQuery.ProgressQuery) //打包进度查询
	http.HandleFunc("/rpmtest", rpmTest.RpmTest)                   //获取打包的rpm包的地址
	http.HandleFunc("/uploadrpm", releaseRpm.ReleaseRpmC)          //发布rpm包
	http.HandleFunc("/selectMysql", dbInterface.SelectMysql)       //数据库中源码信息包的查询

	http.ListenAndServe(":9000", nil)

}
