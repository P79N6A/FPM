package command

import (
	"FPMtestUpload/dataStruct"
	"client-put/model"
	"client-put/tool"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/*
function:client uploadFile
wshell put name1 name2 name3... 可以同时上传多个压缩包
author:wuwei
input:params ...string
output:void
Data:Thu Aur 23 2018
*/
func Put(cmd string, params ...string) {

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
		destName := vartualDir[0] + ".tar.gz"
		localFile = destName
	}
	for i, v := range strings.Split(localFile, ":") {
		if i >= len(params) {
			return
		}
		version, name, iteration, err := fileParse(dirPrefix + "/" + v)
		if err != nil {
			fmt.Println(err)
			return
		}

		//源码包的信息通过post上传到服务器
		var temp dataStruct.FpmBasicInfo
		temp.Name = name
		temp.Version = version
		temp.Iteration = iteration
		if PostFileInfo(temp) != nil {
			fmt.Println("http upload failed")
		}

		//创建s3对象
		client := model.NewS3Client()
		client.S3Bucket = model.S3Bucket
		client.S3AccessKey = model.S3AccessKey
		client.S3SecretKey = model.S3SecretKey
		client.S3Region = model.S3Region
		client.S3EndPoint = model.S3EndPoint

		path, err := client.UploadFile(v, v)

		if err != nil {
			fmt.Println(path)
			fmt.Println(err)
			return
		}

		fmt.Println(name + "-" + version + "-" + iteration + ".tar.gz upload sucess!!!")
	}

}

/*
function:File name analysis
author:wuwei
input:fName string
output:version, baseName string, err error
Data:Thu Aur 29 2018
*/
func fileParse(fName string) (version, baseName, iteration string, err error) {
	bn := filepath.Base(fName)
	tempStr := strings.Split(bn, "-")
	length := len(tempStr)
	if len(tempStr) < 3 {
		err = fmt.Errorf("%s naming not standard. Please use `wshell build` to create package", bn)
		return
	}

	re := regexp.MustCompile(`\.tar\.gz`)
	if matched := re.MatchString(tempStr[length-1]); !matched {
		err = fmt.Errorf("%s naming not standard. ", fName)
		return
	}
	iter := strings.Split(tempStr[length-1], ".")
	version = tempStr[length-2]
	for i := 0; i < length-2; i++ {
		if i == length-3 {
			baseName += tempStr[i]
		} else {
			baseName += tempStr[i] + "-"
		}
	}
	iteration = iter[0]
	var an int
	an, err = strconv.Atoi(iter[0])
	if err != nil {
		err = fmt.Errorf("%s naming not standard. ", fName)
		fmt.Println(an)
	}
	return
}

/*
function:Post File Info
author:wuwei
input:FileInfo dataStruct.FpmBasicInfo
output:
Data:Thu Nov 29 2018
*/
func PostFileInfo(FileInfo dataStruct.FpmBasicInfo) (erro error) {
	data := make(url.Values)

	data["Name"] = []string{FileInfo.Name}
	data["Version"] = []string{FileInfo.Version}
	data["Iteration"] = []string{FileInfo.Iteration}
	res, err := http.PostForm("http://10.209.16.164:9000/UploadFileInfo", data)
	if err != nil {
		fmt.Println(err.Error())
		erro = err
		return
	}
	defer res.Body.Close()
	return
}
