package command

import (
	"client-put/tool"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
function:通过post方式，将打包人、包名和包的描述都传到服务端，服务端将消息进行提取，
如果包名和包的描述为空或者前后不对应（包名中可以将版本号和release号都带上，这样的话就不需要再添加-v和-r的参数了）
则返回错误。每个包都起一个打包的协程。如果有附带的脚本，可以添加参数--preinstall、--postinstall、--preuninstall、
--postuninstall。
wshell pack wuwei -n name1 name2 name3... -d des1 des2 des3...
fpm命令分别获取包名打包
author:wuwei
input:params ...string。可以有一个参数，也可以没有参数，没有参数时默认将“/usr/tmp/wshell/"目录下的文件夹压缩
output:void
Data:Wed Feb 13 2019
*/
func Create(cmd string, params ...string) {

	//先检查是否有打包人pack这一选项，如果没有提示错误，并填写打包人
	if len(params) == 0 {
		CmdHelp(cmd)
	}
	//检查-n后面的参数个数和-d后面的参数个数是否相同。方法是找到-n和-d的位置，利用切片截取长度。
	subN := FindStr(params, "-n")
	subD := FindStr(params, "-d")
	if subD == -1 || subN == -1 {
		fmt.Println("请填写-n与-d参数")
		return
	}
	if len(params[subN:subD]) != len(params[subD:]) {
		fmt.Println("包名称和包的描述不匹配")
	}
	//检查-n后面的参数是否满足包名-版本号-release号的格式
	for _, v := range params[subN+1 : subD] {
		if !CheckFileName(v) {
			fmt.Println("包的命名有误，请检查包名")
		}
	}
	//将信息分多次提交到后端，类似于以前的分多次提交网页版信息。提交后可以批量查看打包的进度，然后可以批量的上传包。
	for i, v := range params[subN+1 : subD] {
		name, version, release, _ := GetRPMInfo(v)
		//通过post方式提交rpm包的信息
		data := make(url.Values)
		data["pkg"] = []string{params[0]}
		data["n"] = []string{name}
		data["vNum"] = []string{version}
		data["iNum"] = []string{release}
		data["dp"] = []string{params[subD+i+1]}
		fmt.Println(name + "打包中...")
		res, err := http.PostForm("http://10.209.16.164:9000/wshellfpmInfo", data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(name + "打包成功")
		defer res.Body.Close()
	}

}

/*
function:查找字符串数组中某个字符串的位置
author：wuwei
input1:字符串数组strList []string
input2:查找的字符串 substr string
output:substr的位置i
*/
func FindStr(strList []string, substr string) int {
	for i := 0; i < len(strList); i++ {
		if strList[i] == substr {
			return i
		}
	}
	return -1
}

/*
function:检查rpm包名是否是以包名-版本号-release号的形式命名的
author：wuwei
input：FileName string
output：true：命名方式正确 false：命名格式错误
*/
func CheckFileName(FileName string) bool {
	fileInfo := strings.Split(FileName, "-")
	length := len(fileInfo)
	if length >= 3 {
		var srcPath string
		for i := 0; i < length-2; i++ {
			if i == length-3 {
				srcPath += fileInfo[i]
			} else {
				srcPath += fileInfo[i] + "-"
			}
		}

		if matched := strings.Index(srcPath, "/"); matched != -1 {
			fmt.Println("Build failed: the name not support absolute path or include '/'.")
			return false
			os.Exit(1)
		}

		version := fileInfo[length-2]
		if matched := strings.Index(version, "-"); matched != -1 {
			fmt.Println("Build failed: the version cannot be include '-'.")
			return false
			os.Exit(1)
		}

		iteration := fileInfo[length-1]
		if tool.IsNumStr(iteration) == false {
			fmt.Println("Build failed: the release must bu integers.")
			return false
			os.Exit(1)
		}
		if matched := strings.Index(iteration, "-"); matched != -1 {
			fmt.Println("Build failed: the release cannot be include '-'.")
			return false
			os.Exit(1)
		}
		return true
	}
	return false
}

/*
function:获取包名中的包名、版本号以及release号
author：wuwei
input：FileName string
output1：name 包名
output2：version 版本号
output3：release 发行号
*/
func GetRPMInfo(FileName string) (name string, version string, release string, ok bool) {
	fileInfo := strings.Split(FileName, "-")
	length := len(fileInfo)
	if length >= 3 {
		for i := 0; i < length-2; i++ {
			if i == length-3 {
				name += fileInfo[i]
			} else {
				name += fileInfo[i] + "-"
			}
		}

		if matched := strings.Index(FileName, "/"); matched != -1 {
			fmt.Println("Build failed: the name not support absolute path or include '/'.")
			ok = false
			return
		}

		version = fileInfo[length-2]
		if matched := strings.Index(version, "-"); matched != -1 {
			fmt.Println("Build failed: the version cannot be include '-'.")
			ok = false
			return
		}

		release = fileInfo[length-1]
		if tool.IsNumStr(release) == false {
			fmt.Println("Build failed: the release must bu integers.")
			ok = false
			return
		}
		if matched := strings.Index(release, "-"); matched != -1 {
			fmt.Println("Build failed: the release cannot be include '-'.")
			ok = false
			return
		}
		ok = true
		return
	}
	return
}
