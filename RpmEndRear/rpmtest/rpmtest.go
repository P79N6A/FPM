package rpmTest

import (
	"FPMtestUpload/dataStruct"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

/*
function:test Rpm
author:wuwei
input:w http.ResponseWriter
input:r *http.Request
output:void
Data:Thu Aur 31 2018
*/
func RpmTest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	name := r.PostFormValue("n")
	version := r.PostFormValue("vNum")
	var result dataStruct.ResponseInfo
	result.Data.Name = name
	result.Data.Version = version
	if dataStruct.Err.ErrGet(name+version) != nil {
		result.ErrorMsg = dataStruct.Err.ErrGet(name + version).Error()
		result.ErrorNo = 105
		jsonResult, er := json.Marshal(result)
		if er == nil {
			w.Write(jsonResult)
		}
		return
	}
	if len(dataStruct.RpmInfoOutput.InfoGet(name+version)) > 0 {
		if bytes.Compare(dataStruct.RpmInfoOutput.InfoGet(name+version), []byte("true")) == 0 {
			result.Data.State = "rpm building..."
		} else {
			rpmName := strings.Split(string(dataStruct.RpmInfoOutput.InfoGet(name+version)), " ")
			result.Data.State = "rpm -i http://k2198v.add.bjyt.qihoo.net:80/" + rpmName[0]
		}
	} else {
		result.ErrorMsg = "no package infomation,please check up the name of package~~~"
		result.ErrorNo = 104
	}
	jsonResult, er := json.Marshal(result)
	if er == nil {
		w.Write(jsonResult)
	}

}
