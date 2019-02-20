package tools

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
function:合并两个字符串数组
author：wuwei
input：s1 []string
input：s2 []string
output：[]string
Data：Thu Aur 15 2018
*/
func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

/*
function:整形转化为字符串数组
author：wuwei
input：i int64
output：[]byte
Data：Thu Aur 16 2018
*/
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

/*
function:Extract .tar.gz file
author:wuwei
input:w http.ResponseWriter
input:r *http.Request
output:void
Data:Thu Aur 22 2018
*/
func TarExtGz(tarFile string) error {

	cmd := exec.Command("tar", "-xzf", tarFile)
	cmd.Dir = "/home/i-wuwei1/FPM/SOURCES/"
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
function:string conversion to json
input:str string
output1:jsonStr []byte
output2:err error
*/
func ToJson(str string) (jsonStr []byte, err error) {
	jsonStr, mErr := json.Marshal(str)
	if mErr != nil {
		err = fmt.Errorf("Marshal account data error, %s", mErr)
		return
	}
	return
}

/*
function:[]byte conversion to string
input:b []byte
output1:string
*/
func ByteToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

/*
function:检查错误输出
input:err error
*/
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
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

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		println("")
	}

	if end < 0 || end > length {
		println("")
	}

	return string(rs[start:end])
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
