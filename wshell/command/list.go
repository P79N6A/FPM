package command

import (
	"client-put/tool"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
)

/*
func List(cmd string, params ...string) {
	if len(params) != 1 && len(params) != 2 {
		CmdHelp(cmd)
		return
	}
	listCmd := ""
	if !tool.Exists(params[0]) {
		fmt.Println(params[0] + " is not exist!")
		return
	}

	if len(params) == 1 {
		listCmd = "tree " + params[0]
	} else if len(params) == 2 {
		listCmd = "tree -L " + params[1] + " " + params[0]
	}

	if tool.ExecCmdRoutine(listCmd) {
		return
	}



	fmt.Println("please check out filename")

}
*/

/*
function:以非命令行的方式列出目录下的所有文件
author:wuwei
input:cmd string 位list参数
input:params ...string 参数信息，可以有两个params[0]为文件夹信息，params[1]位level信息
output:void
Data:Thu Dec 27 2018
*/
func List(cmd string, params ...string) {
	if len(params) != 0 && len(params) != 1 && len(params) != 2 {
		CmdHelp(cmd)
		return
	}

	//读取usr/tmp/wshell目录下的内容，即目的文件夹的名称，vartualDir[0]即为所要拿到的文件夹名称，因为该目录下只有一个文件夹
	vartualDir := tool.GetDirName(dirPrefix) //vartualDir[0]本身不带"/"
	if vartualDir == nil {
		fmt.Println("Please use 'wshell init' create dir")
		return
	}

	if !tool.Exists(dirPrefix + "/" + vartualDir[0]) {
		fmt.Println(vartualDir[0] + " is not exist!")
		return
	}

	fi, _ := os.Lstat(dirPrefix + "/" + vartualDir[0])
	if len(params) == 0 {
		listDir(dirPrefix+"/"+vartualDir[0], fi, "", 0, 100)
		return
	} else if len(params) == 2 {
		maxDepth, _ := strconv.Atoi(params[1])
		listDir(dirPrefix+"/"+vartualDir[0]+"/"+params[0], fi, "", 0, maxDepth) //用户在输入路径时加入根目录”/“和不加入根目录都是可行的，这就需要程序中加入”/“。
		return
	} else {
		//设置正则表达式，判断在输入一个参数时是数字字符串还是纯字母字符串，纯字母字符串说明输入的参数代表目录的含义，因此需要列出该目录下的所有文件
		//pattern := "[-?[0-9]\\d*]"
		rNum := regexp.MustCompile("\\d")
		num := len(rNum.FindAllStringSubmatch(params[0], -1))

		if num != len(params[0]) {
			listDir(dirPrefix+"/"+vartualDir[0]+"/"+params[0], fi, "", 0, 100) //假设目录树的最大深度为100
		} else {
			maxDepth, _ := strconv.Atoi(params[0])
			listDir(dirPrefix+"/"+vartualDir[0], fi, "", 0, maxDepth)
		}
		return

	}

	fmt.Println("please check out filename")

}

/*
function:以树形的方式列出目录下的所有文件
author:wuwei
input:dirName string 文件夹信息
input:fileInfo os.FileInfo 文件的相关属性结构体，由os库中的Lstat或者stat获得
input:prefix string 目录树的前缀，一般为“”
input:depth int 目录的深度
output:void
Data:Thu Dec 28 2018
*/
func listDir(dirname string, fileInfo os.FileInfo, prefix string, depth int, maxDepth int) {
	if !fileInfo.IsDir() {
		return
	}

	if depth > maxDepth {
		return
	}

	fd, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		return
	}
	names, err := fd.Readdirnames(-1)
	fd.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	var fileInfos = make([]os.FileInfo, len(names))
	var num = 0
	sort.Strings(names)
	for i, name := range names {
		fileName := path.Join(dirname, name)
		fileInfo, err = os.Lstat(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileInfos[i] = fileInfo
		if fileInfo.IsDir() {
			num++
		}
	}
	var j, k = 0, len(fileInfos) - num
	var tmpFileInfos = make([]os.FileInfo, len(fileInfos))
	for _, fileInfo = range fileInfos {
		if fileInfo.IsDir() {
			tmpFileInfos[k], k = fileInfo, k+1
		} else {
			tmpFileInfos[j], j = fileInfo, j+1
		}
	}
	fileInfos = tmpFileInfos
	for i, fileInfo := range fileInfos {
		fileName := path.Join(dirname, fileInfo.Name())
		var ss string
		if i == len(fileInfos)-1 {
			fmt.Printf("%v└──%v\n", prefix, fileInfo.Name())
			ss = prefix + "     "
		} else {
			fmt.Printf("%v├──%v\n", prefix, fileInfo.Name())
			ss = prefix + "│   "
		}
		if fileInfo.IsDir() {
			listDir(fileName, fileInfo, ss, depth+1, maxDepth)
		}
	}
}
