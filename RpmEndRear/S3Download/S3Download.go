package S3Download

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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

var Client *S3Client //全局唯一的client

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
