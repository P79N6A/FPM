package uploadFile

import (
	"FPMtestUpload/dataStruct"
	"FPMtestUpload/dbInterface"
	"FPMtestUpload/tools"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
function:upload-callback upload file
author:wuwei
input:http.ResponseWriter
input:r *http.Request
output:void
Data:Thu Aur 22 2018
*/

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {

			fmt.Println(err)
			return
		}
		defer file.Close()

		var result dataStruct.ResponseInfo
		version, name, iteration, err := fileParse(handler.Filename)
		if err != nil {
			result.ErrorNo = 101
			result.ErrorMsg = err.Error()
			jsonResult, er := json.Marshal(result)
			if er == nil {
				w.Write(jsonResult)
			}
			fmt.Println(err)
			return
		}

		result.Data.Name = name
		result.Data.Version = version
		f, err := os.OpenFile("/home/i-wuwei1/FPM/SOURCES/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			result.ErrorNo = 102
			result.Data.State = "upload faild~~~"
			result.ErrorMsg = "upload faild~~~"
			jsonResult, er := json.Marshal(result)
			if er == nil {
				w.Write(jsonResult)
			}
			return
		}
		defer f.Close()

		io.Copy(f, file)
		errTar := tools.TarExtGz(handler.Filename)
		if errTar != nil {
			fmt.Println(errTar)
			result.ErrorNo = 103
			result.ErrorMsg = "extract file error"
			jsonResult, er := json.Marshal(result)
			if er == nil {
				w.Write(jsonResult)
			}
			return
		}
		//	dataStruct.UploadFile.InfoSet(name+version+iteration, []byte("true")) //表示源码文件已经通过http的方式上传

		//源码包的信息插入数据库
		var temp dataStruct.FpmBasicInfo
		temp.Name = name
		temp.Version = version
		temp.Iteration = iteration
		err = dbInterface.InsertMysql(temp)
		if err != nil {
			fmt.Println(err)
			return
		}

		result.Data.State = "upload sucess!!!"
		jsonResult, er := json.Marshal(result)
		if er == nil {
			w.Write(jsonResult)
		}
	}
}

/*
function:upload-callback upload file
author:wuwei
input:http.ResponseWriter
input:r *http.Request
output:void
Data:Thu Sep 22 2018
*/

func UploadFileInfo(w http.ResponseWriter, r *http.Request) {
	var temp dataStruct.FpmBasicInfo
	if r.Method == "POST" {
		temp.Name = r.PostFormValue("Name")
		temp.Version = r.PostFormValue("Version")
		temp.Iteration = r.PostFormValue("Iteration")
	}

	//源码包的信息插入数据库
	err := dbInterface.InsertMysql(temp)
	if err != nil {
		fmt.Println(err)

		return
	}
	fmt.Println(err)

}

/*
function:File name analysis
author:wuwei
input:fName string
output:version, baseName string, err error
Data:Thu Aur 29 2018
*/
func fileParse(fName string) (version, baseName, iteration string, err error) {
	bn := filepath.Base(fName)
	tempStr := strings.Split(bn, "-")
	if len(tempStr) != 3 {
		err = fmt.Errorf("%s naming not standard. Please use `fshell build` to create package", bn)
		return
	}

	re := regexp.MustCompile(`\.tar\.gz`)
	if matched := re.MatchString(tempStr[2]); !matched {
		err = fmt.Errorf("%s naming not standard. ", fName)
		return
	}
	iter := strings.Split(tempStr[2], ".")
	version = tempStr[1]
	baseName = tempStr[0]
	iteration = iter[0]
	var an int
	an, err = strconv.Atoi(iter[0])
	if err != nil {
		err = fmt.Errorf("%s naming not standard. ", fName)
		fmt.Println(an)
	}
	return
}
