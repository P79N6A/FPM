package tool

import (
	"fmt"
	"io/ioutil"
	"os"

	"os/exec"

	"strconv"

	//"furion-stable/fpmfshell/upgrade"
	//"archive/tar"
	//"archive/gzip"
)

func TarGz(path, tarFile string) error {
	cmd := exec.Command("tar", "-zcf", tarFile, "-h", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := ioutil.ReadAll(stdout); err != nil {
		return err
	}
	return nil
}

/*
function:判断字符串是纯数字字符串还是含有字母或者小数点的字符串
author：wuwei
input：str string
output：是数字字符串输出true，否则输出false
*/
func IsNumStr(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	} else {
		return true
	}
}

/*
function:判断所给路径文件/文件夹是否存在
author：wuwei
input：path string
output：存在输出true，否则输出false
*/
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

/*
function:判断所给路径是否为文件夹
author：wuwei
input：path string
output：是文件夹输出true，否则输出false
*/
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

/*
function:判断所给路径是否为文件
author：wuwei
input：path string
output：是文件输出true，否则输出false
*/
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

/*
function:利用exec库的command方法执行linux的命令
author：wuwei
input：Cmd string，要执行的命令行语句
output：执行正确true，失败输出false
*/
func ExecCmdRoutine(Cmd string) bool {
	cmd := exec.Command("/bin/sh", "-c", Cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return false
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return false
	}

	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err)
		return false
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return false
	}
	//fmt.Printf("stdout:\n\n %s", info)
	return true
}

/*
function:读取文件夹下的所有文件夹
author：wuwei
input：path string，输入的文件夹根目录
output：返回文件夹下的所有文件夹的名称数组
*/
func GetDirName(path string) []string {
	var names []string
	dirs, _ := ioutil.ReadDir(path)
	for _, dir := range dirs {
		if dir.IsDir() {
			names = append(names, dir.Name())
		} else {
			continue
		}
	}
	return names
}

/*
function:读取文件夹下的所有文件(文件夹)
author：wuwei
input：path string，输入的文件夹根目录
output：返回文件夹下的所有文件夹的名称数组
*/
func GetFileName(path string) []string {
	var names []string
	dirs, _ := ioutil.ReadDir(path)
	for _, dir := range dirs {
		names = append(names, dir.Name())
	}
	return names
}
