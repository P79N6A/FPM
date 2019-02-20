package command

import (
	"FPMtestUpload/dataStruct"
	"client-put/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
function:上传rpm包，支持批量上传rpm包
wshell upload name1-version1 name2-version2 name3-version3
author:wuwei
input1:cmd string 程序后面的第一个参数params[0]
input2:params ...string。可以有一个参数，也可以没有参数，没有参数时默认上传“/usr/tmp/wshell/"目录下的包。
output:void
Data:Thu Feb 14 2019
*/
func Upload(cmd string, params ...string) {
	//执行发布包的功能，向后台传输要发布的包的名字和版本号，后台将该包复制到指定目录下，然后
	//执行发布的脚本。参数个数为0时，获取默认的虚拟路径下的包名，并发送到后端。参数个数大于
	//等于1时，循环发送所有包的名字到后端，执行多次发布
	var localFile string = ""
	if len(params) >= 1 {
		for i, _ := range params {
			localFile += params[i] + ":"
		}
	} else {
		os.Chdir("/usr/tmp/wshell")
		vartualDir := tool.GetDirName(dirPrefix)
		if vartualDir == nil {
			fmt.Println("Please use 'wshell init' create dir")
			return
		}
		destName := vartualDir[0]
		localFile = destName
	}
	//获取每个包的包名和版本号n和vNum
	for i, v := range strings.Split(localFile, ":") {
		if i >= len(params) {
			return
		}
		name, version, _, ok := GetRPMInfo(v)
		if !ok {
			return
		}
		//将信息分多次提交到后端，类似于以前的分多次提交网页版信息。提交后可以批量查看打包的进度，然后可以批量的上传包。
		data := make(url.Values)
		data["n"] = []string{name}
		data["vNum"] = []string{version}
		res, err := http.PostForm("http://10.209.16.164:9000/uploadrpm", data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer res.Body.Close()
		var result dataStruct.ResponseInfo
		body, err := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &result); err == nil {
			fmt.Println(result.Data.State)
			continue
		}
		fmt.Println(name + version + "包发布失败")

	}
}
