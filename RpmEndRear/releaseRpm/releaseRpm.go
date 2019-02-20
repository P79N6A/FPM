package releaseRpm

import (
	"FPMtestUpload/dataStruct"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

/*
function:HTTP callback function
author:wuwei
input:response http.ResponseWriter
input:request *http.Request
output:void
Data Thu Aur 31 2018
*/
func ReleaseRpmC(w http.ResponseWriter, r *http.Request) {

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
			ReleaseRpm(rpmName[0])
			result.Data.State = rpmName[0] + "upload sucess!!!"
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

/*
function:Release Rpm
author:wuwei
input:rpmName string
output:bool
Data:Thu Aur 15 2018
*/
func ReleaseRpm(rpmName string) bool {

	var releaseCmd = "cp /home/i-wuwei1/FPM/RPMS/" + rpmName + " /data/yum/repository/addops/6/os/x86_64/"
	cmd := exec.Command("/bin/sh", "-c", releaseCmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return false
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return false
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return false
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return false
	}
	fmt.Println(bytes)
	var repoctlu = "/data/yum/tools/repoctl.sh -u" //同步yum源脚本
	var repoctlm = "/data/yum/tools/repoctl.sh -s -d mirror"
	cmd = exec.Command("/bin/sh", "-c", repoctlu)
	stdout, err = cmd.StdoutPipe() //将命令的执行结果放到了stdout中，这个函数只有加了cmd.wait之后才有效，他是在start和wait之间执行的函数
	if err != nil {
		fmt.Printf("Error:release stdout:%s\n", err)
		return false
	}
	if err := cmd.Start(); err != nil { //开始执行命令
		fmt.Println("Error:The release command is err,", err)
		return false
	}

	info, err := ioutil.ReadAll(stdout) //将命令执行中的输出内容读出，放到了info中
	if err != nil {
		fmt.Println("release ReadAll Stdout:", err)
		fmt.Println("info:", info)
		return false
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("release wait:", err.Error())
		return false
	}

	cmd = exec.Command("/bin/sh", "-c", repoctlm)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:release stdout:%s\n", err)
		return false
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The release command is err,", err)
		return false
	}

	info, err = ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("release ReadAll Stdout:", err)
		fmt.Println("info:", info)
		return false
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("release wait:", err.Error())
		return false
	}
	return true
}
