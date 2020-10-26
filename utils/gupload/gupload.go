package gupload

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"julian/goFileConvert/tools/config"
	"julian/goFileConvert/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const(
	FileRootPath	=	"runtime/"
	fileSavePath 	= 	"upload/"
)

type UploadFileRes struct {
	Localpath string
	Remotepath string
}

func UploadHandler(r *http.Request, selfpath string, isCos bool, isLocal bool)  (ue UploadFileRes, e error){

	//接收文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return ue, err
	}

	savePathDir := fileSavePath + selfpath + time.Now().Format("20060102") + "/"
	localpath 	:= FileRootPath + savePathDir + header.Filename
	remotepath 	:= ""

	err = utils.CheckDir(FileRootPath + savePathDir)
	if err != nil {
		return ue, err
	}

	//将文件拷贝到指定路径下，或者其他文件操作
	dst, err := os.Create(localpath)
	if err != nil {
		return ue, err
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		return ue, err
	}

	//关闭文件
	err = dst.Close()
	if err != nil {
		fmt.Printf("close file is Fail %v", err)
	}

	//上传到cos
	if isCos == true {

		u, _ := url.Parse(config.CosConfig.Url)
		b := &cos.BaseURL{BucketURL: u}

		// 永久密钥
		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.CosConfig.Secretid,
				SecretKey: config.CosConfig.Secretkey,
			},
		})
		if client != nil {

			// 对象键（Key）是对象在存储桶中的唯一标识。
			// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
			name := savePathDir + header.Filename

			// 2.通过本地文件上传对象
			_, err = client.Object.PutFromFile(context.Background(), name, "./" + localpath, nil)
			if err == nil {
				remotepath = config.CosConfig.Url + "/" + name
			} else {
				fmt.Printf("COS is Fail %v", err)
			}
		}
	}

	//删除本地文件
	if isLocal == false {
		err := os.Remove("./" + localpath)
		if err != nil {
			fmt.Printf("delete file is fail：%v", err)
		}
		localpath = ""
	}

	ue.Localpath = localpath
	ue.Remotepath = remotepath

	return ue, nil
}

func SaveFileHandler(file io.ReadCloser, fileName string, selfpath string, isCos bool)  (ue UploadFileRes, e error){

	savePathDir := fileSavePath + selfpath + time.Now().Format("20060102") + "/"
	localpath 	:= FileRootPath + savePathDir + fileName
	remotepath 	:= ""

	err := utils.CheckDir(FileRootPath + savePathDir)
	if err != nil {
		return ue, err
	}

	//将文件拷贝到指定路径下，或者其他文件操作
	dst, err := os.Create(localpath)
	if err != nil {
		return ue, err
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		return ue, err
	}

	if isCos == true {
		//上传到cos
		u, _ := url.Parse(config.CosConfig.Url)
		b := &cos.BaseURL{BucketURL: u}

		// 永久密钥
		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.CosConfig.Secretid,
				SecretKey: config.CosConfig.Secretkey,
			},
		})
		if client != nil {

			// 对象键（Key）是对象在存储桶中的唯一标识。
			// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
			name := savePathDir + fileName

			// 2.通过本地文件上传对象
			_, err = client.Object.PutFromFile(context.Background(), name, "./" + localpath, nil)
			if err == nil {
				remotepath = config.CosConfig.Url + "/" + name
			} else {
				fmt.Printf("COS is Fail %v", err)
			}
		}
	}

	ue.Localpath = localpath
	ue.Remotepath = remotepath

	return ue, nil
}

func SaveFileByByteHandler(fileByte []byte, fileName string, selfpath string, isCos bool)  (ue UploadFileRes, e error){

	savePathDir := fileSavePath + selfpath + time.Now().Format("20060102") + "/"
	localpath 	:= FileRootPath + savePathDir + fileName
	remotepath 	:= ""

	err := utils.CheckDir(FileRootPath + savePathDir)
	if err != nil {
		return ue, err
	}

	err = ioutil.WriteFile(localpath, fileByte, 0755)
	if err != nil {
		return ue, err
	}

	if isCos == true {
		//上传到cos
		u, _ := url.Parse(config.CosConfig.Url)
		b := &cos.BaseURL{BucketURL: u}

		// 永久密钥
		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.CosConfig.Secretid,
				SecretKey: config.CosConfig.Secretkey,
			},
		})
		if client != nil {

			// 对象键（Key）是对象在存储桶中的唯一标识。
			// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
			name := savePathDir + fileName

			// 2.通过本地文件上传对象
			_, err = client.Object.PutFromFile(context.Background(), name, "./" + localpath, nil)
			if err == nil {
				remotepath = config.CosConfig.Url + "/" + name
			} else {
				fmt.Printf("COS is Fail %v", err)
			}
		}
	}

	ue.Localpath = localpath
	ue.Remotepath = remotepath

	return ue, nil
}

func UploadFileToCos(filePath string, cosPath string)  (remotePath string, e error){

	//获取文件名
	_, fileName := filepath.Split(filePath)

	//上传到cos
	u, _ := url.Parse(config.CosConfig.Url)
	b := &cos.BaseURL{BucketURL: u}

	// 永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CosConfig.Secretid,
			SecretKey: config.CosConfig.Secretkey,
		},
	})
	if client != nil {

		// 对象键（Key）是对象在存储桶中的唯一标识。
		// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
		name := cosPath + fileName

		// 2.通过本地文件上传对象
		_, err := client.Object.PutFromFile(context.Background(), name, "./" + filePath, nil)
		if err == nil {
			remotePath = config.CosConfig.Url + "/" + name
		} else {
			fmt.Printf("COS is Fail %v\n", err)
		}
	}

	return remotePath, nil
}
