package progressQuery

import (
	"FPMtestUpload/dataStruct"
	"bytes"
	"encoding/json"
	"net/http"
)

/*
function:query progress
author:wuwei
input:http.ResponseWriter
input:r *http.Request
output:void
Data:Thu Aur 22 2018
*/
func ProgressQuery(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	name := r.PostFormValue("n")
	version := r.PostFormValue("vNum")
	var result dataStruct.ResponseInfo
	result.Data.Name = name
	result.Data.Version = version
	if dataStruct.Err.ErrGet(name+version) != nil {
		//result.Data.State = "package infomation error~~~"
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
			result.Data.State = string(dataStruct.RpmInfoOutput.InfoGet(name + version)[:])
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
