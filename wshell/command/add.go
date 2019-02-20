package command

import (
	"client-put/tool"
	"fmt"
	"os"
	//"os"
	"path/filepath"
	"strings"
)

func Add(cmd string, params ...string) {

	currentPwd, _ := os.Getwd() //获取当前目录的目录路径
	//fmt.Println(len(tool.GetFileName(currentPwd)))
	if len(params) != 1 && len(params) != len(tool.GetFileName(currentPwd)) {
		CmdHelp(cmd)
		return

	}

	if !strings.Contains(params[0], "BIPACKAGE") {
		if !tool.Exists(params[0]) {
			fmt.Println(params[0] + " not exist!")
			return
		}
	}

	//读取usr/tmp/wshell目录下的内容，即目的文件夹的名称，vartualDir[0]即为所要拿到的文件夹名称，因为该目录下只有一个文件夹
	vartualDir := tool.GetDirName(dirPrefix)
	if vartualDir == nil {
		fmt.Println("Please use 'wshell init' create dir")
		return
	}

	//获取源主机上的所要打包的文件的路径
	dirPath, _ := filepath.Abs(params[0])
	//由于dirPath是绝对路径，应该去除掉最后一项，因此要将最后一个路径分割，将“/”
	position := strings.LastIndex(dirPath, "/")
	dirPath = dirPath[0:position]
	//fmt.Println("the path is" + dirPath)

	//在虚拟目录下创建获取到的文件的路径
	mkdirCmd := "mkdir -p " + dirPrefix + "/" + vartualDir[0] + dirPath
	//fmt.Println("mkdirCmd: " + mkdirCmd)
	if !tool.ExecCmdRoutine(mkdirCmd) {
		fmt.Println("create virtual path error")
		return
	}

	//将要打包的文件放入虚拟目录的相应的文件夹下
	addCmd := "/bin/cp -Rf " + params[0] + " " + dirPrefix + "/" + vartualDir[0] + dirPath
	if len(params) == len(tool.GetFileName(currentPwd)) {
		var paramsList string = ""
		for i := 0; i < len(tool.GetFileName(currentPwd)); i++ {
			paramsList += params[i] + " "
		}

		addCmd = "cp -R " + paramsList + " " + dirPrefix + "/" + vartualDir[0] + dirPath
	}
	//fmt.Println("addcmd: " + addCmd)
	if tool.ExecCmdRoutine(addCmd) {
		fmt.Println("Add package success!")
		return
	}
	fmt.Println("Add package error!")

}
