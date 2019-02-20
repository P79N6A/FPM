package command

import (
	"client-put/tool"
	"fmt"
	"path/filepath"
	//	"io"
)

func Delete(cmd string, params ...string) {
	if len(params) != 1 && len(params) != 0 {
		CmdHelp(cmd)
		return
	}

	//如果要删除的文件不存在，且传入的参数长度不为0
	if len(params) != 0 && !tool.Exists(params[0]) {
		fmt.Println(params[0] + " not exist!")
		return
	}

	//读取usr/tmp/wshell目录下的内容，即目的文件夹的名称，vartualDir[0]即为所要拿到的文件夹名称，因为该目录下只有一个文件夹
	vartualDir := tool.GetDirName(dirPrefix)
	if vartualDir == nil {
		fmt.Println("Please use 'wshell init' create dir")
		return
	}

	var deleteCmd string = ""
	if len(params) == 0 {
		deleteCmd = "rm -rf " + dirPrefix + "/"
		fmt.Println("delete all files?")
	} else {
		//读取传入文件的绝对路径
		dirPath, _ := filepath.Abs(params[0])
		deleteCmd = "rm -rf " + dirPrefix + "/" + vartualDir[0] + dirPath
		fmt.Println("delete " + dirPath + "?[y/n]")
	}

	var p string
	for true {
		data, _ := fmt.Scanln(&p)
		if data != 1 {
			fmt.Println("input error! Please input 'y' or 'n'")
			continue
		}
		if string(p) == "y" || string(p) == "Y" || string(p) == "n" || string(p) == "N" {
			break
		}
		fmt.Println("input error! Please input again!")
	}

	if string(p) == "y" || string(p) == "Y" {
		if tool.ExecCmdRoutine(deleteCmd) {
			fmt.Println("delete package success!")
			return
		}
		fmt.Println("delete package error!")
	}

}
