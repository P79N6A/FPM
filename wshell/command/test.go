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
function:批量发送包名称和包的版本号，返回测试包的地址，即支持批量获取测试包的地址
wshell test name1-version1 name2-version2 name3-version3...
author:wuwei
input1:cmd string 程序后面的第一个参数params[0]
input2:params ...string。可以有一个参数，也可以没有参数，没有参数时默认获取“/usr/tmp/wshell/"目录下的测试包地址。
output:void
Data:Thu Feb 14 2019
*/
func TEST(cmd string, params ...string) {
	//如果参数个数为0时，默认从虚拟目录中的包名的测试地址，如果参数不为0，则循环获取所有
	//参数的包名的测试地址，那处理思路就是定义一个字符串变量并赋予所有的包名，然后将该字
	//符串以“：”分割循环处理包名
	//当参数个数为0时，从虚拟路径下获取文件夹名称和版本号
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
		res, err := http.PostForm("http://10.209.16.164:9000/rpmtest", data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer res.Body.Close()
		var result dataStruct.ResponseInfo
		body, err := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &result); err == nil {
			fmt.Println(name + "-" + version + "包的测试地址为：" + result.Data.State)
			continue
		}
		fmt.Println("获取" + name + version + "测试地址失败")

	}
}
