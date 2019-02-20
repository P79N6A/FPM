package command

import (
	"FPMtestUpload/dataStruct"
	"fmt"
	"io/ioutil"
	"os/exec"
	//"net/http"
	//"FPMtestUpload/releaseRpm"
	"FPMtestUpload/S3Download"
	"FPMtestUpload/tools"
	"strings"
)

/*
function:deal with sh
author:wuwei
input:rpmInfo dataStruct.FpmBasicInfo
output:preinstall postinstall preuninstall postuninstall string
Data:Thu Aur 21 2018
*/
func Getsh(rpmInfo dataStruct.FpmBasicInfo) (preinstall string,
	postinstall string,
	preuninstall string,
	postuninstall string) {

	if len(rpmInfo.PreInstallSh) > 0 {
		preinstall = " --pre-install scripts/" + rpmInfo.PreInstallSh
	} else {
		preinstall = ""
	}

	if len(rpmInfo.PostInstallSh) > 0 {
		postinstall = " --post-install scripts/" + rpmInfo.PostInstallSh
	} else {
		postinstall = ""
	}
	if len(rpmInfo.PreUninstallSh) > 0 {
		preuninstall = " --pre-uninstall scripts/" + rpmInfo.PreUninstallSh
	} else {
		preuninstall = ""
	}

	if len(rpmInfo.PostUninstallSh) > 0 {
		postuninstall = " --post-uninstall scripts/" + rpmInfo.PostUninstallSh
	} else {
		postuninstall = ""
	}
	return
}

/*
function:exec command goroutine
author:wuwei
input:fpmCmd string
input:rpmInfo dataStruct.FpmBasicInfo
output:int
Data:Fri Aur 24 2018
*/
func ExecCmdRoutine(fpmCmd string, rpmInfo dataStruct.FpmBasicInfo) int64 {
	dataStruct.RpmInfoOutput.InfoSet(rpmInfo.Name+rpmInfo.Version, []byte("true"))
	cmd := exec.Command("/bin/sh", "-c", fpmCmd)
	cmd.Dir = "/home/i-wuwei1/FPM/SOURCES/" + rpmInfo.Name + "-" + rpmInfo.Version + "-" + rpmInfo.Iteration //cmd exec dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return 0
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		dataStruct.Err.ErrSet(rpmInfo.Name+rpmInfo.Version, fmt.Errorf("%s", "rpm infomation error,please check up rpmName|version|iteration"))
		return 0
	}

	info, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err)
		return 0
	}

	if err := cmd.Wait(); err != nil {
		dataStruct.Err.ErrSet(rpmInfo.Name+rpmInfo.Version, fmt.Errorf("%s", "rpm information error"))
		fmt.Println("wait:", err.Error())
		return 0
	}
	fmt.Printf("stdout:\n\n %s", info)
	dataStruct.RpmInfoOutput.InfoSet(rpmInfo.Name+rpmInfo.Version, []byte("addops-auto-"+rpmInfo.Name+"-"+rpmInfo.Version+"-"+rpmInfo.Iteration+".el6.x86_64.rpm has been packaged,please test!!!"))
	return 1
}

/*
function:解压文件之前先检查该目录下是否有该文件夹，有则删掉，避免融合
author：wuwei
input：i int64
output：[]byte
Data：Thu Jan 24 2019
*/
func DeleteDir(localFile string) {

	dirPath := tools.Substr2(localFile, 0, len(localFile)-7)

	if tools.Exists("/home/i-wuwei1/FPM/SOURCES/" + localFile) {
		deleteCmd := "rm -rf " + "/home/i-wuwei1/FPM/SOURCES/" + dirPath + "*"
		println(deleteCmd)
		if tools.ExecCmdRoutine(deleteCmd) {
			fmt.Println("delete package success!")
			return
		}
		fmt.Println("delete package error!")
	}

}

/*
function:处理shell字符串，执行shell命令
author：wuwei
input：rpmInfo
output：void
Data：Thu Aur 16 2018
*/
func ProcessCommand(rpmInfo dataStruct.FpmBasicInfo) int64 {

	//如果源码文件没有通过http方式上传的话，就通过s3下载
	fmt.Println("s3:" + rpmInfo.Name + rpmInfo.Version + rpmInfo.Iteration)
	dataStruct.RpmInfoOutput.InfoSet(rpmInfo.Name+rpmInfo.Version, []byte("true"))

	//s3对象赋值
	S3Download.Client.S3Bucket = S3Download.S3Bucket
	S3Download.Client.S3AccessKey = S3Download.S3AccessKey
	S3Download.Client.S3SecretKey = S3Download.S3SecretKey
	S3Download.Client.S3Region = S3Download.S3Region
	S3Download.Client.S3EndPoint = S3Download.S3EndPoint

	localFile := rpmInfo.Name + "-" + rpmInfo.Version + "-" + rpmInfo.Iteration + ".tar.gz"

	DeleteDir(localFile)

	err := S3Download.Client.DownloadFile(localFile, "/home/i-wuwei1/FPM/SOURCES/"+localFile)

	if err != nil {
		fmt.Println(err)
		dataStruct.Err.ErrSet(rpmInfo.Name+rpmInfo.Version, fmt.Errorf("%s", "No source package information，please use 'wshell put fileName' upload file..."))
		return 0
	}

	//解压文件
	errTar := tools.TarExtGz(localFile)
	if errTar != nil {
		fmt.Println(errTar)
		return 0
	}
	//}

	var dependenceListInfo []string
	if rpmInfo.Dependence != "" {
		dependenceListInfo = strings.Split(rpmInfo.Dependence, "#") //dependence libs
	}
	var dependenceInfo string
	dependenceInfo = ""
	for _, v := range dependenceListInfo {
		dependenceInfo += " -d " + string(v)
	}
	preinstall, postinstall, preuninstall, postuninstall := Getsh(rpmInfo)
	fpmCmd := "fpm -f -s dir -t rpm -n addops-auto-" + rpmInfo.Name + " -v " + rpmInfo.Version + " --iteration " + rpmInfo.Iteration + ".el6" + dependenceInfo + " --description \"" + rpmInfo.Description + "\" -p ../../RPMS " + preinstall + postinstall + preuninstall + postuninstall + " ./"
	ExecCmdRoutine(fpmCmd, rpmInfo)
	return 1
}
