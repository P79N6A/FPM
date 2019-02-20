package command

import (
	"fmt"
	//"io/ioutil"
	"os"
	//"os/exec"
	//"path/filepath"
	"client-put/tool"
	"strings"
)

const (
	dirPrefix = "/usr/tmp/wshell"
)

func Init(cmd string, params ...string) {

	if len(params) != 1 {
		CmdHelp(cmd)
		return
	}
	fileInfo := strings.Split(params[0], "-")
	length := len(fileInfo)
	if length >= 3 {
		var srcPath string
		//后两位是版本号和release号，之前的都是文件名。
		for i := 0; i < length-2; i++ {
			if i == length-3 {
				srcPath += fileInfo[i]
			} else {
				srcPath += fileInfo[i] + "-"
			}
		}

		if matched := strings.Index(srcPath, "/"); matched != -1 {
			fmt.Println("Build failed: the name not support absolute path or include '/'.")
			os.Exit(1)
		}

		version := fileInfo[length-2]
		if matched := strings.Index(version, "-"); matched != -1 {
			fmt.Println("Build failed: the version cannot be include '-'.")
			os.Exit(1)
		}

		iteration := fileInfo[length-1]
		if tool.IsNumStr(iteration) == false {
			fmt.Println("Build failed: the release must bu integers.")
			os.Exit(1)
		}
		if matched := strings.Index(iteration, "-"); matched != -1 {
			fmt.Println("Build failed: the release cannot be include '-'.")
			os.Exit(1)
		}
		srcName := dirPrefix + "/" + srcPath + "-" + version + "-" + iteration

		//创建目录之前先删掉目录下的所有文件
		deleteCmd := "rm -rf " + dirPrefix + "/*"
		tool.ExecCmdRoutine(deleteCmd)

		//os.TempDir()//创建一个临时目录
		mkdirCmd := "mkdir -p " + srcName
		if !tool.ExecCmdRoutine(mkdirCmd) {
			fmt.Println("init error")
			os.Exit(1)
		}
		fmt.Println("Init package", srcPath, "success!")

	} else {
		CmdHelp(cmd)
	}

}
