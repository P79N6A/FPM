package command

import (
	"client-put/tool"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	wshellTgz = "/usr/tmp/wshellTgz"
)

/*
function:在wshellTgz目录下建立虚拟文件夹的压缩包
wshell build name1 name2 name3... 支持同时制作多个压缩包
author:wuwei
input:params ...string。可以有一个参数，也可以没有参数，没有参数时默认将“/usr/tmp/wshell/"目录下的文件夹压缩
output:void
Data:Thu Jan 10 2019
*/
func Build(cmd string, params ...string) {
	//处理参数个数为0的情况
	if len(params) == 0 {
		//读取usr/tmp/wshell目录下的内容，即目的文件夹的名称，vartualDir[0]即为所要拿到的文件夹名称，因为该目录下只有一个文件夹
		vartualDir := tool.GetDirName(dirPrefix)
		if vartualDir == nil {
			fmt.Println("Please use 'wshell init' create dir")
			return
		}
		destName := vartualDir[0] + ".tar.gz"
		os.Chdir(dirPrefix)
		if err := tool.TarGz(vartualDir[0], destName); err != nil {
			fmt.Printf("Build failed: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Build package", vartualDir[0]+".tar.gz", "success!")
		return
	}
	for i, _ := range params {
		fileInfo := strings.Split(params[i], "-")
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
			if _, err := os.Stat(params[0]); os.IsNotExist(err) {
				fmt.Printf("Build failed: %s not exist.\r\n", srcPath)
				os.Exit(1)
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
			srcName := srcPath + "-" + version + "-" + iteration
			destName := filepath.Base(srcPath) + "-" + version + "-" + iteration + ".tar.gz"
			os.Symlink(srcPath, srcName)

			if err := tool.TarGz(srcName, destName); err != nil {
				fmt.Printf("Build failed: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Build package", srcPath, "success!")

		} else {
			CmdHelp(cmd)
		}
	}

}
