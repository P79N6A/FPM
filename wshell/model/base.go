package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Version = "v1.0.6"

type CliFunc func(cmd string, params ...string)

type ApiResponse struct {
	Data       interface{} `json:"data"`
	ErrorMsg   string      `json:"errmsg"`
	StatusCode int         `json:"errno, omitempty"`
	Node       string      `json:"node"`
}

type PutResponse struct {
	ErrorMsg   string `json:"errmsg"`
	StatusCode int    `json:"errno, omitempty"`
	Node       string `json:"node"`
}

type ObjectInfo struct {
	NsId       *string `json:"nsId"`
	TemplateId *string `json:"templateId"`
	Version    *string `json:"version"`
	ObjectKey  *string `json:"objectkey"`
	FileName   *string `json:"filename"`
	Md5        *string `json:"md5"`
	Data       []byte  `json:"data"`
}

type Acc struct {
	Name string `json:"name"`
}

func (a *Acc) ToJson() (jsonStr string, err error) {
	jsonData, mErr := json.Marshal(a)
	if mErr != nil {
		err = fmt.Errorf("Marshal account data error, %s", mErr)
		return
	}
	jsonStr = string(jsonData)
	return
}

func (a *Acc) String() string {
	return fmt.Sprintf("Name: %s", a.Name)
}

///////////////////////////////////////////wuwei start/////////////////////////////////////////////////
/*
Description:连接S3的参数设置，注意域名前要加"http://"，否则自动走https，现在不支持https
*/
const (
	S3Bucket    = "upload_rpm_source_package"                //从hulk获取 bucket名称
	S3AccessKey = "Tc8PHxvIkBdAuRypCD0O"                     //从hulk获取
	S3SecretKey = "phbbxowRWoqR3Bm1xzw88DmSiIxUbHIKinhmEdZe" //从hulk获取
	S3Region    = "us-west-2"
	S3EndPoint  = "http://bjcc.s3.addops.soft.360.cn" //从hulk获取
)

type S3Client struct {
	S3Bucket    string //Bucket名称，从hulk可以获取
	S3AccessKey string //AccessKey名称，从hulk可以获取
	S3SecretKey string //SecretKey名称，从hulk可以获取
	S3Region    string
	S3EndPoint  string //S3的域名，从hulk可以获取
}

func NewS3Client() *S3Client {
	return &S3Client{
		S3Region: "us-west-2",
	}
}

/**
*    上传本地文件到S3
*    params :
*    key S3上保存文件的路径（名字）
*    filename 本地文件名字
 */

func (sc *S3Client) UploadFile(key, filename string) (string, error) {
	creds := credentials.NewStaticCredentials(sc.S3AccessKey, sc.S3SecretKey, "")
	config := &aws.Config{
		Region:      aws.String(sc.S3Region),
		Endpoint:    aws.String(sc.S3EndPoint),
		Credentials: creds,
	}
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open file %q, %v", filename, err)
		return "", err
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(sc.S3Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		log.Printf("failed to upload file, %v", err)
		return "", err
	}
	path := result.Location
	return path, nil
}

/**
*    从S3下载文件到本地并保存
*    params :
*    key S3上需要下载的文件的路径（名字）
*    filename 本地文件名字
 */
func (sc *S3Client) DownloadFile(key, filename string) error {
	creds := credentials.NewStaticCredentials(sc.S3AccessKey, sc.S3SecretKey, "")
	config := &aws.Config{
		Region:      aws.String(sc.S3Region),
		Endpoint:    aws.String(sc.S3EndPoint),
		Credentials: creds,
	}
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("failed to create file %q, %v", filename, err)
		return err
	}

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(sc.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("failed to download file, %v", err)
		return err
	}
	return nil
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

///////////////////////////////////////////wuwei End/////////////////////////////////////////////////
