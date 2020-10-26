package convert

import (
	"github.com/leeli73/goFileView/download"
	"github.com/leeli73/goFileView/utils"
	"io/ioutil"
	"julian/goFileConvert/tools/config"
	"julian/goFileConvert/utils/gupload"
	"julian/goFileConvert/utils/response"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type NowFile struct {
	Md5            string
	Ext            string
	LastActiveTime int64
}

var (
	CosPath 		= "upload/convert/"
	AllOfficeEtx 	= []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}
)

func Handle(w http.ResponseWriter, r *http.Request) {

	//接收参数
	_ = r.ParseForm()
	fileUrl := r.Form.Get("url")
	//var fileList []string

	//返回数组
	returnData := make(map[string]interface{})
	returnData["minIndex"] = 0;
	returnData["maxIndex"] = 0;

	//校验参数
	if fileUrl == "" {
		response.MJson(w, 400, "参数错误", nil)
		return
	}

	//下载文件
	if filePath, err := download.DownloadFile(fileUrl, "cache/download/" + path.Base(fileUrl)); err == nil {

		//获取文件md5名称
		fileName := strings.Split(path.Base(filePath), ".")[0]

		if config.CosConfig.Flag {
			returnData["urlPath"] = config.CosConfig.Url + "/" + CosPath + fileName + "/"
		} else {
			returnData["urlPath"] = config.ApplicationConfig.Domain + "/cache/convert/" + fileName + "/"
		}

		//判断是否已经转换过
		_, err := os.Stat("cache/convert/" + fileName)

		//未转换过，执行转换操作
		if err != nil {
			//上传pdf
			if path.Ext(filePath) == "pdf" {
				if imgPath := utils.ConvertToImg(filePath); imgPath != "" {

				} else {
					response.MJson(w, 400, "转换为图片时出现错误！", nil)
					return
				}
			//上传除pdf外的文件
			} else if utils.IsInArr(path.Ext(filePath), AllOfficeEtx) {
				if pdfPath := utils.ConvertToPDF(filePath); pdfPath != "" {
					if imgPath := utils.ConvertToImg(pdfPath); imgPath != "" {

					} else {
						response.MJson(w, 400, "转换为图片时出现错误！", nil)
						return
					}
				} else {
					response.MJson(w, 400, "转换为PDF时出现错误！", nil)
					return
				}
			//不支持的文件格式
			} else {
				response.MJson(w, 400, "不支持文档格式：" + path.Ext(filePath), nil)
				return
			}
		}

		//获取转化后的所有图片文件
		files, err := ioutil.ReadDir("cache/convert/" + fileName)
		if err != nil {
			response.MJson(w, 400, "读取转化后的文件失败", nil)
			log.Fatal(err)
			return
		}

		//文件上传到cos
		returnData["maxIndex"] = len(files) - 1
		for _, f := range files {
			if config.CosConfig.Flag {
				go uploadImg("cache/convert/" + fileName + "/" + f.Name(), CosPath + fileName + "/")
			}
		}

	//下载失败
	} else {
		log.Println("Error: <", err, "> when download file")
		response.MJson(w, 400, "下载文件失败，请检查你的路径是否正确！", nil)
		return
	}

	response.MJson(w, 200, "success", returnData)
}

//文件上传到cos
func uploadImg(filePath string, cosPath string) {
	_, err := gupload.UploadFileToCos(filePath, cosPath)
	if err != nil {
		go uploadImg(filePath, cosPath)
	}
}
